package dbgo

// Copyright (C) Philip Schlump, 2014-2024.
//
// This file implements an extended printf. On top of the normal fmt verbs,
// the format string may contain %(name) directives:
//
//   - %(LF)              the caller's "File: <path> LineNo:<n>" location.
//   - %(j) / %(json-indent)   the next argument as indented JSON.
//   - %(J) / %(json)          the next argument as compact JSON.
//   - %(Red), %(Green), ...   an ANSI color; see ColorTab for the full set.
//   - %(!) / %(Reset)         reset color back to the terminal default.
//
// Unknown %(name) directives default to red. A trailing color is automatically
// reset at the end of the output.

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"
)

// ColorTab maps color names (used in %(name) directives) to their ANSI escape
// sequences. It is populated in init from the Color* variables.
var ColorTab map[string]string

func init() {
	ColorTab = make(map[string]string)

	ColorTab["red"] = ColorRed
	ColorTab["green"] = ColorGreen
	ColorTab["yellow"] = ColorYellow
	ColorTab["blue"] = ColorBlueOnWhite
	ColorTab["cyan"] = ColorCyan
	ColorTab["magenta"] = ColorMagentaOnWhite
	ColorTab["green_on_white"] = ColorGreenOnWhite
	ColorTab["greenw"] = ColorGreenOnWhite
	ColorTab["reset"] = ColorReset

	ColorTab["Red"] = ColorRed
	ColorTab["Green"] = ColorGreen
	ColorTab["Yellow"] = ColorYellow
	ColorTab["Blue"] = ColorBlueOnWhite
	ColorTab["Cyan"] = ColorCyan
	ColorTab["Reset"] = ColorReset
	ColorTab["Magenta"] = ColorMagentaOnWhite
	ColorTab["Green_on_white"] = ColorGreenOnWhite
	ColorTab["GreenOnWhite"] = ColorGreenOnWhite
	ColorTab["GreenW"] = ColorGreenOnWhite

	ColorTab["!"] = ColorReset
}

// dbOn holds the set of currently enabled debug flags, consulted by DbPrintf,
// DbFprintf and IsDbOn. All access is guarded by dbLock.
var dbOn map[string]bool

// dbLock protects dbOn.
var dbLock = sync.RWMutex{}

func init() {
	dbOn = make(map[string]bool)
}

// SetDbFlag enables or disables the named debug flags. It is safe for
// concurrent use with DbPrintf, DbFprintf and IsDbOn.
func SetDbFlag(f map[string]bool) {
	dbLock.Lock()
	defer dbLock.Unlock()
	for k, v := range f {
		dbOn[k] = v
	}
}

// ProcessFormat expands the dbgo %(name) directives in format and returns the
// resulting format string together with the slice of arguments that should be
// passed to a fmt.*printf call. Directives that consume an argument (%(j) and
// %(J)) pull from a in order; any remaining ordinary %v/%s style verbs consume
// the following arguments. Missing arguments are replaced with a placeholder
// marker instead of causing an out-of-bounds panic.
func ProcessFormat(format string, a []interface{}) (rv string, params []interface{}) {
	var buffer bytes.Buffer
	colorFound := false
	params = make([]interface{}, 0, len(a))
	param_no := 0
	var i, j int
	for i = 0; i < len(format); i++ {
		if format[i] == '%' && i+1 < len(format) && format[i+1] == '(' {
			color := "Red"
			for j = i + 2; j < len(format); j++ {
				if format[j] == ')' && i+2 < j {
					color = format[i+2 : j]
					break
				}
			}
			i = j
			var ct string
			var ok bool
			switch color {
			case "LF":
				ct = LF(3)
			case "J", "json": // compact JSON
				if param_no < len(a) {
					ct = SVar(a[param_no])
					param_no++
				} else {
					ct = "--- Invalid Missing Value for %(J), position " + fmt.Sprintf("%d", param_no) + " ---"
				}
			case "j", "json-indent": // indented JSON
				if param_no < len(a) {
					ct = SVarI(a[param_no])
					param_no++
				} else {
					ct = "--- Invalid Missing Value for %(j), position " + fmt.Sprintf("%d", param_no) + " ---"
				}
			default:
				ct, ok = ColorTab[color]
				if !ok {
					ct = ColorTab["Red"]
				}
				colorFound = true
			}
			buffer.WriteString(ct)

			if color == "!" || color == "Reset" {
				colorFound = false
			}
		} else {
			if format[i] == '%' {
				if i+1 < len(format) && format[i+1] == '%' {
					// literal percent: emit one '%' and skip the next byte, consuming no argument.
					buffer.WriteString("%")
					i++
				} else if param_no < len(a) {
					buffer.WriteByte('%')
					params = append(params, a[param_no])
					param_no++
				} else {
					verb := byte('?')
					if i+1 < len(format) {
						verb = format[i+1]
					}
					buffer.WriteByte('%')
					params = append(params, fmt.Sprintf("--- Invalid Missing Value, format ->%%%c<- position %d ---", verb, param_no))
					param_no++
				}
			} else {
				buffer.WriteByte(format[i])
			}
		}
	}
	if colorFound {
		buffer.WriteString(ColorTab["Reset"])
	}
	return buffer.String(), params
}

// Printf is like fmt.Printf but expands the dbgo %(name) directives described
// at the top of this file.
func Printf(format string, a ...interface{}) (n int, err error) {
	ff, b := ProcessFormat(format, a)
	return fmt.Printf(ff, b...)
}

// Sprintf is like fmt.Sprintf but expands the dbgo %(name) directives.
func Sprintf(format string, a ...interface{}) (ss string) {
	ff, b := ProcessFormat(format, a)
	ss = fmt.Sprintf(ff, b...)
	return
}

// Fprintf is like fmt.Fprintf but expands the dbgo %(name) directives.
func Fprintf(w io.Writer, format string, a ...interface{}) (n int, err error) {
	ff, b := ProcessFormat(format, a)
	return fmt.Fprintf(w, ff, b...)
}

// DbPf behaves like Printf but only when db is true.
func DbPf(db bool, format string, a ...interface{}) (n int, err error) {
	if db {
		ff, b := ProcessFormat(format, a)
		return fmt.Printf(ff, b...)
	}
	return 0, nil
}

// DbFpf behaves like Fprintf (to w) but only when db is true.
func DbFpf(db bool, w io.Writer, format string, a ...interface{}) (n int, err error) {
	if db {
		ff, b := ProcessFormat(format, a)
		return fmt.Fprintf(w, ff, b...)
	}
	return 0, nil
}

// DbPfb behaves like Printf but, when db is true, writes the formatted output
// to both standard output and standard error ("print both").
func DbPfb(db bool, format string, a ...interface{}) (n int, err error) {
	if db {
		ff, b := ProcessFormat(format, a)
		fmt.Fprintf(os.Stderr, ff, b...)
		return fmt.Printf(ff, b...)
	}
	return 0, nil
}

// IsDbOn reports whether the named debug flag is currently enabled.
func IsDbOn(dbflag string) bool {
	dbLock.RLock()
	on := dbOn[dbflag]
	dbLock.RUnlock()
	return on
}

// DbPrintf behaves like Printf but only when the named debug flag has been
// enabled via SetDbFlag.
func DbPrintf(dbflag string, format string, a ...interface{}) (n int, err error) {
	if !IsDbOn(dbflag) {
		return 0, nil
	}
	ff, b := ProcessFormat(format, a)
	return fmt.Printf(ff, b...)
}

// DbFprintf behaves like Fprintf (to w) but only when the named debug flag has
// been enabled via SetDbFlag.
func DbFprintf(dbflag string, w io.Writer, format string, a ...interface{}) (n int, err error) {
	if !IsDbOn(dbflag) {
		return 0, nil
	}
	ff, b := ProcessFormat(format, a)
	return fmt.Fprintf(w, ff, b...)
}

// DbPfe prints to both standard output and standard error (see DbPfb) when the
// named environment variable parses as a true value (see ChkEnv/ParseBool).
func DbPfe(envVar string, format string, a ...interface{}) (n int, err error) {
	return DbPfb(ChkEnv(envVar), format, a...)
}

// BackTickQuote returns ss wrapped in backticks, with any embedded backticks
// escaped as the hex sequence \x60 so the result is safe to use inside a
// shell-style backtick-quoted context.
func BackTickQuote(ss string) string {
	var buffer bytes.Buffer

	buffer.WriteByte('`')
	for i, c := range []byte(ss) {
		if c == '`' {
			buffer.WriteString("\\x60")
		} else {
			buffer.WriteByte(ss[i])
		}
	}
	buffer.WriteByte('`')

	return buffer.String()
}
