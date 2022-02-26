/* See LICENSE file for copyright and license details. */

/* appearance */
static const unsigned int borderpx  = 2;        /* border pixel of windows */
static const unsigned int snap      = 32;       /* snap pixel */
static const int showbar            = 1;        /* 0 means no bar */
static const int swallowfloating    = 0;        /* 1 means swallow floating windows by default */
static const int topbar             = 1;        /* 0 means bottom bar */

static const char *fonts[]          = { "monospace", "emoji" };
static const char dmenufont[]       = "monospace";

static const char col_gray1[]       = "#282828";
static const char col_gray2[]       = "#282828";
static const char col_gray3[]       = "#bbbbbb";
static const char col_gray4[]       = "#eeeeee";
static const char col_red[]         = "#d71513";
static const char col_yellow[]      = "#ffee00";
static const char col_cyan[]        = "#d71513";
static const char *colors[][3]      = {
	/*               fg         bg         border   */
	[SchemeNorm] = { col_gray3, col_gray1, col_gray1 },
	[SchemeSel]  = { col_gray1, col_red,  col_gray4  },
};

/* tagging */
static const char *tags[] = { "1", "2", "3", "4", "5", "6", "7", "8", "9" };

static const Rule rules[] = {
	/* xprop(1):
	 *	WM_CLASS(STRING) = instance, class
	 *	WM_NAME(STRING) = title
	 */
	/* class     instance  title           tags mask  isfloating  isterminal  noswallow  monitor */
	{ "St",      NULL,     NULL,           0,         0,          1,          0,         -1 },
	{ NULL,      NULL,     "Event Tester", 0,         0,          0,          1,         -1 },
};

/* layout(s) */
static const float mfact     = 0.5; /* factor of master area size [0.05..0.95] */
static const int nmaster     = 1;    /* number of clients in master area */
static const int resizehints = 0;    /* 1 means respect size hints in tiled resizals */

#include "fibonacci.c"
#include "tatami.c"

static const Layout layouts[] = {
	/* symbol     arrange function */
	{ "[]=",      tile },    /* first entry is default */
 	{ "[@]",      spiral },
	{ "|+|",      tatami },
	{ "><>",      NULL },    /* no layout function means floating behavior */
	{ "[M]",      monocle },
 	{ "[\\]",     dwindle },
	{ NULL,		NULL },
};

/* key definitions */
#define MODKEY Mod4Mask
#define TAGKEYS(KEY,TAG) \
	{ MODKEY,                       KEY,      view,           {.ui = 1 << TAG} }, \
	{ MODKEY|ControlMask,           KEY,      toggleview,     {.ui = 1 << TAG} }, \
	{ MODKEY|ShiftMask,             KEY,      tag,            {.ui = 1 << TAG} }, \
	{ MODKEY|ControlMask|ShiftMask, KEY,      toggletag,      {.ui = 1 << TAG} },

#define STACKKEYS(MOD,ACTION) \
	{ MOD,	XK_j,	ACTION##stack,	{.i = +1 } }, \
	{ MOD,	XK_k,	ACTION##stack,	{.i = -1 } }, \

/* helper for spawning shell commands in the pre dwm-5.0 fashion */
#define SHCMD(cmd) { .v = (const char*[]){ "/bin/sh", "-c", cmd, NULL } }

/* commands */
static char dmenumon[2] = "0"; /* component of dmenucmd, manipulated in spawn() */
static const char *dmenucmd[] = { "dmenu_run", "-m", dmenumon, "-fn", dmenufont, "-nb", col_gray1, "-nf", col_gray3, "-sb", col_cyan, "-sf", col_gray4, NULL };
static const char *termcmd[]  = { "st", NULL };

/* commands spawned when clicking statusbar, the mouse button pressed is exported as BUTTON */
static const StatusCmd statuscmds[] = {
	{ "notify-send Mouse$BUTTON", 1 },
};
static const char *statuscmd[] = { "/bin/sh", "-c", NULL, NULL };


static Key keys[] = {
	/* modifier                     key        function        argument */
	STACKKEYS(MODKEY,                          focus)
	STACKKEYS(MODKEY|ShiftMask,                push)
	{ MODKEY,                       XK_d,      spawn,          {.v = dmenucmd } },
	{ MODKEY,                       XK_Return, spawnsshaware,  {.v = termcmd } },
	/* { MODKEY,                       XK_b,      togglebar,      {0} }, */
	/* { MODKEY,                       XK_j,      focusstack,     {.i = +1 } }, */
	/* { MODKEY,                       XK_k,      focusstack,     {.i = -1 } }, */
	/* { MODKEY,                       XK_i,      incnmaster,     {.i = +1 } }, */
	/* { MODKEY,                       XK_d,      incnmaster,     {.i = -1 } }, */
	/* { MODKEY,                       XK_h,      setmfact,       {.f = -0.05} }, */
	/* { MODKEY,                       XK_l,      setmfact,       {.f = +0.05} }, */
	/* { MODKEY,                       XK_Return, zoom,           {0} }, */
	/* { MODKEY,                       XK_Tab,    view,           {0} }, */
	/* { MODKEY|ShiftMask,             XK_c,      killclient,     {0} }, */
  { MODKEY|ShiftMask,             XK_l,      spawn,         SHCMD("slock")  },
	{ MODKEY,			                  XK_q,		   killclient,	   {0} },
	/* { MODKEY,                       XK_t,      setlayout,      {.v = &layouts[0]} }, */
	/* { MODKEY,                       XK_f,      setlayout,      {.v = &layouts[1]} }, */
	/* { MODKEY,                       XK_m,      setlayout,      {.v = &layouts[2]} }, */
	/* { MODKEY,                       XK_space,  setlayout,      {0} }, */
	/* { MODKEY|ShiftMask,             XK_space,  togglefloating, {0} }, */

	{ MODKEY,                       XK_r,     spawn,      SHCMD("st -e ranger") },
	{ MODKEY|ShiftMask,             XK_r,     spawn,      SHCMD("st -e lf") },

	{ MODKEY,                       XK_g,     spawn,      SHCMD("chromium") },
	{ MODKEY|ShiftMask,             XK_g,     spawn,      SHCMD("firefox") },

	{ MODKEY,                       XK_m,     spawn,      SHCMD("mirage") },
	{ MODKEY|ShiftMask,             XK_m,     spawn,      SHCMD("element-desktop") },

	{ MODKEY,                       XK_w,     spawn,      SHCMD("st -e sudo nmtui") },
	{ MODKEY,                       XK_h,     spawn,      SHCMD("st -e sudo -E htop") },

	{ MODKEY,			                  XK_minus,	spawn,		SHCMD("pamixer --allow-boost -d 5; pkill -RTMIN+9 dwmblocks") },
	{ MODKEY,			                  XK_equal,	spawn,		SHCMD("pamixer --allow-boost -i 5; pkill -RTMIN+9 dwmblocks") },

	{ MODKEY|ShiftMask,			        XK_minus,	spawn,		SHCMD("sudo brightness down; pkill -RTMIN+10 dwmblocks") },
	{ MODKEY|ShiftMask,			        XK_equal, spawn,		SHCMD("sudo brightness up; pkill -RTMIN+10 dwmblocks") },

	{ MODKEY,			                  XK_f,		   togglefullscr,	 {0} },
	{ MODKEY,                       XK_0,      view,           {.ui = ~0 } },
	{ MODKEY|ShiftMask,             XK_0,      tag,            {.ui = ~0 } },
	{ MODKEY,                       XK_comma,  focusmon,       {.i = -1 } },
	{ MODKEY,                       XK_period, focusmon,       {.i = +1 } },
	{ MODKEY|ShiftMask,             XK_comma,  tagmon,         {.i = -1 } },
	{ MODKEY|ShiftMask,             XK_period, tagmon,         {.i = +1 } },
	TAGKEYS(                        XK_1,                      0)
	TAGKEYS(                        XK_2,                      1)
	TAGKEYS(                        XK_3,                      2)
	TAGKEYS(                        XK_4,                      3)
	TAGKEYS(                        XK_5,                      4)
	TAGKEYS(                        XK_6,                      5)
	TAGKEYS(                        XK_7,                      6)
	TAGKEYS(                        XK_8,                      7)
	TAGKEYS(                        XK_9,                      8)
	{ MODKEY|ShiftMask,             XK_q,      quit,           {0} },
};

/* button definitions */
/* click can be ClkTagBar, ClkLtSymbol, ClkStatusText, ClkWinTitle, ClkClientWin, or ClkRootWin */
static Button buttons[] = {
	/* click                event mask      button          function        argument */
	{ ClkLtSymbol,          0,              Button1,        setlayout,      {0} },
	{ ClkLtSymbol,          0,              Button3,        setlayout,      {.v = &layouts[2]} },
	{ ClkWinTitle,          0,              Button2,        zoom,           {0} },
	{ ClkStatusText,        0,              Button2,        spawn,          {.v = termcmd } },
	{ ClkClientWin,         MODKEY,         Button1,        movemouse,      {0} },
	{ ClkClientWin,         MODKEY,         Button2,        togglefloating, {0} },
	{ ClkClientWin,         MODKEY,         Button3,        resizemouse,    {0} },
	{ ClkTagBar,            0,              Button1,        view,           {0} },
	{ ClkTagBar,            0,              Button3,        toggleview,     {0} },
	{ ClkTagBar,            MODKEY,         Button1,        tag,            {0} },
	{ ClkTagBar,            MODKEY,         Button3,        toggletag,      {0} },
};

