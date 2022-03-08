package main

var blocks = []Block {
  // 0 for signal and interval means don't listen and don't update, respectively
  // you usually don't want more than one 0 in this two fields

  // { signal, interval, icon, command }

  {0,  5,   "🌎", "nmcli con show --active | awk 'NR == 2 { print  $1 }'" },
  {0,  5,   "💿", "df -h | awk '{ if ($6 == \"/home\") print $4 }'" },
  {0,  5,   "☎️",  "adb devices -l | grep -e 'device:' | sed -rn 's/.*?device:(\\w+)\\s+.*/\\1/p' | xargs" },
  {0,  300, "💲", "curl -s rate.sx/1xmr | awk '{print int($1)}'" },
  {10, 0,   "☀️",  "brightness" },
  {9,  0,   "🔊", "echo $(pamixer --get-volume)%" },
  {0,  5,   "🔋", "echo \"$(cat /sys/class/power_supply/BAT0/capacity)%\"" },
  {0,  5,   "💾", "cnt=$(($(lsblk -dn | wc -l) - 1)); [ $cnt -gt 0 ] && echo $cnt 'new device'" }, // use parted instead
  {0,  5,   "🚢", "date +'%A(%u) %d. %B(%-m) ``%y %I:%M %p'" },
}

const delim string = "|"
