package xprint_test

import (
	"errors"
	"fmt"
	"math/rand"
	"testing"

	xprint "gopkg.hlmpn.dev/pkg/xprint"
)

// BenchmarkItem defines a single benchmark test case
type BenchmarkItem struct {
	Name       string
	Format     string
	Args       []any
	Iterations int
}

// Runs benchmarks for various data types
func BenchmarkTypes(b *testing.B) {
	b.Log("Running benchmark for various data types...")

	// Create benchmark cases for different types
	benchmarks := []BenchmarkItem{
		{
			Name:       "String",
			Format:     "String: %s",
			Args:       []any{"Hello, world!"},
			Iterations: 1000,
		},
		{
			Name:       "Int",
			Format:     "Int: %d",
			Args:       []any{42},
			Iterations: 1000,
		},
		{
			Name:       "Int64",
			Format:     "Int64: %d",
			Args:       []any{int64(9223372036854775807)},
			Iterations: 1000,
		},
		{
			Name:       "Float",
			Format:     "Float: %.2f",
			Args:       []any{3.14159},
			Iterations: 1000,
		},
		{
			Name:       "Bool",
			Format:     "Bool: %t",
			Args:       []any{true},
			Iterations: 1000,
		},
		{
			Name:       "Error",
			Format:     "Error: %v",
			Args:       []any{errors.New("test error")},
			Iterations: 1000,
		},
		{
			Name:       "Bytes",
			Format:     "Bytes: %s",
			Args:       []any{[]byte("byte slice")},
			Iterations: 1000,
		},
		{
			Name:       "Complex",
			Format:     "Complex format with multiple args: %s - %d - %t - %f",
			Args:       []any{"text", 42, true, 3.14},
			Iterations: 500,
		},
		{
			Name:       "Map",
			Format:     "Map: %v",
			Args:       []any{map[string]int{"one": 1, "two": 2}},
			Iterations: 500,
		},
		{
			Name:       "Slice",
			Format:     "Slice: %v",
			Args:       []any{[]int{1, 2, 3, 4, 5}},
			Iterations: 500,
		},
		{
			Name:       "Pointer",
			Format:     "Pointer: %p",
			Args:       []any{&struct{}{}},
			Iterations: 500,
		},
	}

	// Shuffle the benchmarks for randomized order
	rand.Shuffle(len(benchmarks), func(i, j int) {
		benchmarks[i], benchmarks[j] = benchmarks[j], benchmarks[i]
	})

	// Run each benchmark
	for _, bench := range benchmarks {
		b.Run(bench.Name, func(b *testing.B) {
			b.Run("xprint", func(b *testing.B) {
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					_ = xprint.Printf(bench.Format, bench.Args...)
				}
			})

			b.Run("fmt", func(b *testing.B) {
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					_ = fmt.Sprintf(bench.Format, bench.Args...)
				}
			})

			// Run one more time to validate results
			xprintResult := xprint.Printf(bench.Format, bench.Args...)
			fmtResult := fmt.Sprintf(bench.Format, bench.Args...)

			// For maps, we don't compare exact strings since order may differ
			skipExactComparison := bench.Name == "Map"

			if xprintResult != fmtResult && !skipExactComparison {
				b.Errorf("[%s] MISMATCH: Results don't match!", bench.Name)
				b.Logf("xprint output: %s", xprintResult)
				b.Logf("fmt output: %s", fmtResult)
			}
		})
	}
}

func TestFloatFormatting(t *testing.T) {
	t.Log("Testing float formatting issues...")

	// Test with different precision formats
	t.Run("Different precision formats", func(t *testing.T) {
		var f64 float64 = 3.14159

		precisions := []string{
			"%.0f", "%.1f", "%.2f", "%.3f", "%.6f", "%f", "%g", "%e",
		}

		for _, prec := range precisions {
			t.Logf("Testing precision format: %s", prec)

			// Use separate format strings for the description and the actual format
			fmt1 := fmt.Sprintf(prec, f64)
			xprint1 := xprint.Printf(prec, f64)

			t.Logf("fmt output: '%s'", fmt1)
			t.Logf("xprint output: '%s'", xprint1)

			if fmt1 != xprint1 {
				t.Errorf("MISMATCH for format %s", prec)
			}
		}
	})

	// Test raw float constants (interface{} type inference)
	t.Run("Raw float constants", func(t *testing.T) {
		// First, test fmt self-consistency
		fmt1 := fmt.Sprintf("Float: %f", 3.14159)
		fmt2 := fmt.Sprintf("Float: %f", 3.14159)

		if fmt1 != fmt2 {
			t.Errorf("fmt inconsistent with itself!")
		} else {
			t.Logf("fmt is consistent with itself")
		}

		// Then, test xprint self-consistency
		xprint1 := xprint.Printf("Float: %f", 3.14159)
		xprint2 := xprint.Printf("Float: %f", 3.14159)

		if xprint1 != xprint2 {
			t.Errorf("xprint inconsistent with itself!")
		} else {
			t.Logf("xprint is consistent with itself")
		}

		// Compare fmt vs xprint
		t.Logf("fmt output: '%s'", fmt1)
		t.Logf("xprint output: '%s'", xprint1)

		if fmt1 != xprint1 {
			t.Errorf("MISMATCH: fmt vs xprint for raw float")
		}
	})

	// Test with explicit float32
	t.Run("Explicit float32 variable", func(t *testing.T) {
		var f32 float32 = 3.14159

		fmt1 := fmt.Sprintf("Float32: %f", f32)
		fmt2 := fmt.Sprintf("Float32: %f", f32)

		if fmt1 != fmt2 {
			t.Errorf("fmt inconsistent with itself for float32!")
		}

		xprint1 := xprint.Printf("Float32: %f", f32)
		xprint2 := xprint.Printf("Float32: %f", f32)

		if xprint1 != xprint2 {
			t.Errorf("xprint inconsistent with itself for float32!")
		}

		t.Logf("fmt float32: '%s'", fmt1)
		t.Logf("xprint float32: '%s'", xprint1)

		if fmt1 != xprint1 {
			t.Errorf("MISMATCH: fmt vs xprint for float32")
		}
	})
}
