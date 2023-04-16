package dbgo

// Copyright (C) Philip Schlump, 2014.

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"
)

/*
	"github.com/pschlump/gc
TODO
	1. Add in JSON format
	2. Add in %(LINE) %(FUNC) etc.
	3. Add in %(TB:4) 	Trace back 4 levels


func x() {
	var buffer bytes.Buffer

	for i := 0; i < 1000; i++ {
		buffer.WriteString("a")
	}

	fmt.Println(buffer.String())
}
*/

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

var dbOn map[string]bool
var dbLock = sync.RWMutex{}

func init() {
	dbOn = make(map[string]bool)
	// dbLock = sync.RwMutex{}
}

func SetDbFlag(f map[string]bool) {
	for k, v := range f {
		dbOn[k] = v
	}
}

func ProcessFormat(format string, a []interface{}) (rv string, params []interface{}) {
	var buffer bytes.Buffer
	colorFound := false
	params = make([]interface{}, 0, len(a))
	param_no := 0
	// fmt.Printf("len(a) = %d, %s\n", len(a), LF())
	var i, j int
	for i = 0; i < len(format); i++ {
		if format[i] == '%' && i+1 < len(format) && format[i+1] == '(' {
			color := "Red"
			for j = i + 2; j < len(format); j++ {
				if format[j] == ')' && i+2 < j {
					// fmt.Printf("Pick of %d:%d\n", i+2, j)
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
			case "J", "json-indent":
				ct = SVar(a[param_no])
				param_no++
			case "j", "json":
				ct = SVarI(a[param_no])
				param_no++
			default:
				// fmt.Printf("--->%s<---, %s\n", color, LF())
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
			// fmt.Printf("i = %d, %s\n", i, LF())
			buffer.WriteByte(format[i])
			// must parse better ! ! !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
			if format[i] == '%' {
				if param_no < len(a) {
					params = append(params, a[param_no])
				} else {
					if i+1 < len(format) {
						params = append(params, fmt.Sprintf("--- Invalid Missing Value, format ->%%%c<- position %d ---", format[i+1], param_no))
					} else {
						params = append(params, fmt.Sprintf("--- Invalid Missing Value, format ->%%%c<- position %d ---", '?', param_no))
					}
				}
				param_no++
			}
		}
	}
	if colorFound {
		buffer.WriteString(ColorTab["Reset"])
	}
	return buffer.String(), params
}

func Printf(format string, a ...interface{}) (n int, err error) {
	ff, b := ProcessFormat(format, a)
	return fmt.Printf(ff, b...)
}

func Sprintf(format string, a ...interface{}) (ss string) {
	ff, b := ProcessFormat(format, a)
	ss = fmt.Sprintf(ff, b...)
	return
}

func Fprintf(w io.Writer, format string, a ...interface{}) (n int, err error) {
	ff, b := ProcessFormat(format, a)
	return fmt.Fprintf(w, ff, b...)
}

func DbPf(db bool, format string, a ...interface{}) (n int, err error) {
	if db {
		ff, b := ProcessFormat(format, a)
		return fmt.Printf(ff, b...)
	}
	return 0, nil
}

func DbFpf(db bool, w io.Writer, format string, a ...interface{}) (n int, err error) {
	if db {
		ff, b := ProcessFormat(format, a)
		return fmt.Fprintf(w, ff, b...)
	}
	return 0, nil
}

func DbPfb(db bool, format string, a ...interface{}) (n int, err error) {
	if db {
		ff, b := ProcessFormat(format, a)
		fmt.Fprintf(os.Stderr, ff, b...)
		return fmt.Printf(ff, b...)
	}
	return 0, nil
}

func IsDbOn(dbflag string) bool {
	dbLock.RLock()
	if dbOn[dbflag] {
		dbLock.RUnlock()
		return true
	}
	dbLock.RUnlock()
	return false
}

func DbPrintf(dbflag string, format string, a ...interface{}) (n int, err error) {
	dbLock.RLock()
	if dbOn[dbflag] {
		dbLock.RUnlock()
		ff, b := ProcessFormat(format, a)
		return fmt.Printf(ff, b...)
	}
	dbLock.RUnlock()
	return 0, nil
}

func DbFprintf(dbflag string, w io.Writer, format string, a ...interface{}) (n int, err error) {
	dbLock.RLock()
	if dbOn[dbflag] {
		dbLock.RUnlock()
		ff, b := ProcessFormat(format, a)
		return fmt.Fprintf(w, ff, b...)
	}
	dbLock.RUnlock()
	return 0, nil
}

func DbPfe(envVar string, format string, a ...interface{}) (n int, err error) {
	return DbPfb(ChkEnv(envVar), format, a...)
}

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
