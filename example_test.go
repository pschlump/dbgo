package dbgo_test

import (
	"fmt"

	"github.com/pschlump/dbgo"
)

// The location helpers report where they were called from. Their output
// depends on the source file and line number, so these examples do not assert
// an exact output.

func ExampleLF() {
	fmt.Println("called from:", dbgo.LF())
}

func ExampleLINE() {
	fmt.Println("line:", dbgo.LINE())
}

func ExampleFILE() {
	fmt.Println("file:", dbgo.FILE())
}

func ExampleFUNCNAME() {
	fmt.Println("in:", dbgo.FUNCNAME())
}

func ExampleSVar() {
	// Compact JSON.
	fmt.Println(dbgo.SVar([]int{1, 2, 3}))
	// Output: [1,2,3]
}

func ExampleSVarI() {
	// Tab-indented JSON.
	fmt.Println(dbgo.SVarI([]int{1, 2, 3}))
	// Output:
	// [
	// 	1,
	// 	2,
	// 	3
	// ]
}

func ExampleParseBool() {
	fmt.Println(dbgo.ParseBool("yes"), dbgo.ParseBool("no"))
	// Output: true false
}

func ExampleInArray() {
	fmt.Println(dbgo.InArray("b", []string{"a", "b", "c"}))
	// Output: 1
}

func ExampleBackTickQuote() {
	// Embedded backticks are escaped as \x60.
	fmt.Println(dbgo.BackTickQuote("a`b"))
	// Output: `a\x60b`
}

// ChkEnv consults a cached boolean derived from the environment; SetFlag lets
// examples and tests set a value deterministically.

func ExampleChkEnv() {
	dbgo.SetFlag("DEBUG_VERBOSE_EXAMPLE", "yes")
	if dbgo.ChkEnv("DEBUG_VERBOSE_EXAMPLE") {
		fmt.Println("verbose mode on")
	}
	// Output: verbose mode on
}

// The extended printf understands %(name) directives. The color example below
// is shown without an asserted output because it contains terminal escapes.

func ExampleSprintf_json() {
	// %(J) renders the next argument as compact JSON.
	fmt.Println(dbgo.Sprintf("data=%(J)", []int{1, 2, 3}))
	// Output: data=[1,2,3]
}

func ExampleSprintf_indent() {
	// %(j) renders the next argument as indented JSON.
	fmt.Println(dbgo.Sprintf("data=%(j)", []int{1, 2, 3}))
	// Output:
	// data=[
	// 	1,
	// 	2,
	// 	3
	// ]
}

func ExamplePrintf() {
	// Color the word "bad" red, then reset. %(LF) inserts the caller's
	// location. (Output not asserted: it contains ANSI escapes and a path.)
	dbgo.Printf("at %(LF), status=%(Red)bad%(!)\n")
}

// Gated debug output: DbPrintf only prints when the flag is enabled.

func ExampleSetDbFlag() {
	dbgo.SetDbFlag(map[string]bool{"demo": true})
	dbgo.DbPrintf("demo", "this prints\n")
	dbgo.SetDbFlag(map[string]bool{"demo": false})
	dbgo.DbPrintf("demo", "this is silent\n")
	// Output: this prints
}
