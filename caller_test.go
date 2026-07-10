package dbgo

import (
	"io"
	"os"
	"strings"
	"testing"
)

// captureStdout temporarily redirects os.Stdout to capture output from
// functions (like IAmAt) that print directly to standard output.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}
	os.Stdout = w

	out := make(chan string)
	go func() {
		b, err := io.ReadAll(r)
		if err != nil {
			out <- ""
			return
		}
		out <- string(b)
	}()

	fn()
	w.Close()
	os.Stdout = old
	return <-out
}

func TestCallerDepth(t *testing.T) {
	if callerDepth(nil) != 1 {
		t.Errorf("callerDepth(nil) = %d, want 1", callerDepth(nil))
	}
	if callerDepth([]int{}) != 1 {
		t.Errorf("callerDepth([]) = %d, want 1", callerDepth([]int{}))
	}
	if got := callerDepth([]int{5}); got != 5 {
		t.Errorf("callerDepth([5]) = %d, want 5", got)
	}
}

func TestLINE(t *testing.T) {
	if LINE() == "LineNo:Unk" {
		t.Errorf("LINE() = Unk, expected a number")
	}
	if LINE(0) == "LineNo:Unk" {
		t.Errorf("LINE(0) = Unk, expected a number")
	}
}

func TestFILE(t *testing.T) {
	f := FILE()
	if !strings.HasSuffix(f, "caller_test.go") {
		t.Errorf("FILE() = %q, want path ending in caller_test.go", f)
	}
}

func TestFUNCNAME(t *testing.T) {
	fn := FUNCNAME()
	if !strings.Contains(fn, "TestFUNCNAME") {
		t.Errorf("FUNCNAME() = %q, want it to contain TestFUNCNAME", fn)
	}
}

func TestLINEnf(t *testing.T) {
	line, file := LINEnf()
	if line <= 0 {
		t.Errorf("LINEnf line = %d, want > 0", line)
	}
	if !strings.HasSuffix(file, "caller_test.go") {
		t.Errorf("LINEnf file = %q, want caller_test.go", file)
	}
}

func TestLF(t *testing.T) {
	s := LF()
	if !strings.Contains(s, "caller_test.go") || !strings.Contains(s, "LineNo:") {
		t.Errorf("LF() = %q, want caller_test.go and LineNo:", s)
	}
	// depth 2 points further up the stack but still yields a valid location.
	if s2 := LF(2); !strings.Contains(s2, "LineNo:") {
		t.Errorf("LF(2) = %q, want LineNo:", s2)
	}
}

func TestLFNegativeDepth(t *testing.T) {
	// Negative depth enables the cross-file walk; the result is non-empty and
	// still carries LineNo markers.
	s := LF(-1)
	if s == "" {
		t.Errorf("LF(-1) = empty, want a non-empty walk-back string")
	}
	if !strings.Contains(s, "LineNo:") {
		t.Errorf("LF(-1) = %q, want LineNo: markers", s)
	}
}

func TestLFj(t *testing.T) {
	s := LFj()
	if !strings.Contains(s, `"File":`) || !strings.Contains(s, `"LineNo":`) {
		t.Errorf("LFj() = %q, want a JSON File/LineNo fragment", s)
	}
}

func TestLF2(t *testing.T) {
	line, file := LF2()
	if line <= 0 {
		t.Errorf("LF2 line = %d, want > 0", line)
	}
	if !strings.HasSuffix(file, "caller_test.go") {
		t.Errorf("LF2 file = %q, want caller_test.go", file)
	}
}

func TestIAmAt(t *testing.T) {
	out := captureStdout(t, func() {
		IAmAt("hello", "world")
	})
	if !strings.Contains(out, "Func:") || !strings.Contains(out, "hello world") {
		t.Errorf("IAmAt output = %q, want Func: and the joined message", out)
	}
}

func TestIAmAt2(t *testing.T) {
	out := captureStdout(t, func() {
		IAmAt2("msg")
	})
	// Reports the current frame and a "called..." frame above it.
	if !strings.Contains(out, "called...") {
		t.Errorf("IAmAt2 output = %q, want it to contain called...", out)
	}
}
