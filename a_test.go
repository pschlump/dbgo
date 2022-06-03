package dbgo

import (
	"fmt"
	"testing"
)

// -----------------------------------------------------------------------------------------------------------------------------------------------
func Test_InArray(t *testing.T) {
	tests := []struct {
		lookFor        string
		ArrayOfStrings []string
		expected       int
	}{
		{"abc", []string{"def", "abc", "ghi"}, 1},
		{"a1c", []string{"def", "abc", "ghi"}, -1},
		{"abc", []string{}, -1},
		{"abc", []string{"abc", "abc", "ghi"}, 0},
		{"abc", []string{"def", "aXc", "abc"}, 2},
	}

	for ii, test := range tests {

		got := InArray(test.lookFor, test.ArrayOfStrings)
		if got != test.expected {
			t.Errorf("Error %2d, got: %d, expected %d\n", ii, got, test.expected)
		}

	}
}

func Test_IsTerminal(t *testing.T) {
	a := StdErrPiped()
	b := StdInPiped()
	c := StdInPiped()
	if a != b || b != c {
		t.Errorf("Error, test of Std...Piped - unlikey result for running tests\n")
	}

	if a {
		if ColorRed != "" {
			t.Errorf("Error, test of ColorRed - should be empty string, got %x\n", ColorRed)
		}
	} else {
		if ColorRed == "" {
			t.Errorf("Error, test of ColorRed - should be non empty string, got empty string\n")
		}
	}
	fmt.Printf("->%s<- ->%s<- ->%s<-\n", ColorRed, ColorGreen, ColorReset)
}

/* vim: set noai ts=4 sw=4: */
