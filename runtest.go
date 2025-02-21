package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"runtime"

	"gopkg.hlmpn.dev/pkg/go-logger"
	"gopkg.hlmpn.dev/pkg/xprint/largeints"
	xprint "gopkg.hlmpn.dev/pkg/xprint/t"
)

// Force GC between operations
func forceGC() {
	runtime.GC()
	time.Sleep(10 * time.Millisecond) // Ensure GC completes
}

// Timer struct to ensure consistent timing
type Timer struct {
	start time.Time
}

// StartTimer initializes and returns a new Timer
func StartTimer() Timer {
	return Timer{start: time.Now()}
}

// Elapsed returns the elapsed time since the timer started
func (t Timer) Elapsed() time.Duration {
	return time.Since(t.start)
}

// Log timing for operations
func logTiming(operation, function string, elapsed time.Duration) {
	logger.LogRedf("================================")
	if function == "xprint.Printf()" {
		logger.LogPurplef("Operation: %s\nFunc: %s\nTiming: %v\n", operation, function, elapsed)
	} else if function == "fmt.Sprintf()" {
		logger.LogOrangef("Operation: %s\nFunc: %s\nTiming: %v\n", operation, function, elapsed)
	}
}

func Runtest() {
	log.Println("Starting benchmark...")

	// === Phase 1: Benchmark xprint.Printf() (First) ===

	// Read JSON files
	json1, _ := os.ReadFile("/var/web_dev/projects/plexus/run-scripts/detect-run/dummy/5MB.json")
	json2, _ := os.ReadFile("/var/web_dev/projects/plexus/run-scripts/detect-run/dummy/1MB.json")

	// Force GC
	forceGC()

	// Measure xprint.Printf performance
	timer := StartTimer()
	StringTestXprint := xprint.Printf("%s \n\nHello world %s", string(json1), string(json2))
	elapsedXprint := timer.Elapsed()

	logTiming("Printing large JSON (string)", "xprint.Printf()", elapsedXprint)

	// === Phase 2: Benchmark fmt.Sprintf() ===

	// Read JSON files again (new variables)
	json1b, _ := os.ReadFile("/var/web_dev/projects/plexus/run-scripts/detect-run/dummy/5MB.json")
	json2b, _ := os.ReadFile("/var/web_dev/projects/plexus/run-scripts/detect-run/dummy/1MB.json")

	// Force GC
	forceGC()

	// Measure fmt.Sprintf performance
	timer = StartTimer()
	StringTestFmt := fmt.Sprintf("%s \n\nHello world %s", string(json1b), string(json2b))
	elapsedFmt := timer.Elapsed()

	logTiming("Printing large JSON (string)", "fmt.Sprintf()", elapsedFmt)

	// Validate outputs
	if StringTestFmt != StringTestXprint {
		log.Fatal("ERROR: Output mismatch between fmt.Sprintf and xprint.Printf!")
	}
	return
	// === Phase 3: Benchmark JSON as []byte ===

	// Read JSON files again
	json1c, _ := os.ReadFile("/var/web_dev/projects/plexus/run-scripts/detect-run/dummy/5MB.json")
	json2c, _ := os.ReadFile("/var/web_dev/projects/plexus/run-scripts/detect-run/dummy/1MB.json")

	// Force GC
	forceGC()

	// Measure xprint.Printf performance for []byte
	timer = StartTimer()
	_ = xprint.Printf("%s \n\nHello world %s", json1c, json2c) // Pass as []byte, no conversion
	elapsedXprint = timer.Elapsed()

	logTiming("Printing large JSON ([]byte)", "xprint.Printf()", elapsedXprint)

	// Force GC
	forceGC()

	// Measure fmt.Sprintf performance for []byte
	timer = StartTimer()
	_ = fmt.Sprintf("%s \n\nHello world %s", json1c, json2c) // Pass as []byte, no conversion
	elapsedFmt = timer.Elapsed()

	logTiming("Printing large JSON ([]byte)", "fmt.Sprintf()", elapsedFmt)

	// === Phase 4: Benchmark LargeInts Helper ===
	li := largeints.TestWrapper{}

	benchmarks := []struct {
		name      string
		function  string
		valueFunc func() interface{}
	}{
		{"Printing Large Int", "xprint.Printf()", func() interface{} { return xprint.Printf("%d", li.Int()) }},
		{"Printing Large Int64", "xprint.Printf()", func() interface{} { return xprint.Printf("%d", li.Int64()) }},
		{"Printing Large Int32", "xprint.Printf()", func() interface{} { return xprint.Printf("%d", li.Int32()) }},
		{"Printing Large Slice", "xprint.Printf()", func() interface{} { return xprint.Printf("%v", li.IntSlice()) }},
		{"Printing Large String Slice", "xprint.Printf()", func() interface{} { return xprint.Printf("%v", li.StringSlice()) }},
		{"Printing Large Bool Slice", "xprint.Printf()", func() interface{} { return xprint.Printf("%v", li.BoolSlice()) }},
		{"Printing Large Mixed Map", "xprint.Printf()", func() interface{} { return xprint.Printf("%v", li.MixedMap()) }},
		{"Printing Large Struct", "xprint.Printf()", func() interface{} { return xprint.Printf("%v", li.SimpleStruct()) }},
		{"Printing Deep Nested Struct", "xprint.Printf()", func() interface{} { return xprint.Printf("%v", li.DeeplyNestedStruct()) }},
	}

	// Run benchmarks
	for _, bench := range benchmarks {
		// Force GC
		forceGC()

		// Measure xprint.Printf
		timer = StartTimer()
		outputXprint := bench.valueFunc()
		elapsedXprint = timer.Elapsed()
		logTiming(bench.name, "xprint.Printf()", elapsedXprint)

		// Force GC
		forceGC()

		// Measure fmt.Sprintf
		timer = StartTimer()
		outputFmt := fmt.Sprintf("%v", outputXprint)
		elapsedFmt = timer.Elapsed()
		logTiming(bench.name, "fmt.Sprintf()", elapsedFmt)

		// Validate outputs
		if fmt.Sprintf("%v", outputXprint) != outputFmt {
			log.Fatalf("ERROR: Output mismatch in %s", bench.name)
		}
	}
}
