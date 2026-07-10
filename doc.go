// Package dbgo is a small toolkit of helpers for debugging Go programs.
//
// It groups together the things that are repeatedly useful when instrumenting
// code: reporting where a log line came from, pretty-printing values as JSON,
// gating debug output behind flags or environment variables, and coloring
// terminal output.
//
// # Where-am-I helpers
//
// The location helpers report their caller via the runtime call stack:
//
//	dbgo.LF()         // "File: /path/foo.go LineNo:42"
//	dbgo.LINE()       // "42"
//	dbgo.FILE()       // "/path/foo.go"
//	dbgo.FUNCNAME()   // "main.handleRequest"
//
// Each takes an optional depth argument (default 1 = the caller of the helper,
// 2 = the caller's caller, and so on). See [LF] for the negative-depth
// "walk back across files" mode.
//
// # JSON helpers
//
// [SVar] and [SVarI] marshal a value to compact and tab-indented JSON
// respectively, returning an "Error:..." string instead of panicking on
// failure:
//
//	dbgo.SVar([]int{1, 2, 3})   // "[1,2,3]"
//
// # Extended printf
//
// [Printf], [Sprintf] and [Fprintf] behave like their fmt counterparts but
// understand extra %(name) directives in the format string:
//
//	dbgo.Printf("at %(LF), data=%(j)\n", someValue)
//
//	%(LF)                caller's file and line
//	%(j) / %(json-indent) next argument as indented JSON
//	%(J) / %(json)        next argument as compact JSON
//	%(Red), %(Green), ... ANSI color (see [ColorTab]); %(!) / %(Reset) resets
//
// Missing arguments produce a visible placeholder rather than an out-of-bounds
// panic.
//
// # Gated debug output
//
// [DbPrintf] and [DbFprintf] print only when a named debug flag has been
// enabled with [SetDbFlag]; query a flag with [IsDbOn]. [DbPfe] gates on an
// environment variable via [ChkEnv].
//
// [ChkEnv] reports whether an environment variable parses as true (see
// [ParseBool]) and caches the result so the underlying os.Getenv is called at
// most once per name. [SetFlag] overrides the environment programmatically.
//
// # Small utilities
//
// [InArray], [InArrayString] and [InArrayInt] search a slice for a value.
// [StdErrPiped], [StdOutPiped] and [StdInPiped] report whether a standard
// stream is redirected to a pipe/file. [BackTickQuote] backtick-quotes a
// string for shell use.
//
// The LogData/LogToFile types in this package are experimental and currently
// no-ops; see the warning in log.go.
package dbgo
