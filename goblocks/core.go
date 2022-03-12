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

const SIGRTMIN int = 34

type Block struct {
  signal int
  interval float64
  fadeout bool
  icon string
  command string
}

type Sig struct { block, button int }
type BlockIndex struct { block *Block; index int }

var (
  results []string
  sig2block map[int]BlockIndex
  index2last_updated = make(map[int]time.Time)
  should_update bool = false
  mutex sync.Mutex
)

func log_err(msg string){
  fmt.Fprintln(os.Stderr, msg)
  os.Exit(1)
}

func set_should_update(v bool){
  mutex.Lock()
  should_update = v
  mutex.Unlock()
}

func sig2int(sig os.Signal) int {
  reg, err := regexp.Compile("[^0-9]")

  if err != nil {
    return -1
  }
  i := -1
  fmt.Sscan(reg.ReplaceAllString(sig.String(), ""), &i)

  return i - SIGRTMIN
}

func initialize() chan os.Signal {
  ch := make(chan os.Signal, 1)
  go bind_button_events(&ch)

  min := math.Inf(1)

  for sig, bi := range sig2block {
    interval := bi.block.interval

    if (interval != 0) && (interval < min) {
      min = interval
    }
    go signal.Notify(ch, syscall.Signal(sig + SIGRTMIN))
    go start_block(bi)
  }
  go start_drawing(min/2)

  return ch
}

func start_block(bi BlockIndex){
  if bi.block.fadeout {
    index2last_updated[bi.index] = time.Now()
  } else{
    go exec_block(bi, 0)
  }
  if interval := time.Duration(bi.block.interval) * time.Second; interval > 0 {
    for {
      time.Sleep(interval)
      exec_block(bi, 0)
    }
  }
}

func update_block(sig Sig) {
  bi := sig2block[sig.block]
  exec_block(bi, sig.button)
  draw_blocks()
}

func exec_command(command string, button_sig int) string {
  cmd := exec.Command(SHELL, "-c", command)
  stdout, err := cmd.StdoutPipe()

  if err != nil {
    return ""
  }
  if button_sig > 0 {
    cmd.Env = os.Environ()
    cmd.Env = append(cmd.Env, fmt.Sprintf("BLOCK_BUTTON=%v", button_sig))
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

func exec_block(bi BlockIndex, button_sig int) {
  res := exec_command(bi.block.command, button_sig)

  if results[bi.index] == res {
    return
  }
  results[bi.index] = res

  if bi.block.fadeout {
    index2last_updated[bi.index] = time.Now()
  }
  set_should_update(true)
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
  ok := false

  for i, last_updated := range index2last_updated {
    if (results[i] != "") && (now.Sub(last_updated) > (2 * time.Second)) {
      results[i] = ""
      ok = true
    }
  }
  if ok {
    set_should_update(true)
  }
}

func draw_blocks(){
  if should_update {
    go update_dwm_status()
    set_should_update(false)
  }
}

func update_dwm_status() {
  d := C.XOpenDisplay(nil);
  defer C.XCloseDisplay(d)

	screen := C.XDefaultScreenOfDisplay(d);
  root := C.XRootWindowOfScreen(screen);
	C.XStoreName(d, root, C.CString(status_string()))
}

func status_string() string {
  status := ""

  for i, res := range results {
    if r := strings.TrimSpace(res); r != "" {
      status += fmt.Sprintf("%c%v%v ", blocks[i].signal,  strings.TrimSpace(blocks[i].icon + " " + r), strings.TrimSpace(" " + DELIM))
    }
  }
  return strings.TrimRight(status, " " + DELIM)
}

func parse_signal(ch *chan os.Signal) (sig Sig) {
  sig.block = sig2int(<-*ch)

  if !is_block_signal(sig.block) {
    return Sig { 0, 0 }
  }
  sig.button = sig2int(<-*ch)

  if sig.block == sig.button {
    sig.button = 0
    return
  }
  if (sig.button > 0) || (sig.button < 7) {
    return
  }
  return Sig { 0, 0 }
}

func is_block_signal(sig int) (yes bool) {
  _, yes = sig2block[sig]
  return
}

func map_blocks(blocks *[]Block) (m map[int]BlockIndex) {
  m = make(map[int]BlockIndex)
  tmp_sig := 6

  for i, block := range *blocks {
    bi := BlockIndex { &((*blocks)[i]), i }

    if block.signal != 0 {
      if block.signal < 7 {
        log_err("use signals in range [7..30]")
      }
      m[block.signal] = bi
      continue
    }
    tmp_sig = next_tmp_sig(tmp_sig, blocks)

    if tmp_sig == 0 {
      log_err("you probably have waaay too many blocks")
    }
    m[tmp_sig] = bi
  }

  return
}

func bind_button_events(ch *chan os.Signal) {
  for i := 1; i < 7; i++ {
    signal.Notify(*ch, syscall.Signal(i + SIGRTMIN))
  }
}

func next_tmp_sig(tmp_sig int, blocks *[]Block) (next_sig int) {
  next_sig = tmp_sig + 1

  for next_sig < 31 {
    ok := true

    for _, block := range *blocks {
      if next_sig == block.signal {
        next_sig++
        ok = false
        break
      }
    }
    if ok {
      return
    }
  }
  return 0
}

func correct_custom_signals(){
  for sig, bi := range sig2block {
    blocks[bi.index].signal = sig
  }
}

func main() {
  sig2block = map_blocks(&blocks)
  correct_custom_signals()
  results = make([]string, len(sig2block))

  ch := initialize()

  for {
    sig := parse_signal(&ch)

    if sig.block > 0 {
      go update_block(sig)
    }
  }
}
