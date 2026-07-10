// dbgo_example is a runnable tour of the github.com/pschlump/dbgo helpers.
//
// Run it with:
//
//	go run ./example
package main

import (
	"fmt"

	"github.com/pschlump/dbgo"
)

func main() {
	// --- Location helpers -------------------------------------------------
	//
	// LF/LINE/FILE/FUNCNAME report where they were called from via the
	// runtime call stack. Each takes an optional depth (default 1 = caller).
	fmt.Println("--- location helpers ---")
	fmt.Printf("LF():        %s\n", dbgo.LF())
	fmt.Printf("LINE():      %s\n", dbgo.LINE())
	fmt.Printf("FILE():      %s\n", dbgo.FILE())
	fmt.Printf("FUNCNAME():  %s\n", dbgo.FUNCNAME())

	// --- JSON helpers -----------------------------------------------------
	fmt.Println("\n--- json helpers ---")
	data := map[string]int{"a": 1, "b": 2}
	fmt.Printf("SVar (compact):   %s\n", dbgo.SVar(data))
	fmt.Printf("SVarI (indented):\n%s\n", dbgo.SVarI(data))

	// --- Extended printf -------------------------------------------------
	//
	// %(LF) location, %(J) compact JSON, %(j) indented JSON, %(Red) color
	// with %(!) to reset.
	fmt.Println("\n--- extended printf ---")
	dbgo.Printf("at %(LF), compact=%(J)\n", []string{"x", "y"})
	dbgo.Printf("%(Red)red text%(!) and %(Green)green text%(!)\n")

	// --- Gated debug output ----------------------------------------------
	//
	// DbPrintf prints only when the named flag is enabled with SetDbFlag.
	fmt.Println("\n--- gated debug output ---")
	dbgo.SetDbFlag(map[string]bool{"network": true, "cache": false})
	dbgo.DbPrintf("network", "sending request, payload=%(J)\n", data)
	dbgo.DbPrintf("cache", "this line is suppressed\n")

	// --- Environment-gated output ----------------------------------------
	//
	// ChkEnv reads (and caches) a boolean from the environment. SetFlag
	// overrides it programmatically.
	fmt.Println("\n--- environment helpers ---")
	dbgo.SetFlag("DEBUG_VERBOSE", "yes")
	if dbgo.ChkEnv("DEBUG_VERBOSE") {
		fmt.Println("DEBUG_VERBOSE is on")
	}
	fmt.Printf("ParseBool(\"on\")  = %v\n", dbgo.ParseBool("on"))
	fmt.Printf("ParseBool(\"no\")  = %v\n", dbgo.ParseBool("no"))

	// --- Small utilities -------------------------------------------------
	fmt.Println("\n--- utilities ---")
	fmt.Printf("InArray(\"b\", [a b c]) = %d\n", dbgo.InArray("b", []string{"a", "b", "c"}))
	fmt.Printf("BackTickQuote(\"a`b\")  = %s\n", dbgo.BackTickQuote("a`b"))
	fmt.Printf("StdOutPiped()          = %v\n", dbgo.StdOutPiped())
}
