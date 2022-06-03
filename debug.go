package dbgo

// ----------------------------------------------------------------------------------------------------------
//
// Simple functions to help with debugging Go (golang) code.
//
// Copyright (C) Philip Schlump, 2013-2017.
// Version: 1.0.2
// See LICENSE file for details. -- Same as Go source code.
// BuildNo: 063
//
// I usually use these like this:
//
//     func something ( j int ) {
//			...
//			fmt.Pritnf ( "Ya someting useful %s\n", debug.LF(1) )
//
// This prints out the line and file that "Ya..." is at - so that it is easier for me to match output
// with code.   The "depth" == 1 parameter is how far up the stack I want to go.  0 is the LF routine.
// 1 is the caller of LF, usually what I want and the default, 2 is the caller of "something".
//
// The most useful functions are:
//    LF 			Return as a string the line number and file name.
//	  IAmAt			Print out current line/file
//	  SVarI			Convert most things to an indented JSON string and return it.
//
// To import put this in your code:
//
//		import (
//			"git.q8s.co/pschlump/dbgo"
//		)
//
// Then
//
//		fmt.Printf ( ".... %s ...\n", dbgo.LF() )
//
// ----------------------------------------------------------------------------------------------------------

import (
	"fmt"
	"github.com/pschlump/json" // modified from "encoding/json" to handle undefined types by ignoring them.
	"os"
	"runtime"
	"strings"
)

// LINE Return the current line number as a string.  Default parameter is 1, must be an integer
// That reflects the depth in the call stack.  A value of 0 would be the LINE() function
// itself.  If you supply more than one parameter, 2..n are ignored.
func LINE(d ...int) string {
	depth := 1
	if len(d) > 0 {
		depth = d[0]
	}
	_, _, line, ok := runtime.Caller(depth)
	if ok {
		return fmt.Sprintf("%d", line)
	}
	return "LineNo:Unk"
}

// LINEnf Returns line number, 0 if error
func LINEnf(d ...int) (int, string) {
	depth := 1
	if len(d) > 0 {
		depth = d[0]
	}
	_, file, line, ok := runtime.Caller(depth)
	if ok {
		return line, file
	}
	return -1, ""
}

// FILE Returns the current file name.
func FILE(d ...int) string {
	depth := 1
	if len(d) > 0 {
		depth = d[0]
	}
	_, file, _, ok := runtime.Caller(depth)
	if ok {
		return file
	} else {
		return "File:Unk"
	}
}

// LF Returns the File name and Line no as a string.
func LF(d ...int) string {
	depth := 1
	if len(d) > 0 {
		depth = d[0]
	}
	loop := false
	nf := 0
	if depth <= -1 { // if <= -1, then number of files to walk back.
		nf = (-depth) + 1
		depth = 1
		loop = true
	}
	if loop {
		_, file0, line, ok := runtime.Caller(depth)
		ss := ""
		for ii := 0; ii < nf; ii++ {
			file := file0
			ln := make([]int, 0, depth)
			for ok && file == file0 {
				ln = append(ln, line)
				depth++
				_, file, line, ok = runtime.Caller(depth)
			}
			ss = ss + fmt.Sprintf("File: %s LineNo:%d ", file0, ln)
			file0 = file
		}
		return ss
	} else {
		_, file, line, ok := runtime.Caller(depth)
		if ok {
			return fmt.Sprintf("File: %s LineNo:%d", file, line)
		} else {
			return fmt.Sprintf("File: Unk LineNo:Unk")
		}
	}
}

// LFj returns the File name and Line no as a string. - for JSON as string
func LFj(d ...int) string {
	depth := 1
	if len(d) > 0 {
		depth = d[0]
	}
	_, file, line, ok := runtime.Caller(depth)
	if ok {
		return fmt.Sprintf("\"File\": \"%s\", \"LineNo\":%d", file, line)
	} else {
		return ""
	}
}

// FUNCNAME returns the current function name as a string.
func FUNCNAME(d ...int) string {
	depth := 1
	if len(d) > 0 {
		depth = d[0]
	}
	pc, _, _, ok := runtime.Caller(depth)
	if ok {
		xfunc := runtime.FuncForPC(pc).Name()
		return xfunc
	} else {
		return fmt.Sprintf("FunctionName:Unk")
	}
}

// IAmAt print out the current Function,File,Line No and an optional set of strings.
func IAmAt(s ...string) {
	pc, file, line, ok := runtime.Caller(1)
	if ok {
		xfunc := runtime.FuncForPC(pc).Name()
		fmt.Printf("Func:%s File:%s LineNo:%d, %s\n", xfunc, file, line, strings.Join(s, " "))
	} else {
		fmt.Printf("Func:Unk File:Unk LineNo:Unk, %s\n", strings.Join(s, " "))
	}
}

// IAmAt2 prints out the current Function,File,Line No and an optional set of strings - do this for 2 levels deep.
func IAmAt2(s ...string) {
	pc, file, line, ok := runtime.Caller(1)
	pc2, file2, line2, ok2 := runtime.Caller(2)
	if ok {
		xfunc := runtime.FuncForPC(pc).Name()
		if ok2 {
			xfunc2 := runtime.FuncForPC(pc2).Name()
			fmt.Printf("Func:%s File: %s LineNo:%d, called...\n", xfunc2, file2, line2)
		} else {
			fmt.Printf("Func:Unk File: unk LineNo:unk, called...\n")
		}
		fmt.Printf("Func:%s File: %s LineNo:%d, %s\n", xfunc, file, line, strings.Join(s, " "))
	} else {
		fmt.Printf("Func:Unk File: Unk LineNo:Unk, %s\n", strings.Join(s, " "))
	}
}

// SVar return the JSON encoded version of the data.
func SVar(v interface{}) string {
	s, err := json.Marshal(v)
	// s, err := json.MarshalIndent ( v, "", "\t" )
	if err != nil {
		return fmt.Sprintf("Error:%s", err)
	} else {
		return string(s)
	}
}

// SVarI return the JSON encoded version of the data with tab indentation.
func SVarI(v interface{}) string {
	// s, err := json.Marshal ( v )
	s, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return fmt.Sprintf("Error:%s", err)
	} else {
		return string(s)
	}
}

// Return 0..n if 's' is in the array arr, else -1.
func InArrayString(s string, arr []string) int {
	for i, v := range arr {
		if v == s {
			return i
		}
	}
	return -1
}

// Return 0..n if 'n' is in the array arr, else -1.
func InArrayInt(s int, arr []int) int {
	for i, v := range arr {
		if v == s {
			return i
		}
	}
	return -1
}

var stdout_on = false

// Hm...
func TrIAmAt(s ...string) {
	if stdout_on {
		pc, file, line, ok := runtime.Caller(1)
		if ok {
			xfunc := runtime.FuncForPC(pc).Name()
			fmt.Printf("Func:%s File:%s LineNo:%d, %s\n", xfunc, file, line, strings.Join(s, " "))
		} else {
			fmt.Printf("Func:Unk File:Unk LineNo:Unk, %s\n", strings.Join(s, " "))
		}
	}
}

// Printf with a true false flag.
func Db2Printf(flag bool, format string, a ...interface{}) (n int, err error) {
	if flag {
		return fmt.Fprintf(os.Stdout, format, a...)
	}
	return
}

// LF2 returns the line/file for the parent.
func LF2(d ...int) (line int, file string) {
	depth := 1
	if len(d) > 0 {
		depth = d[0]
	}
	var ok bool
	_, file, line, ok = runtime.Caller(depth)
	if !ok {
		line = 0
		file = "Unk"
	}
	return
}

/* vim: set noai ts=4 sw=4: */
