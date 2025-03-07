package main

import (
	"fmt"

	"gopkg.hlmpn.dev/pkg/go-logger"
	"gopkg.hlmpn.dev/pkg/xprint"
)

func RunNilCheckTests() {
	LogLine()
	logger.Log("Running nil check tests...")
	s := xprint.Printf("%v", nil)
	ff := fmt.Sprintf("%v", nil)
	if s != ff {
		logger.LogErrorf("[NilCheck] ERROR: Mismatch between fmt.Sprintf and xprint.Printf")
	} else {
		logger.LogSuccessf("[NilCheck] Success: Test passed")
	}
	a := xprint.Printf("Hello world")
	b := fmt.Sprintf("Hello world")
	if a != b {
		logger.LogErrorf("[NilCheck] ERROR: Mismatch between fmt.Sprintf and xprint.Printf")
	} else {
		logger.LogSuccessf("[NilCheck] Success: Test passed")
	}
	LogLine()
}
