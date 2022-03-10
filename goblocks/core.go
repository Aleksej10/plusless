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
  "sync"
)

type Block struct {
  signal int
  interval float64
  icon string
  command string
}

type IndexTimestamp struct {
  index int
  last_changed time.Time
}

var (
  results []string
  sig2decayingblock = make(map[int]IndexTimestamp)
  should_update bool = false
  mutex sync.Mutex
)

const (
  SIGRTMIN int = 34
  SIGRTMAX int = 64
)

func sig2int(sig os.Signal) int {
  reg, err := regexp.Compile("[^0-9]")

  if err != nil {
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

    go bind_channel(&ch, i)
    go start_block(i)
  } 

  go start_drawing(min/2)

  return ch
}

func bind_channel(ch *chan os.Signal, i int) {
  sig := blocks[i].signal

  signal.Notify(*ch, syscall.Signal(sig + SIGRTMIN))

  if blocks[i].interval == 0 {
    sig2decayingblock[sig] =  IndexTimestamp { i, time.Now() }
  }
}

func start_block(i int) {
  interval := time.Duration(blocks[i].interval)

  if interval > 0 {
    for {
      exec_block(i)
      time.Sleep(interval * time.Second)
    }
  } 
}

func update_block(i int) {
  exec_block(i)
  draw_blocks()
}

func exec_command(command string) string {
  cmd := exec.Command(SHELL, "-c", command)
  stdout, err := cmd.StdoutPipe()

  if err != nil {
    return ""
  }

  if err := cmd.Start(); err != nil {
    return ""
  }

  data, err := ioutil.ReadAll(stdout)

  if err != nil {
    return ""
  }

  if err := cmd.Wait(); err != nil {
    return ""
  }

  return string(data)
}

func exec_block(i int) {
  res := exec_command(blocks[i].command)

  if results[i] == res {
    return
  }

  results[i] = res

  for sig, v := range sig2decayingblock {
    if v.index == i {
      sig2decayingblock[sig] = IndexTimestamp { i, time.Now() }
      break
    }
  }

  mutex.Lock()
    should_update = true
  mutex.Unlock()
}

func start_drawing(interval float64){
  go draw_blocks()

  if !math.IsInf(interval, 1) {
    for {
      time.Sleep(time.Duration(interval) * time.Second)
      clear_blocks()
      draw_blocks()
    }
  }
} 

func clear_blocks() {
  now := time.Now()

  for _, v := range sig2decayingblock {
    if blocks[v.index].interval != 0 {
      continue
    }

    if (results[v.index] != "") && (now.Sub(v.last_changed) > (2 * time.Second)) {
      results[v.index] = ""

      mutex.Lock()
        should_update = true
      mutex.Unlock()
    }
  }
}

func draw_blocks(){
  if should_update {
    go update_dwm_status(status_string())
    
    mutex.Lock()
      should_update = false
    mutex.Unlock()
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

    if (r != "") {
      status += fmt.Sprintf("%c%v%v ", blocks[i].signal,  strings.TrimSpace(blocks[i].icon + " " + r), strings.TrimSpace(" " + DELIM))
    }
  }

  return strings.TrimRight(status, " " + DELIM)
}

func main() {
  ch := initialize()

  for {
    sig := <-ch
    
    if sig_num := sig2int(sig); sig_num > 0 {
      for i, block := range blocks {
        if sig_num == block.signal {
          fmt.Println("got fucking signal", sig_num)

          go update_block(i)
          break
        }
      }
    }
  }
}
