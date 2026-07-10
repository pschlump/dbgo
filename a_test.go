package dbgo

import "testing"

// Test_IsTerminal verifies the Std*Piped helpers run without panicking and that
// the color constants are populated. The previous version of this test asserted
// that all three streams reported the same piped state and that ColorRed was
// empty when piped; neither holds in general (streams are independent, and
// colors are always initialized).
func Test_IsTerminal(t *testing.T) {
	// Each helper must return a bool without panicking. The three streams are
	// independent, so their answers need not agree.
	_ = StdErrPiped()
	_ = StdOutPiped()
	_ = StdInPiped()

	// Colors are always initialized regardless of TTY status.
	for name, c := range map[string]string{
		"ColorRed":    ColorRed,
		"ColorGreen":  ColorGreen,
		"ColorReset":  ColorReset,
		"ColorYellow": ColorYellow,
	} {
		if c == "" {
			t.Errorf("%s is empty; expected an ANSI escape sequence", name)
		}
	}
}
