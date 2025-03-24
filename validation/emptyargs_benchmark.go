package main

import (
	"fmt"
	"time"

	"gopkg.hlmpn.dev/pkg/go-logger"
	"gopkg.hlmpn.dev/pkg/xprint"
)

// EmptyArgsBenchmark tests the performance of Printf with no arguments
func EmptyArgsBenchmark() {
	logger.Log("Running empty args benchmark...")

	const iterations = 1000000
	const testStr = "This is a test string with no formatting"

	// Benchmark xprint.Printf with no args
	start := time.Now()
	for i := 0; i < iterations; i++ {
		_ = xprint.Printf(testStr)
	}
	xprintDuration := time.Since(start)

	// Benchmark fmt.Sprintf with no args
	start = time.Now()
	for i := 0; i < iterations; i++ {
		_ = fmt.Sprintf(testStr)
	}
	fmtDuration := time.Since(start)

	logger.LogSuccessf("Empty args benchmark (ns/op):")
	logger.LogSuccessf("xprint.Printf: %.2f ns/op", float64(xprintDuration.Nanoseconds())/float64(iterations))
	logger.LogSuccessf("fmt.Sprintf: %.2f ns/op", float64(fmtDuration.Nanoseconds())/float64(iterations))

	// With a % sign that would normally need to be processed
	const testStrWithPercent = "This has a % sign but no args"

	// Benchmark xprint.Printf with no args but with % sign
	start = time.Now()
	for i := 0; i < iterations; i++ {
		_ = xprint.Printf(testStrWithPercent)
	}
	xprintPercentDuration := time.Since(start)

	// Benchmark fmt.Sprintf with no args but with % sign
	start = time.Now()
	for i := 0; i < iterations; i++ {
		_ = fmt.Sprintf(testStrWithPercent)
	}
	fmtPercentDuration := time.Since(start)

	// Use non-formatted log methods to avoid the % issue
	logger.LogSuccess("Empty args with % benchmark (ns/op):")
	logger.LogSuccessf("xprint.Printf: %.2f ns/op", float64(xprintPercentDuration.Nanoseconds())/float64(iterations))
	logger.LogSuccessf("fmt.Sprintf: %.2f ns/op", float64(fmtPercentDuration.Nanoseconds())/float64(iterations))
}
