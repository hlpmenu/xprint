package main

import (
	"fmt"
	"log"

	xprint "gopkg.hlmpn.dev/pkg/xprint"
)

// quickTest checks if fmt.Sprintf and xprint.Printf produce identical outputs for simple cases
func Quicktest() {
	log.Println("Running quick validation test...")

	testString := "Hello %s, the number is %d and the float is %.2f"
	name := "Alice"
	number := 42
	floatNum := 3.1415

	fmtOutput := fmt.Sprintf(testString, name, number, floatNum)
	xprintOutput := xprint.Printf(testString, name, number, floatNum)

	// Print results for manual inspection
	log.Printf("fmt.Sprintf output: %s", fmtOutput)
	log.Printf("xprint.Printf output: %s", xprintOutput)

	// If outputs differ, fail early
	if fmtOutput != xprintOutput {
		log.Fatal("ERROR: Quick test failed! Outputs are different.")
	}

	log.Println("Quick test passed ✅")
}
