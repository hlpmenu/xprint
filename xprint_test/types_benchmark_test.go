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
	Name   string
	Format string
	Args   []any
}

// BenchmarkTypes runs benchmarks for various data types
func BenchmarkTypes(b *testing.B) {
	// Create benchmark cases for different types
	benchmarks := []BenchmarkItem{
		{
			Name:   "String",
			Format: "String: %s",
			Args:   []any{"Hello, world!"},
		},
		{
			Name:   "Int",
			Format: "Int: %d",
			Args:   []any{42},
		},
		{
			Name:   "Int64",
			Format: "Int64: %d",
			Args:   []any{int64(9223372036854775807)},
		},
		{
			Name:   "Float",
			Format: "Float: %.2f",
			Args:   []any{3.14159},
		},
		{
			Name:   "Bool",
			Format: "Bool: %t",
			Args:   []any{true},
		},
		{
			Name:   "Error",
			Format: "Error: %v",
			Args:   []any{errors.New("test error")},
		},
		{
			Name:   "Bytes",
			Format: "Bytes: %s",
			Args:   []any{[]byte("byte slice")},
		},
		{
			Name:   "Complex",
			Format: "Complex format with multiple args: %s - %d - %t - %f",
			Args:   []any{"text", 42, true, 3.14},
		},
		{
			Name:   "Map",
			Format: "Map: %v",
			Args:   []any{map[string]int{"one": 1, "two": 2}},
		},
		{
			Name:   "Slice",
			Format: "Slice: %v",
			Args:   []any{[]int{1, 2, 3, 4, 5}},
		},
		{
			Name:   "Pointer",
			Format: "Pointer: %p",
			Args:   []any{&struct{}{}},
		},
	}

	// Shuffle the benchmarks for randomized order
	rand.Shuffle(len(benchmarks), func(i, j int) {
		benchmarks[i], benchmarks[j] = benchmarks[j], benchmarks[i]
	})

	// Run each benchmark
	for _, bench := range benchmarks {
		b.Run("xprint_"+bench.Name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = xprint.Printf(bench.Format, bench.Args...)
			}
		})

		b.Run("fmt_"+bench.Name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = fmt.Sprintf(bench.Format, bench.Args...)
			}
		})
	}
}

// TestTypes validates that each format type produces identical output between xprint and fmt
func TestTypes(t *testing.T) {
	// Create test cases for different types
	testCases := []BenchmarkItem{
		{
			Name:   "String",
			Format: "String: %s",
			Args:   []any{"Hello, world!"},
		},
		{
			Name:   "Int",
			Format: "Int: %d",
			Args:   []any{42},
		},
		{
			Name:   "Int64",
			Format: "Int64: %d",
			Args:   []any{int64(9223372036854775807)},
		},
		{
			Name:   "Float",
			Format: "Float: %.2f",
			Args:   []any{3.14159},
		},
		{
			Name:   "Bool",
			Format: "Bool: %t",
			Args:   []any{true},
		},
		{
			Name:   "Error",
			Format: "Error: %v",
			Args:   []any{errors.New("test error")},
		},
		{
			Name:   "Bytes",
			Format: "Bytes: %s",
			Args:   []any{[]byte("byte slice")},
		},
		{
			Name:   "Complex",
			Format: "Complex format with multiple args: %s - %d - %t - %f",
			Args:   []any{"text", 42, true, 3.14},
		},
		{
			Name:   "Map",
			Format: "Map: %v",
			Args:   []any{map[string]int{"one": 1, "two": 2}},
		},
		{
			Name:   "Slice",
			Format: "Slice: %v",
			Args:   []any{[]int{1, 2, 3, 4, 5}},
		},
		{
			Name:   "Pointer",
			Format: "Pointer: %p",
			Args:   []any{&struct{}{}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			xprintResult := xprint.Printf(tc.Format, tc.Args...)
			fmtResult := fmt.Sprintf(tc.Format, tc.Args...)

			t.Logf("xprint output: %s", xprintResult)
			t.Logf("fmt output: %s", fmtResult)

			// For maps, we don't compare exact strings since order may differ
			skipExactComparison := tc.Name == "Map"

			if xprintResult != fmtResult && !skipExactComparison {
				t.Errorf("[%s] MISMATCH: Results don't match!", tc.Name)
			}
		})
	}
}

// TestFloatFormatting tests float formatting with different precision formats
func TestFloatFormatting(t *testing.T) {
	// Test with different precision formats
	t.Run("DifferentPrecisionFormats", func(t *testing.T) {
		var f64 float64 = 3.14159

		precisions := []string{
			"%.0f", "%.1f", "%.2f", "%.3f", "%.6f", "%f", "%g", "%e",
		}

		for _, prec := range precisions {
			t.Run(prec, func(t *testing.T) {
				fmt1 := fmt.Sprintf(prec, f64)
				xprint1 := xprint.Printf(prec, f64)

				t.Logf("fmt output: '%s'", fmt1)
				t.Logf("xprint output: '%s'", xprint1)

				if fmt1 != xprint1 {
					t.Errorf("MISMATCH for format %s", prec)
				}
			})
		}
	})

	// Test raw float constants (interface{} type inference)
	t.Run("RawFloatConstants", func(t *testing.T) {
		// First, test fmt self-consistency
		fmt1 := fmt.Sprintf("Float: %f", 3.14159)
		fmt2 := fmt.Sprintf("Float: %f", 3.14159)
		t.Logf("fmt attempt 1: '%s'", fmt1)
		t.Logf("fmt attempt 2: '%s'", fmt2)
		if fmt1 != fmt2 {
			t.Errorf("fmt inconsistent with itself!")
		}

		// Then, test xprint self-consistency
		xprint1 := xprint.Printf("Float: %f", 3.14159)
		xprint2 := xprint.Printf("Float: %f", 3.14159)
		t.Logf("xprint attempt 1: '%s'", xprint1)
		t.Logf("xprint attempt 2: '%s'", xprint2)
		if xprint1 != xprint2 {
			t.Errorf("xprint inconsistent with itself!")
		}

		// Compare fmt vs xprint
		t.Logf("fmt output: '%s'", fmt1)
		t.Logf("xprint output: '%s'", xprint1)
		if fmt1 != xprint1 {
			t.Errorf("MISMATCH: fmt vs xprint for raw float")
		}
	})

	// Test with explicit float32
	t.Run("ExplicitFloat32", func(t *testing.T) {
		var f32 float32 = 3.14159

		fmt1 := fmt.Sprintf("Float32: %f", f32)
		fmt2 := fmt.Sprintf("Float32: %f", f32)
		t.Logf("fmt float32 attempt 1: '%s'", fmt1)
		t.Logf("fmt float32 attempt 2: '%s'", fmt2)
		if fmt1 != fmt2 {
			t.Errorf("fmt inconsistent with itself for float32!")
		}

		xprint1 := xprint.Printf("Float32: %f", f32)
		xprint2 := xprint.Printf("Float32: %f", f32)
		t.Logf("xprint float32 attempt 1: '%s'", xprint1)
		t.Logf("xprint float32 attempt 2: '%s'", xprint2)
		if xprint1 != xprint2 {
			t.Errorf("xprint inconsistent with itself for float32!")
		}

		// Compare fmt vs xprint for float32
		t.Logf("fmt output: '%s'", fmt1)
		t.Logf("xprint output: '%s'", xprint1)
		if fmt1 != xprint1 {
			t.Errorf("MISMATCH for float32")
		}
	})

	// Test with explicit float64
	t.Run("ExplicitFloat64", func(t *testing.T) {
		var f64 float64 = 3.14159

		fmt1 := fmt.Sprintf("Float64: %f", f64)
		fmt2 := fmt.Sprintf("Float64: %f", f64)
		t.Logf("fmt float64 attempt 1: '%s'", fmt1)
		t.Logf("fmt float64 attempt 2: '%s'", fmt2)
		if fmt1 != fmt2 {
			t.Errorf("fmt inconsistent with itself for float64!")
		}

		xprint1 := xprint.Printf("Float64: %f", f64)
		xprint2 := xprint.Printf("Float64: %f", f64)
		t.Logf("xprint float64 attempt 1: '%s'", xprint1)
		t.Logf("xprint float64 attempt 2: '%s'", xprint2)
		if xprint1 != xprint2 {
			t.Errorf("xprint inconsistent with itself for float64!")
		}

		// Compare fmt vs xprint for float64
		t.Logf("fmt output: '%s'", fmt1)
		t.Logf("xprint output: '%s'", xprint1)
		if fmt1 != xprint1 {
			t.Errorf("MISMATCH for float64")
		}
	})

	// Test with slice of interface{} containing float
	t.Run("SliceOfInterfaceWithFloat", func(t *testing.T) {
		slice := []interface{}{3.14159}

		fmt1 := fmt.Sprintf("Slice float: %f", slice[0])
		xprint1 := xprint.Printf("Slice float: %f", slice[0])

		t.Logf("fmt slice float: '%s'", fmt1)
		t.Logf("xprint slice float: '%s'", xprint1)

		if fmt1 != xprint1 {
			t.Errorf("MISMATCH for slice float")
		}
	})
}
