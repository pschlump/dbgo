package dbgo

// Copyright (C) Philip Schlump, 2014-2024.

import (
	"os"

	"github.com/pschlump/ansi"
)

/*
Color escapes can also be produced directly via the ansi package, for example:

	Color(s, "red")            // red
	Color(s, "red+b")          // red bold
	Color(s, "red+B")          // red blinking
	Color(s, "red+u")          // red underline
	Color(s, "red+bh")         // red bold bright
	Color(s, "red:white")      // red on white
	Color(s, "red+b:white+h")  // red bold on white bright
	Color(s, "red+B:white+h")  // red blink on white bright
	Color(s, "off")            // turn off ansi codes

The named colors are: black, red, green, yellow, blue, magenta, cyan, white.
*/

// ColorRed is the ANSI escape sequence for red foreground on black background.
var ColorRed string

// ColorYellow is the ANSI escape sequence for yellow foreground on black background.
var ColorYellow string

// ColorGreen is the ANSI escape sequence for green foreground on black background.
var ColorGreen string

// ColorBlue is the ANSI escape sequence for blue foreground on black background.
var ColorBlue string

// ColorBlack is the ANSI escape sequence for black foreground on white background.
var ColorBlack string

// ColorMagenta is the ANSI escape sequence for magenta foreground on black background.
var ColorMagenta string

// ColorCyan is the ANSI escape sequence for cyan foreground on black background.
var ColorCyan string

// ColorBlueOnWhite is the ANSI escape sequence for blue foreground on white background.
var ColorBlueOnWhite string

// ColorWhiteOnBlue is the ANSI escape sequence for white foreground on blue background.
var ColorWhiteOnBlue string

// ColorMagentaOnWhite is the ANSI escape sequence for magenta foreground on white background.
var ColorMagentaOnWhite string

// ColorGreenOnWhite is the ANSI escape sequence for bold green foreground on white background.
var ColorGreenOnWhite string

// ColorReset is the ANSI escape sequence that resets the terminal to its default colors.
var ColorReset string

func init() {
	// Colors are always populated. Callers that write to a non-terminal
	// destination and want to suppress escapes can check StdOutPiped first.
	ColorRed = ansi.ColorCode("red:black")
	ColorYellow = ansi.ColorCode("yellow:black")
	ColorGreen = ansi.ColorCode("green:black")
	ColorBlue = ansi.ColorCode("blue:black")
	ColorBlack = ansi.ColorCode("black:white")
	ColorMagenta = ansi.ColorCode("magenta:black")
	ColorCyan = ansi.ColorCode("cyan:black")
	ColorBlueOnWhite = ansi.ColorCode("blue:white")
	ColorWhiteOnBlue = ansi.ColorCode("white:blue")
	ColorMagentaOnWhite = ansi.ColorCode("magenta:white")
	ColorGreenOnWhite = ansi.ColorCode("green+b:white")
	ColorReset = ansi.ColorCode("reset")
}

// InArray returns the index of the first occurrence of s in arr, or -1 if s
// is not present. It is equivalent to InArrayString.
func InArray(s string, arr []string) int {
	return indexOf(s, arr)
}

// StdErrPiped reports whether os.Stderr appears to be redirected to a pipe or
// file (i.e. it is not a character device / terminal).
func StdErrPiped() bool {
	fi, _ := os.Stderr.Stat()
	if (fi.Mode() & os.ModeCharDevice) == 0 {
		return true // output is piped to a file
	}
	return false
}

// StdOutPiped reports whether os.Stdout appears to be redirected to a pipe or
// file (i.e. it is not a character device / terminal).
func StdOutPiped() bool {
	fi, _ := os.Stdout.Stat()
	if (fi.Mode() & os.ModeCharDevice) == 0 {
		return true // output is piped to a file
	}
	return false
}

// StdInPiped reports whether os.Stdin appears to receive data from a pipe or
// file (i.e. it is not a character device / terminal).
func StdInPiped() bool {
	fi, _ := os.Stdin.Stat()
	if (fi.Mode() & os.ModeCharDevice) == 0 {
		return true // data from a pipe
	}
	return false
}
