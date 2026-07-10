package dbgo

import (
	"regexp"
	"testing"
)

// TestFastChop checks LF() reports this file and a line number, and that SVar
// produces compact JSON. The previous version depended on LF() sitting on a
// specific source line and embedded a hardcoded path in the error message.
func TestFastChop(t *testing.T) {
	s := LF()
	// LF() returns "File: <path> LineNo:<n>"; from here the path is this test
	// file and n is a positive number.
	match, err := regexp.MatchString(`debug_test\.go LineNo:[0-9]+`, s)
	if err != nil {
		t.Fatalf("invalid regexp in test: %v", err)
	}
	if !match {
		t.Errorf("LF() = %q, want a string containing this file name and a line number", s)
	}

	if got := SVar([]int{1, 2, 3}); got != "[1,2,3]" {
		t.Errorf("SVar([1,2,3]) = %q, want %q", got, "[1,2,3]")
	}
}
