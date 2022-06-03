package dbgo

import (
	"regexp"
	"testing"
)

//
// Leave Blank Comments Alone - to put 1st test on line 19
//
//
//
//
//
//

func TestFastChop(t *testing.T) {

	s := LF() // Must stay on line 19 - to match regexp on next line.
	match, err := regexp.MatchString(".*debug_test.go.*LineNo:19", s)
	if err != nil {
		t.Errorf("Error 00, Used invalid regular expression in test code, %s\n", err)
	}
	if !match {
		t.Errorf("Error 01, Expected match to 'File: /Users/corwin/lib/go-lib/debug/debug_test.go LineNo:9', got %s\n", s)
	}

	s = SVar([]int{1, 2, 3})
	if s != "[1,2,3]" {
		t.Errorf("Error 02, Expected [1,3,2], got %s\n", s)
	}
}

func TestParseBool(t *testing.T) {
	tf := ParseBool("Yes")
	if !tf {
		t.Errorf("Error ParseBool failed\n")
	}
	tf = ParseBool("")
	if tf {
		t.Errorf("Error ParseBool failed\n")
	}
}

func TestEnvCheck(t *testing.T) {
	tf := ChkEnv("some-undefined-env-var--22323232323232323232323232")
	if tf {
		t.Errorf("Error ChkEnv 1 failed\n")
	}
	for i := 0; i < 50000; i++ {
		tf = ChkEnv("some-undefined-env-var--22323232323232323232323232")
		if tf {
			t.Errorf("Error ChkEnv 1 failed\n")
		}
	}
}
