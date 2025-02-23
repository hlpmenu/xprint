package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"gopkg.hlmpn.dev/pkg/go-logger"
	"gopkg.hlmpn.dev/pkg/xprint/largeints"
	xprint "gopkg.hlmpn.dev/pkg/xprint/t"
)

type StringTetSuite struct {
	Phases []Phase
}
type Phase struct {
	A      string
	B      string
	Middle string
}

func (p *Phase) Bytes() []byte {
	return []byte(p.A + p.Middle + p.B)
}

func (p *Phase) ABytes() []byte {
	return []byte(p.A)
}
func (p *Phase) BBytes() []byte {
	return []byte(p.B)
}
func (p *Phase) MiddleBytes() []byte {
	return []byte(p.Middle)
}

type ResultItem struct {
	Xprint     time.Duration
	Fmt        time.Duration
	StringALen int
	StringBLen int
	MiddleLen  int
	TotalLen   int
}

func GenStrin() string {
	return largeints.LargeString()
}
func GenMiddle() string {
	return largeints.LargeStringN(largeints.RandomNumN(300))
}

var wg = sync.WaitGroup{}

func NewBenchmark() {
	var suite StringTetSuite
	var runs string
	wg.Add(1)
	go func() {
		defer wg.Done()
		suite = StringTetSuite{
			Phases: []Phase{},
		}
		numberOfPhases := largeints.RandomNumN(10)
		for i := range numberOfPhases {
			_ = i
			s := GenStrin()
			s2 := GenStrin()
			mid := GenMiddle()
			suite.Phases = append(suite.Phases, Phase{
				A:      s,
				B:      s2,
				Middle: mid,
			})
		}
		runs = strconv.Itoa(len(suite.Phases))
		LogLine()
		logger.LogSuccessf("Generated %s runs", runs)
	}()
	wg.Wait()
	forceGC()
	var results = []ResultItem{}
	for i, phase := range suite.Phases {
		LogLine()

		var xprintRes string
		var fmtRes string
		xt := xprint.NewTimer()
		ft := xprint.NewTimer()
		if largeints.RandomBool() {
			xt.Start()
			xprintRes = xprint.Printf("%s \n\nHello world %s"+phase.Middle, phase.A, phase.B)
			xt.Stop()
			forceGC()
			ft.Start()
			fmtRes = fmt.Sprintf("%s \n\nHello world %s"+phase.Middle, phase.A, phase.B)
			ft.Stop()
		} else {
			ft.Start()
			fmtRes = fmt.Sprintf("%s \n\nHello world %s"+phase.Middle, phase.B, phase.A)
			ft.Stop()
			forceGC()
			xt.Start()
			xprintRes = xprint.Printf("%s \n\nHello world %s"+phase.Middle, phase.B, phase.A)
			xt.Stop()
		}
		if fmtRes != xprintRes {
			logger.Logf("[Run "+strconv.Itoa(i)+"]len fmt: %d \n len xprint: %d", len(fmtRes), len(xprintRes))
			//nolint:all // Temporarily disabled for string management debugging
			logger.LogErrorf("[Run " + strconv.Itoa(i) + "]ERROR: Output mismatch between fmt.Sprintf and xprint.Printf!")
		}
		xtdebugres := &funcDebugResult{
			Function: "xprint.Printf()",
			Timing:   xt.Duration(),
		}
		fmtdebugres := &funcDebugResult{
			Function: "fmt.Sprintf()",
			Timing:   ft.Duration(),
		}
		logger.LogPurplef("length of output: %d\n", len(xprintRes))
		logTiming("[Run "+strconv.Itoa(i)+"]Printing large JSON (string)", xtdebugres, fmtdebugres)
		results = append(results, ResultItem{
			Xprint:     xt.Duration().Round(time.Millisecond),
			Fmt:        ft.Duration().Round(time.Millisecond),
			StringALen: len(phase.A),
			StringBLen: len(phase.B),
			MiddleLen:  len(phase.Middle),
			TotalLen:   len(xprintRes),
		})
		forceGC()

	}
	// Average results
	avgXprint := time.Duration(0)
	avgFmt := time.Duration(0)
	for _, res := range results {
		avgXprint += res.Xprint
		avgFmt += res.Fmt
	}
	LogLine()
	avgXprint /= time.Duration(len(results))
	avgFmt /= time.Duration(len(results))
	logger.LogSuccessf("Average Xprint: %s \n\n Average Fmt: %s", avgXprint, avgFmt)
	LogLine()
	// Calculate total lengths
	totalLen := 0
	totalStringALen := 0
	totalStringBLen := 0
	for _, res := range results {
		totalLen += res.TotalLen
		totalStringALen += res.StringALen
		totalStringBLen += res.StringBLen
	}

	// Calculate average performance relative to total length
	avgXprintPerTotalLen := float64(avgXprint.Nanoseconds()) / float64(totalLen)
	avgFmtPerTotalLen := float64(avgFmt.Nanoseconds()) / float64(totalLen)
	logger.LogSuccessf("Average Xprint per total length: %.2f ns/byte", avgXprintPerTotalLen)
	logger.LogSuccessf("Average Fmt per total length: %.2f ns/byte", avgFmtPerTotalLen)
	LogLine()

	// Calculate average performance relative to first string argument length
	avgXprintPerStringALen := float64(avgXprint.Nanoseconds()) / float64(totalStringALen)
	avgFmtPerStringALen := float64(avgFmt.Nanoseconds()) / float64(totalStringALen)
	logger.LogSuccessf("Average Xprint per first string length: %.2f ns/byte", avgXprintPerStringALen)
	logger.LogSuccessf("Average Fmt per first string length: %.2f ns/byte", avgFmtPerStringALen)
	LogLine()

	// Calculate average performance relative to last string argument length
	avgXprintPerStringBLen := float64(avgXprint.Nanoseconds()) / float64(totalStringBLen)
	avgFmtPerStringBLen := float64(avgFmt.Nanoseconds()) / float64(totalStringBLen)
	logger.LogSuccessf("Average Xprint per last string length: %.2f ns/byte", avgXprintPerStringBLen)
	logger.LogSuccessf("Average Fmt per last string length: %.2f ns/byte", avgFmtPerStringBLen)
	LogLine()
}

func LogLine() {
	logger.Log("\n=======================================================\n")
}

func BigJSONTest() {
	wg.Add(1)
	var suite StringTetSuite
	var runs string
	go func() {
		defer wg.Done()
		suite = StringTetSuite{
			Phases: []Phase{},
		}
		numberOfPhases := largeints.RandomNumN(10)
		for i := range numberOfPhases {
			jsona, err := os.ReadFile("/var/web_dev/projects/plexus/run-scripts/detect-run/dummy/5MB.json")
			if err != nil {
				logger.LogErrorf("Error reading file: %s", err)
			}
			jsonb, err := os.ReadFile("/var/web_dev/projects/plexus/run-scripts/detect-run/dummy/1MB.json")
			if err != nil {
				logger.LogErrorf("Error reading file: %s", err)
			}

			_ = i
			mid := GenMiddle()
			s := string(jsona) + mid
			s2 := string(jsonb) + mid
			suite.Phases = append(suite.Phases, Phase{
				A:      s,
				B:      s2,
				Middle: mid,
			})
		}
		runs = strconv.Itoa(len(suite.Phases))
		LogLine()
		logger.LogSuccessf("Generated %s runs", runs)
	}()
	wg.Wait()
	forceGC()
	suite.Exec()
	logger.Warn("Running byte version")
	suite.ExecAsBytes()

}

type MixedTypeSuite struct {
	Phases []MixedTypePhase
}
type MixedTypePhase struct {
	A      interface{}
	B      interface{}
	Middle interface{}
}

func MixedTypeTest() {
	li := largeints.TestWrapper{}
	var suite = &MixedTypeSuite{
		Phases: []MixedTypePhase{
			{
				A:      li.Int(),
				B:      li.MixedMap(),
				Middle: li.Int32(),
			},
			{
				A:      li.Int32(),
				B:      li.DeeplyNestedStruct(),
				Middle: li.Int(),
			},
			{
				A:      li.Int64(),
				B:      li.StringSlice(),
				Middle: li.Int32(),
			},
			{
				A:      li.Int32(),
				B:      li.BoolSlice(),
				Middle: li.Int64(),
			},
			{
				A:      li.IntSlice(),
				B:      li.IntSlice(),
				Middle: li.Int64(),
			},
			{
				A:      li.Int64(),
				B:      li.Int64(),
				Middle: li.IntSlice(),
			},
			{
				A:      li.Int64(),
				B:      li.MixedMap(),
				Middle: li.Int64(),
			},
		},
	}
	var results = []ResultItem{}
	for i, phase := range suite.Phases {
		LogLine()

		var xprintRes string
		var fmtRes string
		xt := xprint.NewTimer()
		ft := xprint.NewTimer()
		if largeints.RandomBool() {
			xt.Start()
			xprintRes = xprint.Printf("%s \n\nHello world %s"+"\n\nHello world %s", phase.A, phase.B, phase.Middle)
			xt.Stop()
			forceGC()
			ft.Start()
			fmtRes = fmt.Sprintf("%s \n\nHello world %s"+"\n\nHello world %s", phase.A, phase.B, phase.Middle)
			ft.Stop()
		} else {
			ft.Start()
			fmtRes = fmt.Sprintf("%s \n\nHello world %s"+"\n\nHello world %s", phase.A, phase.B, phase.Middle)
			ft.Stop()
			forceGC()
			xt.Start()
			xprintRes = xprint.Printf("%s \n\nHello world %s"+"\n\nHello world %s", phase.A, phase.B, phase.Middle)
			xt.Stop()
		}
		if fmtRes != xprintRes {
			logger.Logf("[Run "+strconv.Itoa(i)+"]len fmt: %d \n len xprint: %d", len(fmtRes), len(xprintRes))
			logtype(phase.A, phase.B, phase.Middle)
			//nolint:all // Temporarily disabled for string management debugging
			logger.LogErrorf("[Run " + strconv.Itoa(i) + "]ERROR: Output mismatch between fmt.Sprintf and xprint.Printf!")
		}
		xtdebugres := &funcDebugResult{
			Function: "xprint.Printf()",
			Timing:   xt.Duration(),
		}
		fmtdebugres := &funcDebugResult{
			Function: "fmt.Sprintf()",
			Timing:   ft.Duration(),
		}
		logger.LogPurplef("length of output: %d\n", len(xprintRes))
		logTiming("[Run "+strconv.Itoa(i)+"]Printing large JSON (string)", xtdebugres, fmtdebugres)
		results = append(results, ResultItem{
			Xprint:   xt.Duration().Round(time.Millisecond),
			Fmt:      ft.Duration().Round(time.Millisecond),
			TotalLen: len(xprintRes),
		})
		forceGC()
	}
	// Average results
	avgXprint := time.Duration(0)
	avgFmt := time.Duration(0)
	for _, res := range results {
		avgXprint += res.Xprint
		avgFmt += res.Fmt
	}
	LogLine()
	avgXprint /= time.Duration(len(results))
	avgFmt /= time.Duration(len(results))
	logger.LogSuccessf("Average Xprint: %s \n\n Average Fmt: %s", avgXprint, avgFmt)
	LogLine()

}

func (suite *StringTetSuite) Exec() []ResultItem {
	var results = []ResultItem{}
	for i, phase := range suite.Phases {
		LogLine()

		var xprintRes string
		var fmtRes string
		xt := xprint.NewTimer()
		ft := xprint.NewTimer()
		if largeints.RandomBool() {
			xt.Start()
			xprintRes = xprint.Printf("%s \n\nHello world %s"+phase.Middle, phase.A, phase.B)
			xt.Stop()
			forceGC()
			ft.Start()
			fmtRes = fmt.Sprintf("%s \n\nHello world %s"+phase.Middle, phase.A, phase.B)
			ft.Stop()
		} else {
			ft.Start()
			fmtRes = fmt.Sprintf("%s \n\nHello world %s"+phase.Middle, phase.B, phase.A)
			ft.Stop()
			forceGC()
			xt.Start()
			xprintRes = xprint.Printf("%s \n\nHello world %s"+phase.Middle, phase.B, phase.A)
			xt.Stop()
		}
		if fmtRes != xprintRes {
			logger.Logf("[Run "+strconv.Itoa(i)+"]len fmt: %d \n len xprint: %d", len(fmtRes), len(xprintRes))
			//nolint:all // Temporarily disabled for string management debugging
			logger.LogErrorf("[Run " + strconv.Itoa(i) + "]ERROR: Output mismatch between fmt.Sprintf and xprint.Printf!")
		}
		xtdebugres := &funcDebugResult{
			Function: "xprint.Printf()",
			Timing:   xt.Duration(),
		}
		fmtdebugres := &funcDebugResult{
			Function: "fmt.Sprintf()",
			Timing:   ft.Duration(),
		}
		logger.LogPurplef("length of output: %d\n", len(xprintRes))
		logTiming("[Run "+strconv.Itoa(i)+"]Printing large JSON (string)", xtdebugres, fmtdebugres)
		results = append(results, ResultItem{
			Xprint:     xt.Duration().Round(time.Millisecond),
			Fmt:        ft.Duration().Round(time.Millisecond),
			StringALen: len(phase.A),
			StringBLen: len(phase.B),
			MiddleLen:  len(phase.Middle),
			TotalLen:   len(xprintRes),
		})
		forceGC()

	}
	// Average results
	avgXprint := time.Duration(0)
	avgFmt := time.Duration(0)
	for _, res := range results {
		avgXprint += res.Xprint
		avgFmt += res.Fmt
	}
	LogLine()
	avgXprint /= time.Duration(len(results))
	avgFmt /= time.Duration(len(results))
	logger.LogSuccessf("Average Xprint: %s \n\n Average Fmt: %s", avgXprint, avgFmt)
	LogLine()
	// Calculate total lengths
	totalLen := 0
	totalStringALen := 0
	totalStringBLen := 0
	for _, res := range results {
		totalLen += res.TotalLen
		totalStringALen += res.StringALen
		totalStringBLen += res.StringBLen
	}

	// Calculate average performance relative to total length
	avgXprintPerTotalLen := float64(avgXprint.Nanoseconds()) / float64(totalLen)
	avgFmtPerTotalLen := float64(avgFmt.Nanoseconds()) / float64(totalLen)
	logger.LogSuccessf("Average Xprint per total length: %.2f ns/byte", avgXprintPerTotalLen)
	logger.LogSuccessf("Average Fmt per total length: %.2f ns/byte", avgFmtPerTotalLen)
	LogLine()

	// Calculate average performance relative to first string argument length
	avgXprintPerStringALen := float64(avgXprint.Nanoseconds()) / float64(totalStringALen)
	avgFmtPerStringALen := float64(avgFmt.Nanoseconds()) / float64(totalStringALen)
	logger.LogSuccessf("Average Xprint per first string length: %.2f ns/byte", avgXprintPerStringALen)
	logger.LogSuccessf("Average Fmt per first string length: %.2f ns/byte", avgFmtPerStringALen)
	LogLine()

	// Calculate average performance relative to last string argument length
	avgXprintPerStringBLen := float64(avgXprint.Nanoseconds()) / float64(totalStringBLen)
	avgFmtPerStringBLen := float64(avgFmt.Nanoseconds()) / float64(totalStringBLen)
	logger.LogSuccessf("Average Xprint per last string length: %.2f ns/byte", avgXprintPerStringBLen)
	logger.LogSuccessf("Average Fmt per last string length: %.2f ns/byte", avgFmtPerStringBLen)
	LogLine()
	return results
}

func (suite *StringTetSuite) ExecAsBytes() []ResultItem {
	var results = []ResultItem{}
	for i, phase := range suite.Phases {
		LogLine()

		var xprintRes string
		var fmtRes string
		Ab := []byte(phase.A)
		Bb := []byte(phase.B)
		xt := xprint.NewTimer()
		ft := xprint.NewTimer()
		if largeints.RandomBool() {
			xt.Start()
			xprintRes = xprint.Printf("%s \n\nHello world %s"+phase.Middle, Ab, Bb)
			xt.Stop()
			forceGC()
			ft.Start()
			fmtRes = fmt.Sprintf("%s \n\nHello world %s"+phase.Middle, Ab, Bb)
			ft.Stop()
		} else {
			ft.Start()
			fmtRes = fmt.Sprintf("%s \n\nHello world %s"+phase.Middle, Ab, Bb)
			ft.Stop()
			forceGC()
			xt.Start()
			xprintRes = xprint.Printf("%s \n\nHello world %s"+phase.Middle, Ab, Bb)
			xt.Stop()
		}
		if fmtRes != xprintRes {
			logger.Logf("[Run "+strconv.Itoa(i)+"]len fmt: %d \n len xprint: %d", len(fmtRes), len(xprintRes))
			//nolint:all // Temporarily disabled for string management debugging
			logger.LogErrorf("[Run " + strconv.Itoa(i) + "]ERROR: Output mismatch between fmt.Sprintf and xprint.Printf!")
		}
		xtdebugres := &funcDebugResult{
			Function: "xprint.Printf()",
			Timing:   xt.Duration(),
		}
		fmtdebugres := &funcDebugResult{
			Function: "fmt.Sprintf()",
			Timing:   ft.Duration(),
		}
		logger.LogPurplef("length of output: %d\n", len(xprintRes))
		logTiming("[Run "+strconv.Itoa(i)+"]Printing large JSON ([]byte)", xtdebugres, fmtdebugres)
		results = append(results, ResultItem{
			Xprint:     xt.Duration().Round(time.Millisecond),
			Fmt:        ft.Duration().Round(time.Millisecond),
			StringALen: len(phase.ABytes()),
			StringBLen: len(phase.BBytes()),
			MiddleLen:  len(phase.MiddleBytes()),
			TotalLen:   len([]byte(xprintRes)),
		})
		forceGC()

	}
	// Average results
	avgXprint := time.Duration(0)
	avgFmt := time.Duration(0)
	for _, res := range results {
		avgXprint += res.Xprint
		avgFmt += res.Fmt
	}
	LogLine()
	avgXprint /= time.Duration(len(results))
	avgFmt /= time.Duration(len(results))
	logger.LogSuccessf("Average Xprint: %s \n\n Average Fmt: %s", avgXprint, avgFmt)
	LogLine()
	// Calculate total lengths
	totalLen := 0
	totalStringALen := 0
	totalStringBLen := 0
	for _, res := range results {
		totalLen += res.TotalLen
		totalStringALen += res.StringALen
		totalStringBLen += res.StringBLen
	}

	// Calculate average performance relative to total length
	avgXprintPerTotalLen := float64(avgXprint.Nanoseconds()) / float64(totalLen)
	avgFmtPerTotalLen := float64(avgFmt.Nanoseconds()) / float64(totalLen)
	logger.LogSuccessf("Average Xprint per total byte: %.2f ns/byte", avgXprintPerTotalLen)
	logger.LogSuccessf("Average Fmt per total byte: %.2f ns/byte", avgFmtPerTotalLen)
	LogLine()

	// Calculate average performance relative to first string argument length
	avgXprintPerStringALen := float64(avgXprint.Nanoseconds()) / float64(totalStringALen)
	avgFmtPerStringALen := float64(avgFmt.Nanoseconds()) / float64(totalStringALen)
	logger.LogSuccessf("Average Xprint per first byte length: %.2f ns/byte", avgXprintPerStringALen)
	logger.LogSuccessf("Average Fmt per first byte length: %.2f ns/byte", avgFmtPerStringALen)
	LogLine()

	// Calculate average performance relative to last string argument length
	avgXprintPerStringBLen := float64(avgXprint.Nanoseconds()) / float64(totalStringBLen)
	avgFmtPerStringBLen := float64(avgFmt.Nanoseconds()) / float64(totalStringBLen)
	logger.LogSuccessf("Average Xprint per last bytelength: %.2f ns/byte", avgXprintPerStringBLen)
	logger.LogSuccessf("Average Fmt per last byt e length: %.2f ns/byte", avgFmtPerStringBLen)
	LogLine()
	return results
}
