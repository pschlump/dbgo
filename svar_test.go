package dbgo

import (
	"strings"
	"testing"
)

func TestSVar(t *testing.T) {
	cases := []struct {
		name string
		v    interface{}
		want string
	}{
		{"int slice", []int{1, 2, 3}, "[1,2,3]"},
		{"string", "hi", `"hi"`},
		{"int", 42, "42"},
		{"nil", nil, "null"},
		{"bool", true, "true"},
	}
	for _, c := range cases {
		if got := SVar(c.v); got != c.want {
			t.Errorf("SVar(%s) = %q, want %q", c.name, got, c.want)
		}
	}
}

func TestSVarIIndented(t *testing.T) {
	got := SVarI([]int{1, 2, 3})
	if !strings.Contains(got, "\n") || !strings.Contains(got, "\t1") {
		t.Errorf("SVarI = %q, want tab-indented JSON", got)
	}
}

func TestSVarStruct(t *testing.T) {
	type person struct {
		Name string
		Age  int
	}
	got := SVar(person{Name: "Ada", Age: 36})
	if !strings.Contains(got, `"Name":"Ada"`) || !strings.Contains(got, `"Age":36`) {
		t.Errorf("SVar(struct) = %q, want Name and Age fields", got)
	}
}

func TestSVarUnsupportedTypeNoPanic(t *testing.T) {
	// Channels are not JSON-marshalable with encoding/json. The custom json
	// fork may either error or elide the value; either way SVar must not
	// panic.
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("SVar panicked on chan: %v", r)
		}
	}()
	if got := SVar(make(chan int)); got == "" {
		t.Errorf("SVar(chan) = empty, want a non-empty result")
	}
}
