package dbgo

import (
	"strings"
	"testing"
)

func TestStdPipedHelpers(t *testing.T) {
	// Each must return a bool without panicking; the actual value depends on
	// how the test runner wires the standard streams.
	_ = StdErrPiped()
	_ = StdOutPiped()
	_ = StdInPiped()
}

func TestColorTab(t *testing.T) {
	for _, name := range []string{"red", "Red", "green", "Green", "yellow", "reset", "Reset", "!", "magenta", "cyan"} {
		v, ok := ColorTab[name]
		if !ok {
			t.Errorf("ColorTab[%q] missing", name)
			continue
		}
		if v == "" {
			t.Errorf("ColorTab[%q] is empty", name)
		}
	}
	if ColorTab["!"] != ColorReset {
		t.Errorf("ColorTab[!] = %q, want ColorReset (%q)", ColorTab["!"], ColorReset)
	}
	if ColorTab["Red"] != ColorRed {
		t.Errorf("ColorTab[Red] != ColorRed")
	}
	if ColorTab["Green"] != ColorGreen {
		t.Errorf("ColorTab[Green] != ColorGreen")
	}
}

func TestBackTickQuote(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"plain", "`plain`"},
		{"", "``"},
		{"a`b", "`a\\x60b`"},
		{"x``y", "`x\\x60\\x60y`"},
	}
	for _, c := range cases {
		if got := BackTickQuote(c.in); got != c.want {
			t.Errorf("BackTickQuote(%q) = %q, want %q", c.in, got, c.want)
		}
	}
	// No raw backtick should survive in the interior of the result.
	trimmed := strings.Trim(BackTickQuote("a`b"), "`")
	if strings.Contains(trimmed, "`") {
		t.Errorf("BackTickQuote leaked an interior backtick: %q", trimmed)
	}
}
