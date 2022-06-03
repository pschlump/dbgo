package dbgo

import (
	"os"

	"github.com/pschlump/ansi"
)

/*
```go
Color(s, "red")            // red
Color(s, "red+b")          // red bold
Color(s, "red+B")          // red blinking
Color(s, "red+u")          // red underline
Color(s, "red+bh")         // red bold bright
Color(s, "red:white")      // red on white
Color(s, "red+b:white+h")  // red bold on white bright
Color(s, "red+B:white+h")  // red blink on white bright
Color(s, "off")            // turn off ansi codes
```
Colors

* black
* red
* green
* yellow
* blue
* magenta
* cyan
* white
*/

// Red
var ColorRed string

// Yellow
var ColorYellow string

// Green
var ColorGreen string

// Blue
var ColorBlue string

// Black
var ColorBlack string

// Magenta
var ColorMagenta string

// Cyan
var ColorCyan string

var ColorBlueOnWhite string
var ColorWhiteOnBlue string
var ColorMagentaOnWhite string
var ColorGreenOnWhite string

// Reset to default color
var ColorReset string

func init() {
	ColorRed = ""
	ColorYellow = ""
	ColorGreen = ""
	ColorBlue = ""
	ColorBlack = ""
	ColorMagenta = ""
	ColorCyan = ""
	ColorBlueOnWhite = ""
	ColorWhiteOnBlue = ""
	ColorMagentaOnWhite = ""
	ColorGreenOnWhite = ""
	ColorReset = ""
	// if !StdErrPiped() { // check if stderr is terminal
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
	// }
}

// InArray returns the subscript if the string 's' is found in the array 'arr', else -1
func InArray(s string, arr []string) int {
	for i, v := range arr {
		if v == s {
			return i
		}
	}
	return -1
}

// StdErrPiped returns true if os.Stderr apears to be send to a pipe
func StdErrPiped() bool {
	fi, _ := os.Stderr.Stat() // get the FileInfo struct describing the standard input.

	if (fi.Mode() & os.ModeCharDevice) == 0 {
		return true // output is piped to file
	}

	return false
}

// StdErrPiped returns true if os.Stdout apears to be send to a pipe
func StdOutPiped() bool {
	fi, _ := os.Stdout.Stat() // get the FileInfo struct describing the standard input.

	if (fi.Mode() & os.ModeCharDevice) == 0 {
		return true // output is piped to file
	}

	return false
}

// StdErrPiped returns true if os.Stdin apears to be send to a pipe
func StdInPiped() bool {
	fi, _ := os.Stdin.Stat() // get the FileInfo struct describing the standard input.

	if (fi.Mode() & os.ModeCharDevice) == 0 {
		return true // Data from pipe
	}

	return false
}
