package dbgo

import (
	"fmt"
	"testing"
)

func TestPrintf(t *testing.T) {
	e := `abc[0;31;40mdef
[0m`
	a := Sprintf("abc%(Red)def\n")
	fmt.Printf("->%s<-\n", a)
	if a != e {
		t.Errorf("Failed to produce red")
	}
	a = Sprintf("abc%(Green)%(j)def\n", []string{"xxx", "YYY"})
	fmt.Printf("->%s<-\n", a)
	e = `abc[0;32;40m[
	"xxx",
	"YYY"
]def
[0m`
	if a != e {
		t.Errorf("Failed to produce json")
	}
}
