package main

import (
	"fmt"

	"git.q8s.co/pschlump/dbgo"
)

func main() {
	x := dbgo.ChkEnv("YepYep")
	fmt.Printf("%v\n", x)
}
