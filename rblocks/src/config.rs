use crate::CfgBlk;

#[allow(unused)]
pub const BLOCKS: &'static [CfgBlk] = &[
  // 0 for signal and interval means don't listen and don't update, respectively
  // zero signals will be filled with remaining numbers in order to support button events
  // t field is interval in seconds

  // signals [1..6] are reserved for button events
  // use numbers 7-30 for custom signals
  CfgBlk { sig: 0,    t: 600,   ico: "",     cmd: "bash $HOME/tmp/tariffs/automatic/download_tariff.sh" },
  CfgBlk { sig: 0,    t: 5,     ico: "",     cmd: "sb-cpu" },
  CfgBlk { sig: 0,    t: 5,     ico: "",     cmd: "sb-internet" },
  CfgBlk { sig: 0,    t: 5,     ico: "💿",   cmd: "df -h | awk '{ if ($6 == \"/home\") print $4 }'" },
  CfgBlk { sig: 0,    t: 3,     ico: "☎️",    cmd: "adb devices -l | grep -e 'device:' | sed -rn 's/.*?device:(\\w+)\\s+.*/\\1/p' | xargs" },
  CfgBlk { sig: 0,    t: 300,   ico: "💲",   cmd: "curl -s rate.sx/1xmr | awk '{print int($1)}'" },
  CfgBlk { sig: 10,   t: 0,     ico: "☀️",    cmd: "sb-brightness" },
  CfgBlk { sig: 9,    t: 0,     ico: "🔊",   cmd: "sb-volume" },
  CfgBlk { sig: 0,    t: 10,    ico: "",     cmd: "sb-battery" },
  CfgBlk { sig: 0,    t: 3,     ico: "💾",   cmd: "cnt=$(($(lsblk -dn | wc -l) - 1)); [ $cnt -gt 0 ] && echo $cnt 'new device'" }, // use parted instead
  CfgBlk { sig: 0,    t: 3,     ico: "",     cmd: "sb-date" },
  CfgBlk { sig: 0,    t: 0,     ico: "",     cmd: "sb-sys" },
];

#[allow(unused)]
pub const DELIM: &'static str = " ";
#[allow(unused)]
pub const SHELL: &'static str = "dash";

// button signals:
// 1 - left click
// 2 - middle click
// 3 - right click
// 4 - scroll up
// 5 - scroll down
// 6 - shift + left click

// NOTE: scroll direction might differ on touchpads based on your settings
