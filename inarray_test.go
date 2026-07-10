package dbgo

import "testing"

func TestInArray(t *testing.T) {
	cases := []struct {
		lookFor string
		arr     []string
		want    int
	}{
		{"abc", []string{"def", "abc", "ghi"}, 1},
		{"a1c", []string{"def", "abc", "ghi"}, -1},
		{"abc", []string{}, -1},
		{"abc", []string{"abc", "abc", "ghi"}, 0}, // first match wins
		{"abc", []string{"def", "aXc", "abc"}, 2},
		{"last", []string{"a", "b", "last"}, 2},
	}
	for i, c := range cases {
		if got := InArray(c.lookFor, c.arr); got != c.want {
			t.Errorf("%d: InArray(%q, %v) = %d, want %d", i, c.lookFor, c.arr, got, c.want)
		}
	}
}

func TestInArrayStringAndInt(t *testing.T) {
	if got := InArrayString("b", []string{"a", "b", "c"}); got != 1 {
		t.Errorf("InArrayString = %d, want 1", got)
	}
	if got := InArrayString("z", []string{"a", "b", "c"}); got != -1 {
		t.Errorf("InArrayString(missing) = %d, want -1", got)
	}
	if got := InArrayInt(5, []int{1, 3, 5, 7}); got != 2 {
		t.Errorf("InArrayInt = %d, want 2", got)
	}
	if got := InArrayInt(99, []int{1, 3, 5}); got != -1 {
		t.Errorf("InArrayInt(missing) = %d, want -1", got)
	}
	if got := InArrayInt(1, []int{}); got != -1 {
		t.Errorf("InArrayInt(empty) = %d, want -1", got)
	}
}

func TestIndexOf(t *testing.T) {
	if got := indexOf("x", []string{"a", "x"}); got != 1 {
		t.Errorf("indexOf(string) = %d, want 1", got)
	}
	if got := indexOf(7, []int{7}); got != 0 {
		t.Errorf("indexOf(int) = %d, want 0", got)
	}
	if got := indexOf("z", []string{}); got != -1 {
		t.Errorf("indexOf(empty) = %d, want -1", got)
	}
}
