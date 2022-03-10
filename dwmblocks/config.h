//Modify this file to change what commands output to your statusbar, and recompile using the make command.
static const Block blocks[] = {
	/*Icon*/	/*Command*/		/*Update Interval*/	/*Update Signal*/
	/* {"⌨", "sb-kbselect", 0, 30}, */
	/* {"", "cat /tmp/recordingicon 2>/dev/null",	0,	9}, */
	/* {"",	"sb-tasks",	10,	26}, */
	/* {"",	"sb-music",	0,	11}, */
	/* {"",	"sb-pacpackages",	0,	8}, */
	/* {"",	"sb-news",		0,	6}, */
	/* {"",	"sb-price lbc \"LBRY Token\" 📚",			9000,	22}, */
	/* {"",	"sb-price bat \"Basic Attention Token\" 🦁",	9000,	20}, */
	/* {"",	"sb-price link \"Chainlink\" 🔗",			300,	25}, */
	/* {"",	"sb-price xmr \"Monero\" 🔒",			9000,	24}, */
	/* {"",	"sb-price eth Ethereum 🍸",	9000,	23}, */
	/* {"",	"sb-price btc Bitcoin 💰",				9000,	21}, */
	/* {"",	"sb-torrent",	20,	7}, */
	/* {"",	"sb-memory",	10,	14}, */
	/* {"",	"sb-cpu",		10,	18}, */
	/* {"",	"sb-moonphase",	18000,	17}, */
	/* {"",	"sb-forecast",	18000,	5}, */
	/* {"",	"sb-mailbox",	180,	12}, */
	/* {"",	"sb-nettraf",	1,	16}, */
	/* {"",	"sb-volume",	0,	10}, */
	/* {"",	"sb-battery",	5,	3}, */


  {"", "sb-internet", 60, 12},
  /* {" 💿 ", "df -h | awk '{ if ($6 == \"/home\") print $4 }'", 300, 51}, */
	/* {" 💲 ", "curl -s rate.sx/1xmr | awk '{print int($1)}'",	900,	22}, */
  /* {" ☀️ ", "brightness", 0, 10 }, */
  /* {" 🔊 ", "echo $(pamixer --get-volume)%", 0, 9}, */
  /* {" 🔋 ", "echo \"$(cat /sys/class/power_supply/BAT0/capacity)%\"", 5, 3}, */
  {" 💾 ", "cnt=$(($(lsblk -dn | wc -l) - 1)); [ $cnt -gt 0 ] && echo $cnt 'new device'", 5, 2}, // use parted in future
	{" 🚢 ", "date +'%A(%u) %d. %B(%-m) ``%y %I:%M %p'",	60,	1}
	/* {" ⌨️ ", "setxkbmap -query | awk '/layout/{ print $2 }'",	0,	50}, */
	/* {"",	"sb-internet",	5,	4}, */
	/* {"",	"sb-help-icon",	0,	15}, */
};

//Sets delimiter between status commands. NULL character ('\0') means no delimiter.
static char const *delim = "|";

// Have dwmblocks automatically recompile and run when you edit this file in
// vim with the following line in your vimrc/init.vim:

// autocmd BufWritePost ~/.local/src/dwmblocks/config.h !cd ~/.local/src/dwmblocks/; sudo make install && { killall -q dwmblocks;setsid dwmblocks & }

