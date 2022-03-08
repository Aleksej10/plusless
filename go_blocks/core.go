package main

// #cgo LDFLAGS: -lX11
// #include <X11/Xlib.h>
import "C"

import (
  "fmt"
  "os"
  "os/signal"
  "os/exec"
  "io/ioutil"
  "syscall"
  "regexp"
  "time"
  "math"
  "strings"
)

type Block struct {
  signal int
  interval float64
  icon string
  command string
}


var (
  results []string
  last_result string
  sig_2_block = make(map[int]int)
)

const (
  SIGRTMIN int = 34
  SIGRTMAX int = 64
)

func sig_to_int(sig os.Signal) int {
  reg, err := regexp.Compile("[^0-9]")

  if err != nil {
    fmt.Fprintln(os.Stderr, err)
    return -1
  }

  var i int = -1
  fmt.Sscan(reg.ReplaceAllString(sig.String(), ""), &i)

  return i - SIGRTMIN
}

func initialize() chan os.Signal {
  results = make([]string, len(blocks))
  ch := make(chan os.Signal, 1)

  min := math.Inf(1)

  for i, block := range blocks {
    if (block.interval != 0) && (block.interval < min) {
      min = block.interval
    }

    bind_channel(&ch, block.signal, i)
    go start_block(i)
  } 

  go start_drawing(min/2)

  return ch
}

func bind_channel(ch *chan os.Signal, sig int, i int) {
  if sig > 0 {
    signal.Notify(*ch, syscall.Signal(sig + SIGRTMIN))
    sig_2_block[sig] = i
  }
}

func start_block(i int) {
  exec_block(i)

  interval := time.Duration(blocks[i].interval)

  if interval > 0 {
    for {
      time.Sleep(interval * time.Second)
      exec_block(i)
    }
  } 
}

func update_block(i int) {
  exec_block(i)
  draw_blocks()
}

func exec_command(command string) string {
  cmd := exec.Command("dash", "-c", command)
  stdout, err := cmd.StdoutPipe()

  if err != nil {
    fmt.Fprintln(os.Stderr, err)
    return ""
  }

  if err := cmd.Start(); err != nil {
    fmt.Fprintln(os.Stderr, err)
    return ""
  }

  data, err := ioutil.ReadAll(stdout)

  if err != nil {
    fmt.Fprintln(os.Stderr, err)
    return ""
  }

  if err := cmd.Wait(); err != nil {
    fmt.Fprintln(os.Stderr, err)
    return ""
  }

  return string(data)
}

func exec_block(i int) {
  results[i] = exec_command(blocks[i].command)
}

func start_drawing(interval float64){
  draw_blocks()

  if !math.IsInf(interval, 1) {
    for {
      time.Sleep(time.Duration(interval) * time.Second)
      draw_blocks()
    }
  }
} 

func draw_blocks(){
  s := status_string()

  if s != last_result {
    last_result = s
    update_dwm_status(s)
  }
}

func update_dwm_status(status string) {
  d := C.XOpenDisplay(nil);
  defer C.XCloseDisplay(d)

	screen := C.XDefaultScreenOfDisplay(d);
  root := C.XRootWindowOfScreen(screen);
	C.XStoreName(d, root, C.CString(status))
}

func status_string() string {
  status := ""

  for i, res := range results {
    r := strings.TrimSpace(res)
    if r != "" {
      status += fmt.Sprintf("%v %v %v ", blocks[i].icon, r, delim)
    }
  }

  status = strings.TrimRight(status, " " + delim)

  return status
}

func main() {
  ch := initialize()

  for {
    sig := <-ch
    i := sig_to_int(sig)

    if i > 0 {
      update_block(sig_2_block[i])
    }
  }
}
