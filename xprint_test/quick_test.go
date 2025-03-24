package xprint_test

import (
	"fmt"
	"testing"

	xprint "gopkg.hlmpn.dev/pkg/xprint"
)

// TestQuick checks if fmt.Sprintf and xprint.Printf produce identical outputs for simple cases
func TestQuick(t *testing.T) {
	t.Log("Running quick validation test...")

	testString := "Hello %s, the number is %d and the float is %.2f"
	name := "Alice"
	number := 42
	floatNum := 3.1415

	fmtOutput := fmt.Sprintf(testString, name, number, floatNum)
	xprintOutput := xprint.Printf(testString, name, number, floatNum)

	// Print results for inspection
	t.Logf("fmt.Sprintf output: %s", fmtOutput)
	t.Logf("xprint.Printf output: %s", xprintOutput)

	// If outputs differ, fail the test
	if fmtOutput != xprintOutput {
		t.Fatalf("ERROR: Quick test failed! Outputs are different.")
	}

	t.Log("Quick test passed ✅")
}
