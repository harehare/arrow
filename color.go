package main

import (
	"os"
	"regexp"
	"strconv"

	"github.com/samber/mo"
)

type Color string

var rgbPattern = regexp.MustCompile("^#?([0-9a-fA-F]{6})$")

var ForegroundColor = foreground()
var HighlightColor = highlight()
var CursorColor = cursor()
var DisabledColor = disabled()
var SymLinkColor = symLink()
var CurrentDirectoryColor = currentDirectory()
var PromptColor = prompt()
var BorderColor = border()

func (c Color) String() string {
	return string(c)
}

func symLink() string {
	e := os.Getenv("ARROW_SYMLINK_COLOR")
	if e == "" {
		return "36"
	}

	return ColorFromString(e).OrElse(Color("36")).String()
}

func prompt() string {
	e := os.Getenv("ARROW_PROMPT_COLOR")
	if e == "" {
		return "36"
	}

	return ColorFromString(e).OrElse(Color("36")).String()
}

func currentDirectory() string {
	e := os.Getenv("ARROW_CURRENT_DIRECTORY_COLOR")
	if e == "" {
		return "57"
	}

	return ColorFromString(e).OrElse(Color("57")).String()
}

func disabled() string {
	e := os.Getenv("ARROW_DISABLED_COLOR")
	if e == "" {
		return "240"
	}

	return ColorFromString(e).OrElse(Color("240")).String()
}

func cursor() string {
	e := os.Getenv("ARROW_CURSOR_COLOR")
	if e == "" {
		return "57"
	}

	return ColorFromString(e).OrElse(Color("57")).String()
}

func foreground() string {
	e := os.Getenv("ARROW_FOREGROUND_COLOR")
	if e == "" {
		return "15"
	}

	return ColorFromString(e).OrElse(Color("15")).String()
}

func highlight() string {
	e := os.Getenv("ARROW_HIGHLIGHT_COLOR")
	if e == "" {
		return "80"
	}

	return ColorFromString(e).OrElse(Color("80")).String()
}

func border() string {
	e := os.Getenv("ARROW_BORDER_COLOR")
	if e == "" {
		return "80"
	}

	return ColorFromString(e).OrElse(Color("80")).String()
}

func ColorFromString(s string) mo.Option[Color] {
	i, err := strconv.Atoi(s)

	if err == nil && i >= 0 && i <= 255 {
		return mo.Some(Color(s))
	}

	if rgbPattern.MatchString(s) {
		return mo.Some(Color(s))
	}

	return mo.None[Color]()
}
