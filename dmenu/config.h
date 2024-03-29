static int topbar = 1;                      /* -b  option; if 0, dmenu appears at bottom     */
static int fuzzy = 1;
static int managed = 0;
/* -fn option overrides fonts[0]; default X11 font or font set */

static const char *fonts[] = { "monospace:pixelsize=14" };

static const char *prompt      = NULL;      /* -p  option; prompt to the left of input field */
static const char *colors[SchemeLast][2] = {
	/*     fg         bg       */
  [SchemeNorm] = { "#dedede", "#000000" },
  [SchemeSel] = { "#efefef", "#ce2029" },
  [SchemeOut] = { "#000000", "#ffd800" },
};
/* -l option; if nonzero, dmenu uses vertical list with given number of lines */
static unsigned int lines      = 6;

/*
 * Characters not considered part of a word while deleting words
 * for example: " /?\"&[]"
 */
static const char worddelimiters[] = " ";
