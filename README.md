# dbgo

[![license](https://img.shields.io/badge/license-BSD--3-blue.svg?style=flat)](https://github.com/pschlump/dbgo/blob/main/LICENSE) [![Go Reference](https://pkg.go.dev/badge/github.com/pschlump/dbgo.svg)](https://pkg.go.dev/github.com/pschlump/dbgo)

A small library of helpers for debugging Go programs: report where a log line
came from, pretty-print values as JSON, gate debug output behind flags or
environment variables, and color terminal output.

## Install

```sh
go get github.com/pschlump/dbgo
```

```go
import "github.com/pschlump/dbgo"
```

## Quick start

```go
package main

import (
	"fmt"

	"github.com/pschlump/dbgo"
)

func main() {
	// "I am at: File: /path/main.go LineNo:10"
	fmt.Printf("I am at: %s\n", dbgo.LF())

	// Extended printf: %(LF) location, %(j) indented JSON, %(Green) color.
	dbgo.Printf("at %(LF), data=%(Green)%(j)%(!)\n", []string{"a", "b"})

	// Debug output gated behind a flag.
	dbgo.SetDbFlag(map[string]bool{"network": true})
	dbgo.DbPrintf("network", "sending %(J)\n", 42)
}
```

## Where-am-I helpers

The location helpers report their caller via the runtime call stack. Each takes
an optional `depth` argument: `1` (the default) is the caller of the helper,
`2` is the caller's caller, `0` is the helper itself.

| Function     | Returns                                         |
|--------------|-------------------------------------------------|
| `LF(depth)`  | `"File: /path/foo.go LineNo:42"`                |
| `LINE(depth)`| `"42"`                                          |
| `FILE(depth)`| `"/path/foo.go"`                                |
| `FUNCNAME(depth)` | `"main.handleRequest"`                     |
| `LF2(depth)` | `(line int, file string)`                       |
| `LFj(depth)` | `"File": "/path/foo.go", "LineNo":42` (JSON fragment) |
| `IAmAt(s...)`| prints `Func:… File:… LineNo:…, <s>` to stdout |

`LF` also accepts a **negative** depth, which switches on a "walk back across
files" mode: it reports every line number in the run of stack frames that share
the starting source file. `-1` reports the current file's run, `-2 the current
file's run plus the next file's, and so on. See the package docs for details.

## JSON helpers

`SVar` and `SVarI` marshal a value to compact and tab-indented JSON. On error
they return an `"Error:…"` string rather than panicking.

```go
dbgo.SVar([]int{1, 2, 3})  // "[1,2,3]"
dbgo.SVarI([]int{1, 2, 3}) // "[\n\t1,\n\t2,\n\t3\n]"
```

## Extended printf

`Printf`, `Sprintf` and `Fprintf` behave like their `fmt` counterparts but
understand extra `%(name)` directives:

| Directive                | Description                                              |
|--------------------------|----------------------------------------------------------|
| `%(LF)`                  | caller's `"File: … LineNo:…"` location (no argument)     |
| `%(j)` / `%(json-indent)`| next argument as **indented** JSON                       |
| `%(J)` / `%(json)`       | next argument as **compact** JSON                        |
| `%(Red)` `%(Green)` …    | ANSI color (see `ColorTab`); unknown names default to red |
| `%(!)` / `%(Reset)`      | reset color to the terminal default                      |

Any color directive is automatically reset at the end of the output. A missing
argument for `%(j)`/`%(J)` or an ordinary verb produces a visible placeholder
instead of panicking.

```go
dbgo.Sprintf("oops %(Red)bad%(!) %(J)", map[string]int{"x": 1})
```

Use `ProcessFormat` directly if you need the expanded format and argument
slice without printing.

## Gated debug output

`DbPrintf` and `DbFprintf` print only when a named debug flag has been enabled
with `SetDbFlag`; test a flag with `IsDbOn`.

```go
dbgo.SetDbFlag(map[string]bool{"db": true, "net": false})
dbgo.DbPrintf("db", "query=%s\n", q) // prints
dbgo.DbPrintf("net", "dial %s\n", a) // silent
```

`DbPf`, `DbFpf` gate on a plain `bool`; `DbPfb` writes to both stdout and
stderr; `DbPfe` gates on an environment variable via `ChkEnv`.

## Environment helpers

`ChkEnv` reports whether an environment variable parses as a true value (see
`ParseBool`) and **caches** the result, so `os.Getenv` is called at most once
per variable name. `SetFlag` overrides the environment programmatically (and is
concurrency-safe).

```go
if dbgo.ChkEnv("DEBUG_NET") {
	// ... verbose networking ...
}
```

`ParseBool` treats these as true (case variants included):
`t`, `yes`, `1`, `true`, `on`. Everything else (including the empty string) is
false.

## Small utilities

| Function                          | Description                                              |
|-----------------------------------|----------------------------------------------------------|
| `InArray(s, arr)` / `InArrayString` | index of `s` in `arr`, or `-1`                         |
| `InArrayInt(n, arr)`              | index of `n` in `arr`, or `-1`                           |
| `StdErrPiped` / `StdOutPiped` / `StdInPiped` | is the stream redirected to a pipe/file?    |
| `BackTickQuote(s)`                | backtick-quote `s` for shell use, escaping embedded backticks |

## Colors

The `Color*` variables hold ANSI escape sequences (red, green, yellow, blue,
black, magenta, cyan, plus `ColorReset` and several on-white/on-blue
combinations). They are always populated; callers that write to a
non-terminal destination and want to suppress escapes can check `StdOutPiped`
first. The `ColorTab` map powers the `%(name)` color directives.

## Status

The `LogData`/`LogToFile` types are **experimental** and currently no-ops (see
the warning in `log.go`).

## License

BSD 3-Clause — see [LICENSE](LICENSE).
