package benchmark_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	xprint "gopkg.hlmpn.dev/pkg/xprint"
)

// ANSI color codes for improved readability.
const (
	colorBlue    = "\033[1;34m"
	colorYellow  = "\033[1;33m"
	colorGreen   = "\033[1;32m"
	colorMagenta = "\033[1;35m"
	colorReset   = "\033[0m"
)

// Initialize random seed.
func init() {
	rand.Seed(time.Now().UnixNano())
}

// Helpers for generating reproducible random data.
func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func randomInt() int {
	return rand.Intn(1000000)
}

func randomUint8() uint8 {
	return uint8(rand.Intn(256))
}

func randomBytes(n int) []byte {
	b := make([]byte, n)
	for i := range b {
		b[i] = byte(rand.Intn(256))
	}
	return b
}

// runSubBenchmark is a helper to run a single sub-test.
// It takes a label (either "fmt" or "xprint"), a formatting function,
// a format string and three arguments (the same pattern is used in all tests).
func runSubBenchmark(iterations int, label string, format string, fn func(string, ...interface{}) string, args ...interface{}) (time.Duration, int64) {
	start := time.Now()
	var totalBytes int64
	// Cycle through the given arguments as needed.
	n := len(args)
	for i := 0; i < iterations; i++ {
		// For a simple three-argument test, we use mod arithmetic.
		result := fn(format, args[i%n], args[(i+1)%n], args[(i+2)%n])
		totalBytes += int64(len(result))
	}
	return time.Since(start), totalBytes
}

// shuffleOrder randomly returns an order of the two strings.
func shuffleOrder(a, b string) []string {
	if rand.Intn(2) == 0 {
		return []string{a, b}
	}
	return []string{b, a}
}

// BenchmarkString tests string formatting performance with multiple sub-tests.
func BenchmarkString(b *testing.B) {
	// Create test data as a slice of interface{}.
	testData := make([]interface{}, 100)
	for i := range testData {
		testData[i] = randomString(100)
	}

	// Decide on a random number of sub-tests (1 to 5).
	numSubTests := 1 + rand.Intn(5)
	fmt.Printf("\n%s=== BenchmarkString: %d sub-tests ===%s\n", colorBlue, numSubTests, colorReset)

	// Possible format variants for strings.
	formatOptions := []string{"%s %s %s", "%v %v %v", "%q %q %q"}

	// Run each sub-test.
	for sub := 0; sub < numSubTests; sub++ {
		// Pick a format at random.
		format := formatOptions[rand.Intn(len(formatOptions))]
		// Random iterations between 500 and 1000.
		iterations := 500 + rand.Intn(501)
		fmt.Printf("\n%s-- Sub-test %d: format = %q, iterations = %d --%s\n", colorYellow, sub+1, format, iterations, colorReset)

		// Randomize the order of the two tests.
		order := shuffleOrder("fmt", "xprint")
		results := make(map[string]struct {
			duration time.Duration
			bytes    int64
		})

		for _, testName := range order {
			if testName == "fmt" {
				// Pass the testData slice directly (no spreading needed because it's already []interface{}).
				d, bts := runSubBenchmark(iterations, "fmt", format, fmt.Sprintf, testData)
				results["fmt"] = struct {
					duration time.Duration
					bytes    int64
				}{d, bts}
				fmt.Printf("%sfmt:%s    %v (%.2f ns/op, %.2f MB/s)\n", colorGreen, colorReset, d,
					float64(d.Nanoseconds())/float64(iterations),
					float64(bts)/d.Seconds()/1024/1024)
			} else { // "xprint"
				d, bts := runSubBenchmark(iterations, "xprint", format, xprint.Printf, testData)
				results["xprint"] = struct {
					duration time.Duration
					bytes    int64
				}{d, bts}
				fmt.Printf("%sxprint:%s %v (%.2f ns/op, %.2f MB/s)\n", colorGreen, colorReset, d,
					float64(d.Nanoseconds())/float64(iterations),
					float64(bts)/d.Seconds()/1024/1024)
			}
		}

		speedup := float64(results["fmt"].duration.Nanoseconds()) / float64(results["xprint"].duration.Nanoseconds())
		fmt.Printf("%sSpeedup (fmt/xprint): %.2fx%s\n", colorMagenta, speedup, colorReset)
	}
}

// BenchmarkInt tests integer formatting performance.
func BenchmarkInt(b *testing.B) {
	// Create test data as a slice of interface{}.
	testData := make([]interface{}, 100)
	for i := range testData {
		testData[i] = randomInt()
	}

	numSubTests := 1 + rand.Intn(5)
	fmt.Printf("\n%s=== BenchmarkInt: %d sub-tests ===%s\n", colorBlue, numSubTests, colorReset)

	// Format options for ints.
	formatOptions := []string{"%d %d %d", "%v %v %v", "%x %x %x"}

	for sub := 0; sub < numSubTests; sub++ {
		format := formatOptions[rand.Intn(len(formatOptions))]
		iterations := 500 + rand.Intn(501)
		fmt.Printf("\n%s-- Sub-test %d: format = %q, iterations = %d --%s\n", colorYellow, sub+1, format, iterations, colorReset)
		order := shuffleOrder("fmt", "xprint")
		results := make(map[string]struct {
			duration time.Duration
			bytes    int64
		})
		for _, testName := range order {
			if testName == "fmt" {
				d, bts := runSubBenchmark(iterations, "fmt", format, fmt.Sprintf, testData)
				results["fmt"] = struct {
					duration time.Duration
					bytes    int64
				}{d, bts}
				fmt.Printf("%sfmt:%s    %v (%.2f ns/op, %.2f MB/s)\n", colorGreen, colorReset, d,
					float64(d.Nanoseconds())/float64(iterations),
					float64(bts)/d.Seconds()/1024/1024)
			} else {
				d, bts := runSubBenchmark(iterations, "xprint", format, xprint.Printf, testData)
				results["xprint"] = struct {
					duration time.Duration
					bytes    int64
				}{d, bts}
				fmt.Printf("%sxprint:%s %v (%.2f ns/op, %.2f MB/s)\n", colorGreen, colorReset, d,
					float64(d.Nanoseconds())/float64(iterations),
					float64(bts)/d.Seconds()/1024/1024)
			}
		}
		speedup := float64(results["fmt"].duration.Nanoseconds()) / float64(results["xprint"].duration.Nanoseconds())
		fmt.Printf("%sSpeedup (fmt/xprint): %.2fx%s\n", colorMagenta, speedup, colorReset)
	}
}

// BenchmarkUint8 tests uint8 formatting performance.
func BenchmarkUint8(b *testing.B) {
	// Create test data as a slice of interface{}.
	testData := make([]interface{}, 100)
	for i := range testData {
		testData[i] = randomUint8()
	}

	numSubTests := 1 + rand.Intn(5)
	fmt.Printf("\n%s=== BenchmarkUint8: %d sub-tests ===%s\n", colorBlue, numSubTests, colorReset)

	// Use int formatting for uint8.
	formatOptions := []string{"%d %d %d", "%v %v %v", "%x %x %x"}

	for sub := 0; sub < numSubTests; sub++ {
		format := formatOptions[rand.Intn(len(formatOptions))]
		iterations := 500 + rand.Intn(501)
		fmt.Printf("\n%s-- Sub-test %d: format = %q, iterations = %d --%s\n", colorYellow, sub+1, format, iterations, colorReset)
		order := shuffleOrder("fmt", "xprint")
		results := make(map[string]struct {
			duration time.Duration
			bytes    int64
		})
		for _, testName := range order {
			if testName == "fmt" {
				d, bts := runSubBenchmark(iterations, "fmt", format, fmt.Sprintf, testData)
				results["fmt"] = struct {
					duration time.Duration
					bytes    int64
				}{d, bts}
				fmt.Printf("%sfmt:%s    %v (%.2f ns/op, %.2f MB/s)\n", colorGreen, colorReset, d,
					float64(d.Nanoseconds())/float64(iterations),
					float64(bts)/d.Seconds()/1024/1024)
			} else {
				d, bts := runSubBenchmark(iterations, "xprint", format, xprint.Printf, testData)
				results["xprint"] = struct {
					duration time.Duration
					bytes    int64
				}{d, bts}
				fmt.Printf("%sxprint:%s %v (%.2f ns/op, %.2f MB/s)\n", colorGreen, colorReset, d,
					float64(d.Nanoseconds())/float64(iterations),
					float64(bts)/d.Seconds()/1024/1024)
			}
		}
		speedup := float64(results["fmt"].duration.Nanoseconds()) / float64(results["xprint"].duration.Nanoseconds())
		fmt.Printf("%sSpeedup (fmt/xprint): %.2fx%s\n", colorMagenta, speedup, colorReset)
	}
}

// BenchmarkBytes tests byte slice formatting performance.
func BenchmarkBytes(b *testing.B) {
	// Create test data as a slice of interface{}.
	testData := make([]interface{}, 100)
	for i := range testData {
		testData[i] = randomBytes(100)
	}

	numSubTests := 1 + rand.Intn(5)
	fmt.Printf("\n%s=== BenchmarkBytes: %d sub-tests ===%s\n", colorBlue, numSubTests, colorReset)

	// Format options for byte slices.
	formatOptions := []string{"%s %s %s", "%v %v %v", "%q %q %q"}

	for sub := 0; sub < numSubTests; sub++ {
		format := formatOptions[rand.Intn(len(formatOptions))]
		iterations := 500 + rand.Intn(501)
		fmt.Printf("\n%s-- Sub-test %d: format = %q, iterations = %d --%s\n", colorYellow, sub+1, format, iterations, colorReset)
		order := shuffleOrder("fmt", "xprint")
		results := make(map[string]struct {
			duration time.Duration
			bytes    int64
		})
		for _, testName := range order {
			if testName == "fmt" {
				d, bts := runSubBenchmark(iterations, "fmt", format, fmt.Sprintf, testData)
				results["fmt"] = struct {
					duration time.Duration
					bytes    int64
				}{d, bts}
				fmt.Printf("%sfmt:%s    %v (%.2f ns/op, %.2f MB/s)\n", colorGreen, colorReset, d,
					float64(d.Nanoseconds())/float64(iterations),
					float64(bts)/d.Seconds()/1024/1024)
			} else {
				d, bts := runSubBenchmark(iterations, "xprint", format, xprint.Printf, testData)
				results["xprint"] = struct {
					duration time.Duration
					bytes    int64
				}{d, bts}
				fmt.Printf("%sxprint:%s %v (%.2f ns/op, %.2f MB/s)\n", colorGreen, colorReset, d,
					float64(d.Nanoseconds())/float64(iterations),
					float64(bts)/d.Seconds()/1024/1024)
			}
		}
		speedup := float64(results["fmt"].duration.Nanoseconds()) / float64(results["xprint"].duration.Nanoseconds())
		fmt.Printf("%sSpeedup (fmt/xprint): %.2fx%s\n", colorMagenta, speedup, colorReset)
	}
}

// BenchmarkMixed tests mixed-type formatting performance.
func BenchmarkMixed(b *testing.B) {
	// Generate test data for each type.
	strData := make([]string, 100)
	intData := make([]int, 100)
	uint8Data := make([]uint8, 100)
	byteData := make([][]byte, 100)
	for i := 0; i < 100; i++ {
		strData[i] = randomString(100)
		intData[i] = randomInt()
		uint8Data[i] = randomUint8()
		byteData[i] = randomBytes(100)
	}

	numSubTests := 1 + rand.Intn(5)
	fmt.Printf("\n%s=== BenchmarkMixed: %d sub-tests ===%s\n", colorBlue, numSubTests, colorReset)

	// Format options for a mixed test.
	formatOptions := []string{
		"%s %d %d %s",
		"%v %v %v %v",
		"%q %d %d %q",
	}

	for sub := 0; sub < numSubTests; sub++ {
		format := formatOptions[rand.Intn(len(formatOptions))]
		iterations := 500 + rand.Intn(501)
		fmt.Printf("\n%s-- Sub-test %d: format = %q, iterations = %d --%s\n", colorYellow, sub+1, format, iterations, colorReset)
		order := shuffleOrder("fmt", "xprint")
		results := make(map[string]struct {
			duration time.Duration
			bytes    int64
		})
		// For mixed, we have four arguments.
		// We'll cycle through each slice independently.
		runMixed := func(fn func(string, ...interface{}) string) (time.Duration, int64) {
			start := time.Now()
			var totalBytes int64
			for i := 0; i < iterations; i++ {
				res := fn(format,
					strData[i%len(strData)],
					intData[i%len(intData)],
					uint8Data[i%len(uint8Data)],
					byteData[i%len(byteData)],
				)
				totalBytes += int64(len(res))
			}
			return time.Since(start), totalBytes
		}
		for _, testName := range order {
			if testName == "fmt" {
				d, bts := runMixed(fmt.Sprintf)
				results["fmt"] = struct {
					duration time.Duration
					bytes    int64
				}{d, bts}
				fmt.Printf("%sfmt:%s    %v (%.2f ns/op, %.2f MB/s)\n", colorGreen, colorReset, d,
					float64(d.Nanoseconds())/float64(iterations),
					float64(bts)/d.Seconds()/1024/1024)
			} else {
				d, bts := runMixed(xprint.Printf)
				results["xprint"] = struct {
					duration time.Duration
					bytes    int64
				}{d, bts}
				fmt.Printf("%sxprint:%s %v (%.2f ns/op, %.2f MB/s)\n", colorGreen, colorReset, d,
					float64(d.Nanoseconds())/float64(iterations),
					float64(bts)/d.Seconds()/1024/1024)
			}
		}
		speedup := float64(results["fmt"].duration.Nanoseconds()) / float64(results["xprint"].duration.Nanoseconds())
		fmt.Printf("%sSpeedup (fmt/xprint): %.2fx%s\n", colorMagenta, speedup, colorReset)
	}
}
