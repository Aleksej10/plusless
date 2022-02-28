/* user and group to drop privileges to */
static const char *user  = "nobody";
static const char *group = "nogroup";

static const char *colorname[NUMCOLS] = {
	[INIT] =   "#000000",   /* after initialization */
	[INPUT] =  "#282828",   /* during input */
	[FAILED] = "#d71513",   /* wrong password */
};

/* lock screen opacity */
static const float alpha = 0.2;

/* treat a cleared input like a wrong password (color) */
static const int failonclear = 1;

/* default message */
static const char * message = "";

/* text color */
static const char * text_color = "#eeeeee";

/* text size (must be a valid size) */
static const char * text_size = "fixed";
