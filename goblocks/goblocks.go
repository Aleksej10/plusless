package main

var blocks = []Block {
  // 0 for signal and interval means don't listen and don't update, respectively
  // you usually don't want more than one 0 in this two fields

  // UPDATE: don't put 0 for signal even if you don't intent to use it
  // this is to support button events for all blocks

  // { signal, interval, icon, command }

  {20, 5,   "",   "sb-cpu" },
  {21, 5,   "",   "sb-internet" },
  {1,  5,   "💿", "df -h | awk '{ if ($6 == \"/home\") print $4 }'" },
  {2,  5,   "☎️",  "adb devices -l | grep -e 'device:' | sed -rn 's/.*?device:(\\w+)\\s+.*/\\1/p' | xargs" },
  {3,  300, "💲", "curl -s rate.sx/1xmr | awk '{print int($1)}'" },
  {10, 0,   "☀️",  "brightness" },
  {9,  0,   "🔊", "echo $(pamixer --get-volume)%" },
  {22, 5,   "",   "sb-battery" },
  {4,  5,   "💾", "cnt=$(($(lsblk -dn | wc -l) - 1)); [ $cnt -gt 0 ] && echo $cnt 'new device'" }, // use parted instead
  {5,  5,   "🚢", "date +'%A(%u) %d. %B(%-m) `%y %I:%M %p'" },
}

const (
  DELIM string = ""
  SHELL string = "dash"
)
