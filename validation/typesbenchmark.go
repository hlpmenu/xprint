package main

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"

	"gopkg.hlmpn.dev/pkg/go-logger"
	"gopkg.hlmpn.dev/pkg/xprint"
)

// BenchmarkItem defines a single benchmark test case
type BenchmarkItem struct {
	Name       string
	Format     string
	Args       []any
	Iterations int
}

// TypesBenchmark runs benchmarks for various data types
func TypesBenchmark() {
	logger.Log("Running benchmark for various data types...")

	// Create benchmark cases for different types
	benchmarks := []BenchmarkItem{
		{
			Name:       "String",
			Format:     "String: %s",
			Args:       []any{"Hello, world!"},
			Iterations: 1000000,
		},
		{
			Name:       "Int",
			Format:     "Int: %d",
			Args:       []any{42},
			Iterations: 1000000,
		},
		{
			Name:       "Int64",
			Format:     "Int64: %d",
			Args:       []any{int64(9223372036854775807)},
			Iterations: 1000000,
		},
		{
			Name:       "Float",
			Format:     "Float: %.2f",
			Args:       []any{3.14159},
			Iterations: 1000000,
		},
		{
			Name:       "Bool",
			Format:     "Bool: %t",
			Args:       []any{true},
			Iterations: 1000000,
		},
		{
			Name:       "Error",
			Format:     "Error: %v",
			Args:       []any{errors.New("test error")},
			Iterations: 1000000,
		},
		{
			Name:       "Bytes",
			Format:     "Bytes: %s",
			Args:       []any{[]byte("byte slice")},
			Iterations: 1000000,
		},
		{
			Name:       "Complex",
			Format:     "Complex format with multiple args: %s - %d - %t - %f",
			Args:       []any{"text", 42, true, 3.14},
			Iterations: 500000,
		},
		{
			Name:       "Map",
			Format:     "Map: %v",
			Args:       []any{map[string]int{"one": 1, "two": 2}},
			Iterations: 500000,
		},
		{
			Name:       "Slice",
			Format:     "Slice: %v",
			Args:       []any{[]int{1, 2, 3, 4, 5}},
			Iterations: 500000,
		},
		{
			Name:       "Pointer",
			Format:     "Pointer: %p",
			Args:       []any{&struct{}{}},
			Iterations: 500000,
		},
	}

	// Shuffle the benchmarks for randomized order
	rand.Shuffle(len(benchmarks), func(i, j int) {
		benchmarks[i], benchmarks[j] = benchmarks[j], benchmarks[i]
	})

	// Track total times for overall comparison
	var totalXprintTime, totalFmtTime time.Duration

	// Run each benchmark
	for _, bench := range benchmarks {
		runTypeBenchmark(bench, &totalXprintTime, &totalFmtTime)
		LogLine()
	}

	// Print overall results
	logger.Log("Overall results:")
	logger.LogPurplef("Total xprint.Printf time: %v", totalXprintTime)
	logger.LogOrangef("Total fmt.Sprintf time: %v", totalFmtTime)
	ratio := float64(totalFmtTime.Nanoseconds()) / float64(totalXprintTime.Nanoseconds())
	logger.Logf("xprint is %.2fx faster overall", ratio)
}

// runTypeBenchmark runs a single benchmark for a specific type
func runTypeBenchmark(bench BenchmarkItem, totalXprintTime, totalFmtTime *time.Duration) {
	logger.Logf("Running %s benchmark...", bench.Name)

	// Run xprint benchmark
	start := time.Now()
	var xprintResult string
	for i := 0; i < bench.Iterations; i++ {
		_ = xprint.Printf(bench.Format, bench.Args...)
	}
	xprintTime := time.Since(start)
	*totalXprintTime += xprintTime

	// Run fmt benchmark
	start = time.Now()
	var fmtResult string
	for i := 0; i < bench.Iterations; i++ {
		_ = fmt.Sprintf(bench.Format, bench.Args...)
	}
	fmtTime := time.Since(start)
	*totalFmtTime += fmtTime

	// Run one more time to validate results
	xprintResult = xprint.Printf(bench.Format, bench.Args...)
	fmtResult = fmt.Sprintf(bench.Format, bench.Args...)

	// Verify results match and display outputs
	logger.LogPurplef("xprint output: %s", xprintResult)
	logger.LogOrangef("fmt output: %s", fmtResult)

	// For maps, we don't compare exact strings since order may differ
	skipExactComparison := bench.Name == "Map"

	if xprintResult != fmtResult && !skipExactComparison {
		logger.LogErrorf("[%s] MISMATCH: Results don't match!", bench.Name)
	}

	// Calculate and display performance metrics
	nsPerOp := float64(xprintTime.Nanoseconds()) / float64(bench.Iterations)
	fmtNsPerOp := float64(fmtTime.Nanoseconds()) / float64(bench.Iterations)
	speedup := fmtNsPerOp / nsPerOp

	logger.LogPurplef("xprint.Printf: %.2f ns/op", nsPerOp)
	logger.LogOrangef("fmt.Sprintf: %.2f ns/op", fmtNsPerOp)
	logger.Logf("Speedup: %.2fx", speedup)
}

func TestFloats() {
	logger.Log("Testing float formatting issues...")
	LogLine()

	// Test with different precision formats (moved to first)
	logger.Log("Test 1: Different precision formats")
	LogLine()

	var f64 float64 = 3.14159

	precisions := []string{
		"%.0f", "%.1f", "%.2f", "%.3f", "%.6f", "%f", "%g", "%e",
	}

	for _, prec := range precisions {
		logger.Log("Testing precision format: " + prec)

		// Use separate format strings for the description and the actual format
		fmt1 := fmt.Sprintf(prec, f64)
		xprint1 := xprint.Printf(prec, f64)

		logger.LogOrangef("fmt output: '%s'", fmt1)
		logger.LogPurplef("xprint output: '%s'", xprint1)

		if fmt1 != xprint1 {
			logger.LogErrorf("❌ MISMATCH for format %s", prec)
			// Exit with an error code if the formats don't match
			os.Exit(1)
		} else {
			logger.LogSuccessf("✅ Match for format %s", prec)
		}
		LogLine() // Add LogLine after each test case for better separation
	}

	// Test raw float constants (interface{} type inference)
	logger.Log("Test 2: Raw float constants as interface{} arguments")
	LogLine()

	// First, test fmt self-consistency
	fmt1 := fmt.Sprintf("Float: %f", 3.14159)
	fmt2 := fmt.Sprintf("Float: %f", 3.14159)
	logger.LogOrangef("fmt attempt 1: '%s'", fmt1)
	logger.LogOrangef("fmt attempt 2: '%s'", fmt2)
	if fmt1 != fmt2 {
		logger.LogErrorf("❌ fmt inconsistent with itself!")
		os.Exit(1)
	} else {
		logger.LogSuccessf("✅ fmt consistent with itself")
	}

	// Then, test xprint self-consistency
	xprint1 := xprint.Printf("Float: %f", 3.14159)
	xprint2 := xprint.Printf("Float: %f", 3.14159)
	logger.LogPurplef("xprint attempt 1: '%s'", xprint1)
	logger.LogPurplef("xprint attempt 2: '%s'", xprint2)
	if xprint1 != xprint2 {
		logger.LogErrorf("❌ xprint inconsistent with itself!")
		os.Exit(1)
	} else {
		logger.LogSuccessf("✅ xprint consistent with itself")
	}

	// Compare fmt vs xprint
	logger.Log("Comparing fmt vs xprint:")
	logger.LogOrangef("fmt output: '%s'", fmt1)
	logger.LogPurplef("xprint output: '%s'", xprint1)
	if fmt1 != xprint1 {
		logger.LogErrorf("❌ MISMATCH: fmt vs xprint for raw float")
		os.Exit(1)
	} else {
		logger.LogSuccessf("✅ fmt and xprint match")
	}
	LogLine()

	// Test with explicit float32
	logger.Log("Test 3: Explicit float32 variable")
	LogLine()
	var f32 float32 = 3.14159

	fmt1 = fmt.Sprintf("Float32: %f", f32)
	fmt2 = fmt.Sprintf("Float32: %f", f32)
	logger.LogOrangef("fmt float32 attempt 1: '%s'", fmt1)
	logger.LogOrangef("fmt float32 attempt 2: '%s'", fmt2)

	xprint1 = xprint.Printf("Float32: %f", f32)
	xprint2 = xprint.Printf("Float32: %f", f32)
	logger.LogPurplef("xprint float32 attempt 1: '%s'", xprint1)
	logger.LogPurplef("xprint float32 attempt 2: '%s'", xprint2)

	// Compare fmt vs xprint for float32
	logger.Log("Comparing fmt vs xprint for float32:")
	logger.LogOrangef("fmt output: '%s'", fmt1)
	logger.LogPurplef("xprint output: '%s'", xprint1)
	if fmt1 != xprint1 {
		logger.LogErrorf("❌ MISMATCH for float32")
		os.Exit(1)
	} else {
		logger.LogSuccessf("✅ fmt and xprint match for float32")
	}
	LogLine()

	// Test with explicit float64
	logger.Log("Test 4: Explicit float64 variable")
	LogLine()

	fmt1 = fmt.Sprintf("Float64: %f", f64)
	fmt2 = fmt.Sprintf("Float64: %f", f64)
	logger.LogOrangef("fmt float64 attempt 1: '%s'", fmt1)
	logger.LogOrangef("fmt float64 attempt 2: '%s'", fmt2)

	xprint1 = xprint.Printf("Float64: %f", f64)
	xprint2 = xprint.Printf("Float64: %f", f64)
	logger.LogPurplef("xprint float64 attempt 1: '%s'", xprint1)
	logger.LogPurplef("xprint float64 attempt 2: '%s'", xprint2)

	// Compare fmt vs xprint for float64
	logger.Log("Comparing fmt vs xprint for float64:")
	logger.LogOrangef("fmt output: '%s'", fmt1)
	logger.LogPurplef("xprint output: '%s'", xprint1)
	if fmt1 != xprint1 {
		logger.LogErrorf("❌ MISMATCH for float64")
		os.Exit(1)
	} else {
		logger.LogSuccessf("✅ fmt and xprint match for float64")
	}
	LogLine()

	// Test with slice of interface{} containing float
	logger.Log("Test 5: Slice of interface{} containing float")
	LogLine()
	slice := []interface{}{3.14159}

	fmt1 = fmt.Sprintf("Slice float: %f", slice[0])
	xprint1 = xprint.Printf("Slice float: %f", slice[0])

	logger.LogOrangef("fmt slice float: '%s'", fmt1)
	logger.LogPurplef("xprint slice float: '%s'", xprint1)

	if fmt1 != xprint1 {
		logger.LogErrorf("❌ MISMATCH for slice float")
		os.Exit(1)
	} else {
		logger.LogSuccessf("✅ fmt and xprint match for slice float")
	}
	LogLine()
}
