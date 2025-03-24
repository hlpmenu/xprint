package xprint_test

import (
	"fmt"
	"testing"

	xprint "gopkg.hlmpn.dev/pkg/xprint"
)

// BenchmarkEmptyArgs tests the performance of Printf with no arguments
func BenchmarkEmptyArgs(b *testing.B) {
	b.Run("NoFormat", func(b *testing.B) {
		const testStr = "This is a test string with no formatting"

		// Benchmark xprint.Printf with no args
		b.Run("xprint", func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = xprint.Printf(testStr)
			}
		})

		// Benchmark fmt.Sprintf with no args
		b.Run("fmt", func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = fmt.Sprintf(testStr)
			}
		})
	})

	b.Run("WithPercentSign", func(b *testing.B) {
		// With a % sign that would normally need to be processed
		const testStrWithPercent = "This has a % sign but no args"

		// Benchmark xprint.Printf with no args but with % sign
		b.Run("xprint", func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = xprint.Printf(testStrWithPercent)
			}
		})

		// Benchmark fmt.Sprintf with no args but with % sign
		b.Run("fmt", func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = fmt.Sprintf(testStrWithPercent)
			}
		})
	})
}

// TestEmptyArgs verifies that xprint.Printf and fmt.Sprintf produce identical output
func TestEmptyArgs(t *testing.T) {
	t.Run("NoFormat", func(t *testing.T) {
		const testStr = "This is a test string with no formatting"

		xprintResult := xprint.Printf(testStr)
		fmtResult := fmt.Sprintf(testStr)

		if xprintResult != fmtResult {
			t.Errorf("Results don't match for string with no formatting\nxprint: %s\nfmt: %s",
				xprintResult, fmtResult)
		} else {
			t.Logf("Results match for string with no formatting")
		}
	})

	t.Run("WithPercentSign", func(t *testing.T) {
		const testStrWithPercent = "This has a % sign but no args"

		xprintResult := xprint.Printf(testStrWithPercent)
		fmtResult := fmt.Sprintf(testStrWithPercent)

		if xprintResult != fmtResult {
			t.Errorf("Results don't match for string with %% sign\nxprint: %s\nfmt: %s",
				xprintResult, fmtResult)
		} else {
			t.Logf("Results match for string with %% sign")
		}
	})
}
