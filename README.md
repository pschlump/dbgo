# godebug

 [![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/pschlump/Go-FTL/master/LICENSE)

Small library of tools for debugging Go (golang) programs.

The most useful of these is LF().  It returns a line number
and file name.  The parameter is an optional number.
1 indicates that you want it for the current call. 2 would
be the parent of the current call.

Example:

```golang

	package main

	import (
		"fmt"

		"git.q8s.co/pschlump/dbgo"
	)

	func main() {
		dbgo.Printf("I am at: %(LF), data %(J)\n", []string{1,2,3})
	}

```

`LF()` takes an optional parameter, 1 indicates the current
function, 2 is the parent, 3 the grandparent.

I will add complete documentation tomorrow.

Formats:

| Format     | Description                                                |
|------------|------------------------------------------------------------|
| `%(LF)`    | print out line number                                      |
| `%(j)`     | convert data using json.Encode and output with indentation |
| `%(J)`     | convert data using json.Encode and output                  |


## ChkEnv

Return true if the passed environment variable can be parsed to be a `true` value.
Values are cached so that the os.Getenv is only called once for each variable.

There is a test in ./test/test.sh to test this (or use the Makefile).

