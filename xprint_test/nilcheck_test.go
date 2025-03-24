package xprint_test

import (
	"fmt"
	"testing"

	xprint "gopkg.hlmpn.dev/pkg/xprint"
)

func TestNilChecks2(t *testing.T) {
	t.Log("Running nil check tests...")

	// Test nil with %v format
	t.Run("NilWithVerb", func(t *testing.T) {
		xprintResult := xprint.Printf("%v", nil)
		fmtResult := fmt.Sprintf("%v", nil)

		if xprintResult != fmtResult {
			t.Errorf("Mismatch between fmt.Sprintf and xprint.Printf with nil\nxprint: %q\nfmt: %q",
				xprintResult, fmtResult)
		} else {
			t.Logf("Success: Test passed for nil with %%v format")
		}
	})

	// Test string with no format verbs
	t.Run("PlainString", func(t *testing.T) {
		xprintResult := xprint.Printf("Hello world")
		fmtResult := fmt.Sprintf("Hello world")

		if xprintResult != fmtResult {
			t.Errorf("Mismatch between fmt.Sprintf and xprint.Printf with plain string\nxprint: %q\nfmt: %q",
				xprintResult, fmtResult)
		} else {
			t.Logf("Success: Test passed for plain string")
		}
	})
}
