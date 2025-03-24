package xprint_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	xprint "gopkg.hlmpn.dev/pkg/xprint"
)

// Timer struct to ensure consistent timing
type Timer struct {
	start time.Time
}

type funcDebugResult struct {
	Function string
	Timing   time.Duration
}

// StartTimer initializes and returns a new Timer
func StartTimer() Timer {
	return Timer{start: time.Now()}
}

// Elapsed returns the elapsed time since the timer started
func (t Timer) Elapsed() time.Duration {
	return time.Since(t.start)
}

func BenchmarkLargeStringFormatting(b *testing.B) {
	// Read JSON files
	json1, err := os.ReadFile("/var/web_dev/projects/plexus/run-scripts/detect-run/dummy/5MB.json")
	if err != nil {
		b.Skip("Couldn't read JSON test file: " + err.Error())
		return
	}

	json2, err := os.ReadFile("/var/web_dev/projects/plexus/run-scripts/detect-run/dummy/1MB.json")
	if err != nil {
		b.Skip("Couldn't read JSON test file: " + err.Error())
		return
	}

	aa := string(json1)
	bb := string(json2)

	b.Run("xprint.Printf-string", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = xprint.Printf("%s \n\nHello world %s", aa, bb)
		}
	})

	b.Run("fmt.Sprintf-string", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = fmt.Sprintf("%s \n\nHello world %s", aa, bb)
		}
	})

	// Byte slices benchmark
	b.Run("xprint.Printf-bytes", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = xprint.Printf("%s \n\nHello world %s", json1, json2)
		}
	})

	b.Run("fmt.Sprintf-bytes", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = fmt.Sprintf("%s \n\nHello world %s", json1, json2)
		}
	})
}

func BenchmarkLargeTypesFormatting(b *testing.B) {
	li := TestWrapper{}

	benchmarks := []struct {
		name       string
		xprintFunc func() string
		fmtFunc    func() string
	}{
		{
			name:       "LargeInt",
			xprintFunc: func() string { return xprint.Printf("%d", li.Int()) },
			fmtFunc:    func() string { return fmt.Sprintf("%d", li.Int()) },
		},
		{
			name:       "LargeInt64",
			xprintFunc: func() string { return xprint.Printf("%d", li.Int64()) },
			fmtFunc:    func() string { return fmt.Sprintf("%d", li.Int64()) },
		},
		{
			name:       "LargeInt32",
			xprintFunc: func() string { return xprint.Printf("%d", li.Int32()) },
			fmtFunc:    func() string { return fmt.Sprintf("%d", li.Int32()) },
		},
		{
			name:       "LargeSlice",
			xprintFunc: func() string { return xprint.Printf("%v", li.IntSlice()) },
			fmtFunc:    func() string { return fmt.Sprintf("%v", li.IntSlice()) },
		},
		{
			name:       "LargeStringSlice",
			xprintFunc: func() string { return xprint.Printf("%v", li.StringSlice()) },
			fmtFunc:    func() string { return fmt.Sprintf("%v", li.StringSlice()) },
		},
		{
			name:       "LargeBoolSlice",
			xprintFunc: func() string { return xprint.Printf("%v", li.BoolSlice()) },
			fmtFunc:    func() string { return fmt.Sprintf("%v", li.BoolSlice()) },
		},
		{
			name:       "LargeMixedMap",
			xprintFunc: func() string { return xprint.Printf("%v", li.MixedMap()) },
			fmtFunc:    func() string { return fmt.Sprintf("%v", li.MixedMap()) },
		},
		{
			name:       "LargeStruct",
			xprintFunc: func() string { return xprint.Printf("%v", li.SimpleStruct()) },
			fmtFunc:    func() string { return fmt.Sprintf("%v", li.SimpleStruct()) },
		},
		{
			name:       "DeepNestedStruct",
			xprintFunc: func() string { return xprint.Printf("%v", li.DeeplyNestedStruct()) },
			fmtFunc:    func() string { return fmt.Sprintf("%v", li.DeeplyNestedStruct()) },
		},
	}

	for _, bm := range benchmarks {
		b.Run("xprint-"+bm.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = bm.xprintFunc()
			}
		})

		b.Run("fmt-"+bm.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = bm.fmtFunc()
			}
		})
	}
}

// This test ensures that fmt.Sprintf and xprint.Printf produce the same output
// for various types
func TestFormattingConsistency(t *testing.T) {
	// Test string consistency
	json1, err := os.ReadFile("/var/web_dev/projects/plexus/run-scripts/detect-run/dummy/5MB.json")
	if err != nil {
		t.Skip("Couldn't read JSON test file: " + err.Error())
		return
	}

	json2, err := os.ReadFile("/var/web_dev/projects/plexus/run-scripts/detect-run/dummy/1MB.json")
	if err != nil {
		t.Skip("Couldn't read JSON test file: " + err.Error())
		return
	}

	aa := string(json1)
	bb := string(json2)

	xprintResult := xprint.Printf("%s \n\nHello world %s", aa, bb)
	fmtResult := fmt.Sprintf("%s \n\nHello world %s", aa, bb)

	if xprintResult != fmtResult {
		t.Errorf("[STRING] Output mismatch between fmt.Sprintf and xprint.Printf")
		t.Logf("fmt length: %d, xprint length: %d", len(fmtResult), len(xprintResult))
	}

	// Test byte slice consistency
	xprintByteResult := xprint.Printf("%s \n\nHello world %s", json1, json2)
	fmtByteResult := fmt.Sprintf("%s \n\nHello world %s", json1, json2)

	if xprintByteResult != fmtByteResult {
		t.Errorf("[BYTE] Output mismatch between fmt.Sprintf and xprint.Printf")
		t.Logf("fmt length: %d, xprint length: %d", len(fmtByteResult), len(xprintByteResult))
	}

	// Test with largeints package
	li := TestWrapper{}

	testCases := []struct {
		name       string
		xprintFunc func() string
		fmtFunc    func() string
	}{
		{"Int", func() string { return xprint.Printf("%d", li.Int()) }, func() string { return fmt.Sprintf("%d", li.Int()) }},
		{"Int64", func() string { return xprint.Printf("%d", li.Int64()) }, func() string { return fmt.Sprintf("%d", li.Int64()) }},
		{"Int32", func() string { return xprint.Printf("%d", li.Int32()) }, func() string { return fmt.Sprintf("%d", li.Int32()) }},
		{"IntSlice", func() string { return xprint.Printf("%v", li.IntSlice()) }, func() string { return fmt.Sprintf("%v", li.IntSlice()) }},
		{"StringSlice", func() string { return xprint.Printf("%v", li.StringSlice()) }, func() string { return fmt.Sprintf("%v", li.StringSlice()) }},
		{"BoolSlice", func() string { return xprint.Printf("%v", li.BoolSlice()) }, func() string { return fmt.Sprintf("%v", li.BoolSlice()) }},
		{"MixedMap", func() string { return xprint.Printf("%v", li.MixedMap()) }, func() string { return fmt.Sprintf("%v", li.MixedMap()) }},
		{"SimpleStruct", func() string { return xprint.Printf("%v", li.SimpleStruct()) }, func() string { return fmt.Sprintf("%v", li.SimpleStruct()) }},
		{"DeeplyNestedStruct", func() string { return xprint.Printf("%v", li.DeeplyNestedStruct()) }, func() string { return fmt.Sprintf("%v", li.DeeplyNestedStruct()) }},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			xpOut := tc.xprintFunc()
			fmtOut := tc.fmtFunc()

			if xpOut != fmtOut {
				t.Errorf("[%s] Output mismatch between fmt.Sprintf and xprint.Printf", tc.name)
			}
		})
	}
}
