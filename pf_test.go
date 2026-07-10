package dbgo

import (
	"bytes"
	"strings"
	"testing"
)

// TestSprintfColorAndJSON verifies color directives, the auto-reset behavior,
// and the %(j)/%(J) JSON directives. Expected values are built from the
// package's own color constants so the assertions do not depend on hand-typed
// escape sequences.
func TestSprintfColorAndJSON(t *testing.T) {
	// A color directive turns colorFound on, so a Reset is appended at the end.
	got := Sprintf("abc%(Red)def\n")
	want := "abc" + ColorRed + "def\n" + ColorReset
	if got != want {
		t.Errorf("color directive:\n got = %q\nwant = %q", got, want)
	}

	// Green followed by indented JSON via %(j), then the trailing auto-reset.
	got = Sprintf("abc%(Green)%(j)def\n", []string{"xxx", "YYY"})
	want = "abc" + ColorGreen + SVarI([]string{"xxx", "YYY"}) + "def\n" + ColorReset
	if got != want {
		t.Errorf("green + json:\n got = %q\nwant = %q", got, want)
	}

	// %(!) / %(Reset) reset immediately, so no second reset is appended at end.
	got = Sprintf("x%(Red)y%(!)z")
	want = "x" + ColorRed + "y" + ColorReset + "z"
	if got != want {
		t.Errorf("reset directive:\n got = %q\nwant = %q", got, want)
	}
}

func TestProcessFormatLF(t *testing.T) {
	got, params := ProcessFormat("at %(LF) end", nil)
	if len(params) != 0 {
		t.Errorf("LF should consume no params, got %v", params)
	}
	if !strings.Contains(got, "LineNo:") {
		t.Errorf("LF output = %q, want it to contain LineNo:", got)
	}
}

func TestProcessFormatJSONAliases(t *testing.T) {
	// %(J) and %(json) render compact JSON.
	for _, dir := range []string{"%(J)", "%(json)"} {
		got, params := ProcessFormat("x"+dir, []interface{}{[]int{1, 2, 3}})
		if !strings.HasSuffix(got, "[1,2,3]") {
			t.Errorf("%s: got %q, want suffix [1,2,3]", dir, got)
		}
		if len(params) != 0 {
			t.Errorf("%s consumed a params slot: %v", dir, params)
		}
	}
	// %(j) and %(json-indent) render indented JSON.
	for _, dir := range []string{"%(j)", "%(json-indent)"} {
		got, _ := ProcessFormat("x"+dir, []interface{}{[]int{1, 2, 3}})
		if !strings.Contains(got, "\n\t1,") {
			t.Errorf("%s: got %q, want indented JSON", dir, got)
		}
	}
}

func TestProcessFormatMissingJSONArg(t *testing.T) {
	// Missing arguments for %(j)/%(J) must not panic.
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("ProcessFormat panicked on missing arg: %v", r)
		}
	}()
	got, _ := ProcessFormat("v=%(J)", nil)
	if !strings.Contains(got, "Missing Value") {
		t.Errorf("missing compact-JSON arg: got %q, want a Missing Value placeholder", got)
	}
	got, _ = ProcessFormat("v=%(j)", nil)
	if !strings.Contains(got, "Missing Value") {
		t.Errorf("missing indented-JSON arg: got %q, want a Missing Value placeholder", got)
	}
}

func TestProcessFormatPercentLiteral(t *testing.T) {
	got, params := ProcessFormat("100%% done", nil)
	if want := "100% done"; got != want {
		t.Errorf("%% handling: got %q, want %q", got, want)
	}
	if len(params) != 0 {
		t.Errorf("%% should consume no params, got %v", params)
	}
}

func TestProcessFormatMissingPlainArg(t *testing.T) {
	got, params := ProcessFormat("hi %s %d", []interface{}{"world"})
	if !strings.Contains(got, "%s") || !strings.Contains(got, "%d") {
		t.Errorf("plain verbs should remain in format: got %q", got)
	}
	if len(params) != 2 {
		t.Fatalf("params = %v (len %d), want len 2", params, len(params))
	}
	if params[0] != "world" {
		t.Errorf("params[0] = %v, want world", params[0])
	}
	placeholder, ok := params[1].(string)
	if !ok || !strings.Contains(placeholder, "Missing Value") {
		t.Errorf("params[1] = %v, want a Missing Value placeholder", params[1])
	}
}

func TestFprintf(t *testing.T) {
	var buf bytes.Buffer
	if _, err := Fprintf(&buf, "v=%(J)", []int{1}); err != nil {
		t.Fatalf("Fprintf error: %v", err)
	}
	if !strings.HasSuffix(buf.String(), "[1]") {
		t.Errorf("Fprintf output = %q, want suffix [1]", buf.String())
	}
}

func TestDbFlags(t *testing.T) {
	SetDbFlag(map[string]bool{"flag-a": true})
	if !IsDbOn("flag-a") {
		t.Errorf("IsDbOn(flag-a) = false, want true")
	}
	if IsDbOn("flag-b") {
		t.Errorf("IsDbOn(flag-b) = true, want false")
	}
	SetDbFlag(map[string]bool{"flag-a": false})
	if IsDbOn("flag-a") {
		t.Errorf("IsDbOn(flag-a) = true after disable, want false")
	}
}

func TestDbFprintf(t *testing.T) {
	var buf bytes.Buffer

	// Flag off: nothing written.
	DbFprintf("pf-off", &buf, "hello %(LF)\n")
	if buf.Len() != 0 {
		t.Errorf("expected no output when flag off, got %q", buf.String())
	}

	// Flag on: output written and byte count reported.
	SetDbFlag(map[string]bool{"pf-on": true})
	n, err := DbFprintf("pf-on", &buf, "num=%d\n", 7)
	if err != nil {
		t.Fatalf("DbFprintf error: %v", err)
	}
	if n != buf.Len() {
		t.Errorf("reported n=%d != buf.Len()=%d", n, buf.Len())
	}
	if !strings.Contains(buf.String(), "num=7\n") {
		t.Errorf("output = %q, want it to contain num=7\\n", buf.String())
	}
}

func TestDbPfOff(t *testing.T) {
	// Gated printers must be silent and error-free when their flag is false.
	if n, err := DbPf(false, "x"); n != 0 || err != nil {
		t.Errorf("DbPf(false) = (%d,%v), want (0,nil)", n, err)
	}
	if n, err := DbPfb(false, "x"); n != 0 || err != nil {
		t.Errorf("DbPfb(false) = (%d,%v), want (0,nil)", n, err)
	}
}
