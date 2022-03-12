package main

var blocks = []Block {
  // 0 for signal and interval means don't listen and don't update, respectively
  // fadeout fields will fade out FADEOUT_TIME seconds after beign drawn
  // zero signals will be filled with remaining numbers in order to support button events

  // signals [1..6] are reserved for button events
  // use numbers 7-30 for custom signals

  /*
  signal interval fadeout icon  command
  */

  {0,    5,       false,  "",   "sb-cpu" },
  {0,    5,       false,  "",   "sb-internet" },
  {0,    5,       false,  "💿", "df -h | awk '{ if ($6 == \"/home\") print $4 }'" },
  {0,    5,       false,  "☎️",  "adb devices -l | grep -e 'device:' | sed -rn 's/.*?device:(\\w+)\\s+.*/\\1/p' | xargs" },
  {0,    300,     false,  "💲", "curl -s rate.sx/1xmr | awk '{print int($1)}'" },
  {10,   0,       false,  "☀️",  "sb-brightness" },
  {9,    0,       false,  "🔊", "sb-volume" },
  {0,    5,       false,  "",   "sb-battery" },
  {0,    5,       false,  "💾", "cnt=$(($(lsblk -dn | wc -l) - 1)); [ $cnt -gt 0 ] && echo $cnt 'new device'" }, // use parted instead
  {0,    5,       false,  "",   "sb-date" },
}

const (
  DELIM string = ""
  SHELL string = "dash"
  FADEOUT_TIME float64 = 2
)

// button signals:
// 1 - left click
// 2 - middle click
// 3 - right click
// 4 - scroll down
// 5 - scroll up
// 6 - shift + left click
