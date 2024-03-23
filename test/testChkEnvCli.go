package main

import (
	"fmt"

	"github.com/pschlump/dbgo"
)

func main() {
	x := dbgo.ChkEnv("YepYep")
	fmt.Printf("%v\n", x)
}
