package xprint_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	xprint "gopkg.hlmpn.dev/pkg/xprint"
)

// BenchmarkLargeJSONString benchmarks formatting of large JSON strings
func BenchmarkLargeJSONString(b *testing.B) {
	// Read JSON files
	json1, err := os.ReadFile("/var/web_dev/projects/plexus/run-scripts/detect-run/dummy/5MB.json")
	if err != nil {
		b.Skip("Could not read 5MB.json: " + err.Error())
		return
	}

	json2, err := os.ReadFile("/var/web_dev/projects/plexus/run-scripts/detect-run/dummy/1MB.json")
	if err != nil {
		b.Skip("Could not read 1MB.json: " + err.Error())
		return
	}

	// Convert to string for string-based test
	jsonStr1 := string(json1)
	jsonStr2 := string(json2)

	b.Run("xprint", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = xprint.Printf("%s \n\nHello world %s", jsonStr1, jsonStr2)
		}
	})

	b.Run("fmt", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = fmt.Sprintf("%s \n\nHello world %s", jsonStr1, jsonStr2)
		}
	})
}

// BenchmarkLargeJSONBytes benchmarks formatting of large JSON as byte slices
func BenchmarkLargeJSONBytes(b *testing.B) {
	// Read JSON files
	json1, err := os.ReadFile("/var/web_dev/projects/plexus/run-scripts/detect-run/dummy/5MB.json")
	if err != nil {
		b.Skip("Could not read 5MB.json: " + err.Error())
		return
	}

	json2, err := os.ReadFile("/var/web_dev/projects/plexus/run-scripts/detect-run/dummy/1MB.json")
	if err != nil {
		b.Skip("Could not read 1MB.json: " + err.Error())
		return
	}

	b.Run("xprint", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = xprint.Printf("%s \n\nHello world %s", json1, json2)
		}
	})

	b.Run("fmt", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = fmt.Sprintf("%s \n\nHello world %s", json1, json2)
		}
	})
}

// TestLargeJSONFormatting verifies that the output is identical between fmt and xprint
func TestLargeJSONFormatting(t *testing.T) {
	// Read JSON files
	json1, err := os.ReadFile("/var/web_dev/projects/plexus/run-scripts/detect-run/dummy/5MB.json")
	if err != nil {
		t.Skip("Could not read 5MB.json: " + err.Error())
		return
	}

	json2, err := os.ReadFile("/var/web_dev/projects/plexus/run-scripts/detect-run/dummy/1MB.json")
	if err != nil {
		t.Skip("Could not read 1MB.json: " + err.Error())
		return
	}

	t.Run("StringFormat", func(t *testing.T) {
		jsonStr1 := string(json1)
		jsonStr2 := string(json2)

		start := time.Now()
		xprintResult := xprint.Printf("%s \n\nHello world %s", jsonStr1, jsonStr2)
		xprintTime := time.Since(start)

		start = time.Now()
		fmtResult := fmt.Sprintf("%s \n\nHello world %s", jsonStr1, jsonStr2)
		fmtTime := time.Since(start)

		t.Logf("xprint time: %v, fmt time: %v", xprintTime, fmtTime)

		if xprintResult != fmtResult {
			t.Errorf("Output mismatch when formatting strings")
			t.Logf("len fmt: %d, len xprint: %d", len(fmtResult), len(xprintResult))
		} else {
			t.Logf("String formatting output matches")
		}
	})

	t.Run("ByteFormat", func(t *testing.T) {
		start := time.Now()
		xprintResult := xprint.Printf("%s \n\nHello world %s", json1, json2)
		xprintTime := time.Since(start)

		start = time.Now()
		fmtResult := fmt.Sprintf("%s \n\nHello world %s", json1, json2)
		fmtTime := time.Since(start)

		t.Logf("xprint time: %v, fmt time: %v", xprintTime, fmtTime)

		if xprintResult != fmtResult {
			t.Errorf("Output mismatch when formatting byte slices")
			t.Logf("len fmt: %d, len xprint: %d", len(fmtResult), len(xprintResult))
		} else {
			t.Logf("Byte slice formatting output matches")
		}
	})
}

// BenchmarkLargeInts benchmarks various integer and complex data structure formatting
func BenchmarkLargeInts(b *testing.B) {
	li := TestWrapper{}

	benchmarks := []struct {
		name     string
		xprintFn func() string
		fmtFn    func() string
	}{
		{
			name:     "LargeInt",
			xprintFn: func() string { return xprint.Printf("%d", li.Int()) },
			fmtFn:    func() string { return fmt.Sprintf("%d", li.Int()) },
		},
		{
			name:     "LargeInt64",
			xprintFn: func() string { return xprint.Printf("%d", li.Int64()) },
			fmtFn:    func() string { return fmt.Sprintf("%d", li.Int64()) },
		},
		{
			name:     "LargeInt32",
			xprintFn: func() string { return xprint.Printf("%d", li.Int32()) },
			fmtFn:    func() string { return fmt.Sprintf("%d", li.Int32()) },
		},
		{
			name:     "LargeSlice",
			xprintFn: func() string { return xprint.Printf("%v", li.IntSlice()) },
			fmtFn:    func() string { return fmt.Sprintf("%v", li.IntSlice()) },
		},
		{
			name:     "LargeStringSlice",
			xprintFn: func() string { return xprint.Printf("%v", li.StringSlice()) },
			fmtFn:    func() string { return fmt.Sprintf("%v", li.StringSlice()) },
		},
		{
			name:     "LargeBoolSlice",
			xprintFn: func() string { return xprint.Printf("%v", li.BoolSlice()) },
			fmtFn:    func() string { return fmt.Sprintf("%v", li.BoolSlice()) },
		},
		{
			name:     "LargeMixedMap",
			xprintFn: func() string { return xprint.Printf("%v", li.MixedMap()) },
			fmtFn:    func() string { return fmt.Sprintf("%v", li.MixedMap()) },
		},
		{
			name:     "LargeStruct",
			xprintFn: func() string { return xprint.Printf("%v", li.SimpleStruct()) },
			fmtFn:    func() string { return fmt.Sprintf("%v", li.SimpleStruct()) },
		},
		{
			name:     "DeepNestedStruct",
			xprintFn: func() string { return xprint.Printf("%v", li.DeeplyNestedStruct()) },
			fmtFn:    func() string { return fmt.Sprintf("%v", li.DeeplyNestedStruct()) },
		},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name+"-xprint", func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = bm.xprintFn()
			}
		})

		b.Run(bm.name+"-fmt", func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = bm.fmtFn()
			}
		})
	}
}

// TestLargeIntsFormatting verifies that output is identical for complex data structures
func TestLargeIntsFormatting(t *testing.T) {
	li := TestWrapper{}

	testCases := []struct {
		name     string
		xprintFn func() string
		fmtFn    func() string
	}{
		{
			name:     "Int",
			xprintFn: func() string { return xprint.Printf("%d", li.Int()) },
			fmtFn:    func() string { return fmt.Sprintf("%d", li.Int()) },
		},
		{
			name:     "Int64",
			xprintFn: func() string { return xprint.Printf("%d", li.Int64()) },
			fmtFn:    func() string { return fmt.Sprintf("%d", li.Int64()) },
		},
		{
			name:     "Int32",
			xprintFn: func() string { return xprint.Printf("%d", li.Int32()) },
			fmtFn:    func() string { return fmt.Sprintf("%d", li.Int32()) },
		},
		{
			name:     "IntSlice",
			xprintFn: func() string { return xprint.Printf("%v", li.IntSlice()) },
			fmtFn:    func() string { return fmt.Sprintf("%v", li.IntSlice()) },
		},
		{
			name:     "StringSlice",
			xprintFn: func() string { return xprint.Printf("%v", li.StringSlice()) },
			fmtFn:    func() string { return fmt.Sprintf("%v", li.StringSlice()) },
		},
		{
			name:     "MixedMap",
			xprintFn: func() string { return xprint.Printf("%v", li.MixedMap()) },
			fmtFn:    func() string { return fmt.Sprintf("%v", li.MixedMap()) },
		},
		{
			name:     "SimpleStruct",
			xprintFn: func() string { return xprint.Printf("%v", li.SimpleStruct()) },
			fmtFn:    func() string { return fmt.Sprintf("%v", li.SimpleStruct()) },
		},
		{
			name:     "DeeplyNestedStruct",
			xprintFn: func() string { return xprint.Printf("%v", li.DeeplyNestedStruct()) },
			fmtFn:    func() string { return fmt.Sprintf("%v", li.DeeplyNestedStruct()) },
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			xprintResult := tc.xprintFn()
			fmtResult := tc.fmtFn()

			if xprintResult != fmtResult {
				t.Errorf("Output mismatch for %s test", tc.name)
			}
		})
	}
}
