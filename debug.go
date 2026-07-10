package dbgo

// ----------------------------------------------------------------------------------------------------------
//
// Simple functions to help with debugging Go (golang) code.
//
// Copyright (C) Philip Schlump, 2013-2024.
// See LICENSE file for details. -- Same as Go source code.
//
// The location helpers (LF, LINE, FILE, FUNCNAME, ...) report where they were
// called from via the runtime call stack. By convention their optional "depth"
// argument is 1 for the immediate caller (the default), 2 for the caller's
// caller, and so on.
//
// ----------------------------------------------------------------------------------------------------------

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/pschlump/json" // modified from "encoding/json" to handle undefined types by ignoring them.
)

// callerDepth resolves the optional call-stack depth argument used by the
// runtime.Caller-based helpers. With no argument it defaults to 1, meaning
// "the caller of the helper". 0 is the helper itself, 2 the caller's caller.
func callerDepth(d []int) int {
	if len(d) > 0 {
		return d[0]
	}
	return 1
}

// LINE returns the current line number as a string. The optional depth
// argument selects how far up the call stack to report: 0 is LINE itself,
// 1 (the default) is the caller of LINE, 2 the caller's caller, and so on.
// Extra arguments are ignored. "LineNo:Unk" is returned if the frame cannot
// be recovered.
func LINE(d ...int) string {
	_, _, line, ok := runtime.Caller(callerDepth(d))
	if ok {
		return fmt.Sprintf("%d", line)
	}
	return "LineNo:Unk"
}

// LINEnf returns the line number and source file at the given stack depth.
// The line number is -1 (and the file "") if the frame cannot be recovered.
func LINEnf(d ...int) (int, string) {
	_, file, line, ok := runtime.Caller(callerDepth(d))
	if ok {
		return line, file
	}
	return -1, ""
}

// FILE returns the source file name at the given stack depth, or "File:Unk"
// if the frame cannot be recovered.
func FILE(d ...int) string {
	_, file, _, ok := runtime.Caller(callerDepth(d))
	if ok {
		return file
	}
	return "File:Unk"
}

// LF returns a "File: <file> LineNo:<line>" string for the caller. It is the
// most commonly used of the location helpers.
//
// The optional depth argument selects the stack frame:
//   - 0: the LF function itself.
//   - 1 (the default): the caller of LF.
//   - 2, 3, ...: further callers up the stack.
//
// A negative depth switches on a "walk back across files" mode: frames are
// scanned forward until the source file changes, and every line number in the
// run of frames sharing the starting file is reported. A depth of -1 reports
// the current file's run, -2 additionally reports the next file's run, and so
// on. This is useful when several wrapper functions live in the same file and
// you want to see all of their line numbers at once.
func LF(d ...int) string {
	depth := callerDepth(d)
	loop := false
	nf := 0
	if depth <= -1 { // negative: walk back across this many distinct files
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
			ss += fmt.Sprintf("File: %s LineNo:%d ", file0, ln)
			file0 = file
		}
		return ss
	}
	_, file, line, ok := runtime.Caller(depth)
	if ok {
		return fmt.Sprintf("File: %s LineNo:%d", file, line)
	}
	return "File: Unk LineNo:Unk"
}

// LFj returns the file name and line number formatted as a JSON object
// fragment, for example `"File": "/path/x.go", "LineNo":42`. This is handy
// when embedding location information inside a JSON log line. It returns ""
// if the frame cannot be recovered.
func LFj(d ...int) string {
	_, file, line, ok := runtime.Caller(callerDepth(d))
	if !ok {
		return ""
	}
	return fmt.Sprintf("\"File\": \"%s\", \"LineNo\":%d", file, line)
}

// FUNCNAME returns the name of the function at the given stack depth, or
// "FunctionName:Unk" if the frame cannot be recovered.
func FUNCNAME(d ...int) string {
	pc, _, _, ok := runtime.Caller(callerDepth(d))
	if !ok {
		return "FunctionName:Unk"
	}
	return runtime.FuncForPC(pc).Name()
}

// IAmAt prints the current function, file and line number (one frame up),
// followed by the optional message strings, to standard output. It is a
// quick "where am I" trace marker.
func IAmAt(s ...string) {
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		fmt.Printf("Func:Unk File:Unk LineNo:Unk, %s\n", strings.Join(s, " "))
		return
	}
	xfunc := runtime.FuncForPC(pc).Name()
	fmt.Printf("Func:%s File:%s LineNo:%d, %s\n", xfunc, file, line, strings.Join(s, " "))
}

// IAmAt2 is like IAmAt but also reports the caller's caller ("called...")
// before the current frame, making a two-level trace.
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

// SVar returns v as a compact JSON string. If marshaling fails it returns
// "Error:" followed by the error message instead of panicking.
func SVar(v interface{}) string {
	s, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprintf("Error:%s", err)
	}
	return string(s)
}

// SVarI returns v as a tab-indented JSON string. If marshaling fails it
// returns "Error:" followed by the error message instead of panicking.
func SVarI(v interface{}) string {
	s, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return fmt.Sprintf("Error:%s", err)
	}
	return string(s)
}

// indexOf returns the index of the first element of s equal to v, or -1 if v
// is not present. It is the shared implementation behind the InArray* helpers.
func indexOf[T comparable](v T, s []T) int {
	for i, x := range s {
		if x == v {
			return i
		}
	}
	return -1
}

// InArrayString returns the index of the first occurrence of s in arr, or -1
// if s is not present.
func InArrayString(s string, arr []string) int {
	return indexOf(s, arr)
}

// InArrayInt returns the index of the first occurrence of n in arr, or -1 if
// n is not present.
func InArrayInt(n int, arr []int) int {
	return indexOf(n, arr)
}

// stdoutOn gates TrIAmAt. Nothing in this package sets it to true, so
// TrIAmAt is inert unless a caller (from within the package) enables it.
var stdoutOn = false

// TrIAmAt is a "trace" variant of IAmAt that only prints when stdoutOn is
// true. Because nothing in this package enables stdoutOn, TrIAmAt is inert by
// default; set stdoutOn = true to turn it on.
func TrIAmAt(s ...string) {
	if !stdoutOn {
		return
	}
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		fmt.Printf("Func:Unk File:Unk LineNo:Unk, %s\n", strings.Join(s, " "))
		return
	}
	xfunc := runtime.FuncForPC(pc).Name()
	fmt.Printf("Func:%s File:%s LineNo:%d, %s\n", xfunc, file, line, strings.Join(s, " "))
}

// Db2Printf behaves like fmt.Fprintf to standard output, but only when flag
// is true; otherwise it does nothing.
func Db2Printf(flag bool, format string, a ...interface{}) (n int, err error) {
	if flag {
		return fmt.Fprintf(os.Stdout, format, a...)
	}
	return
}

// LF2 returns the line number and file name (as separate values, unlike LF)
// for the caller. line is 0 and file is "Unk" if the frame cannot be
// recovered.
func LF2(d ...int) (line int, file string) {
	var ok bool
	_, file, line, ok = runtime.Caller(callerDepth(d))
	if !ok {
		line = 0
		file = "Unk"
	}
	return
}

/* vim: set noai ts=4 sw=4: */
