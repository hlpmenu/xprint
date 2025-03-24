package main

import (
	"errors"
	"fmt"

	"gopkg.hlmpn.dev/pkg/go-logger"
	"gopkg.hlmpn.dev/pkg/xprint"
)

// Using LogErrorfLine to avoid conflict with existing LogLine function
func LogErrorfLine() {
	logger.Log("---------------------------------------")
}

func TestSimpleError() {
	msg := "This is a simple error"
	xf := xprint.Errorf("%s", msg)
	ff := fmt.Errorf("%s", msg)

	if xf.Error() != ff.Error() {
		logger.LogPurplef("fmt output: %s", ff.Error())
		logger.LogOrangef("xprint output: %s", xf.Error())
		logger.LogErrorf("[SimpleError] ERROR: Mismatch between fmt.Errorf and xprint.Errorf")
	} else {
		logger.LogSuccessf("[SimpleError] Success: Test passed")
	}
}

func TestWrappedError() {
	err := errors.New("Original error")
	xf := xprint.Errorf("Wrapped: %w", err)
	ff := fmt.Errorf("Wrapped: %w", err)

	if xf.Error() != ff.Error() {
		logger.LogPurplef("fmt output: %s", ff.Error())
		logger.LogOrangef("xprint output: %s", xf.Error())
		logger.LogErrorf("[WrappedError] ERROR: Mismatch between fmt.Errorf and xprint.Errorf")
	} else {
		logger.LogSuccessf("[WrappedError] Success: Error message test passed")
	}

	// Test unwrapping
	xUnwrapped := errors.Unwrap(xf)
	fUnwrapped := errors.Unwrap(ff)

	if xUnwrapped.Error() != fUnwrapped.Error() {
		logger.LogPurplef("fmt unwrapped: %s", fUnwrapped.Error())
		logger.LogOrangef("xprint unwrapped: %s", xUnwrapped.Error())
		logger.LogErrorf("[WrappedError] ERROR: Mismatch after unwrapping")
	} else {
		logger.LogSuccessf("[WrappedError] Success: Unwrap test passed")
	}
}

func TestMultipleWrappedErrors() {
	err1 := errors.New("First error")
	err2 := errors.New("Second error")
	xf := xprint.Errorf("Multiple errors: %w and %w", err1, err2)
	ff := fmt.Errorf("Multiple errors: %w and %w", err1, err2)

	if xf.Error() != ff.Error() {
		logger.LogPurplef("fmt output: %s", ff.Error())
		logger.LogOrangef("xprint output: %s", xf.Error())
		logger.LogErrorf("[MultipleWrappedErrors] ERROR: Mismatch between fmt.Errorf and xprint.Errorf")
	} else {
		logger.LogSuccessf("[MultipleWrappedErrors] Success: Error message test passed")
	}

	// Test unwrapping multiple errors (Go 1.20+ feature)
	// Access the underlying errors using type assertions on the concrete types
	xfWrapErrors, ok1 := xf.(interface{ Unwrap() []error })
	ffWrapErrors, ok2 := ff.(interface{ Unwrap() []error })

	if !ok1 || !ok2 {
		logger.LogErrorf("[MultipleWrappedErrors] ERROR: Failed to get multiple error unwrapper")
	} else {
		xUnwrapped := xfWrapErrors.Unwrap()
		fUnwrapped := ffWrapErrors.Unwrap()

		if len(xUnwrapped) != len(fUnwrapped) {
			logger.LogPurplef("fmt unwrapped count: %d", len(fUnwrapped))
			logger.LogOrangef("xprint unwrapped count: %d", len(xUnwrapped))
			logger.LogErrorf("[MultipleWrappedErrors] ERROR: Unwrapped error count mismatch")
		} else {
			allMatch := true
			for i := range xUnwrapped {
				if xUnwrapped[i].Error() != fUnwrapped[i].Error() {
					allMatch = false
					logger.LogErrorf("[MultipleWrappedErrors] ERROR: Unwrapped error %d mismatch", i)
				}
			}
			if allMatch {
				logger.LogSuccessf("[MultipleWrappedErrors] Success: Unwrap multiple test passed")
			}
		}
	}
}

func TestMixedFormatting() {
	err := errors.New("error part")
	str := "string part"
	num := 42

	xf := xprint.Errorf("Mixed: %s, %d, %w", str, num, err)
	ff := fmt.Errorf("Mixed: %s, %d, %w", str, num, err)

	if xf.Error() != ff.Error() {
		logger.LogPurplef("fmt output: %s", ff.Error())
		logger.LogOrangef("xprint output: %s", xf.Error())
		logger.LogErrorf("[MixedFormatting] ERROR: Mismatch between fmt.Errorf and xprint.Errorf")
	} else {
		logger.LogSuccessf("[MixedFormatting] Success: Error message test passed")
	}

	// Test unwrapping
	xUnwrapped := errors.Unwrap(xf)
	fUnwrapped := errors.Unwrap(ff)

	if xUnwrapped == nil || fUnwrapped == nil || xUnwrapped.Error() != fUnwrapped.Error() {
		logger.LogPurplef("fmt unwrapped: %v", fUnwrapped)
		logger.LogOrangef("xprint unwrapped: %v", xUnwrapped)
		logger.LogErrorf("[MixedFormatting] ERROR: Mismatch after unwrapping")
	} else {
		logger.LogSuccessf("[MixedFormatting] Success: Unwrap test passed")
	}
}

func TestFastPath() {
	// Test fast path for simple string
	simpleStr := "Simple string test"
	xf1 := xprint.Errorf("%s", simpleStr)
	ff1 := fmt.Errorf("%s", simpleStr)

	if xf1.Error() != ff1.Error() {
		logger.LogErrorf("[FastPath] ERROR: Simple string test failed")
	} else {
		logger.LogSuccessf("[FastPath] Success: Simple string test passed")
	}

	// Test fast path for simple error wrapping
	err := errors.New("Original error for fast path")
	xf2 := xprint.Errorf("Wrapped fast path: %w", err)
	ff2 := fmt.Errorf("Wrapped fast path: %w", err)

	if xf2.Error() != ff2.Error() {
		logger.LogErrorf("[FastPath] ERROR: Simple wrap test failed")
	} else {
		logger.LogSuccessf("[FastPath] Success: Simple wrap test passed")
	}

	// Validate unwrapping works with fast path
	xUnwrapped := errors.Unwrap(xf2)
	fUnwrapped := errors.Unwrap(ff2)

	if xUnwrapped == nil || fUnwrapped == nil || xUnwrapped.Error() != fUnwrapped.Error() {
		logger.LogErrorf("[FastPath] ERROR: Unwrap test failed")
	} else {
		logger.LogSuccessf("[FastPath] Success: Unwrap test passed")
	}
}

func RunErrorfTests() {
	logger.Log("Running Errorf validation tests...")
	LogErrorfLine()

	TestSimpleError()
	LogErrorfLine()

	TestWrappedError()
	LogErrorfLine()

	TestMultipleWrappedErrors()
	LogErrorfLine()

	TestMixedFormatting()
	LogErrorfLine()

	TestFastPath()
	LogErrorfLine()

	// If we get here without logger.LogErrorf exiting the program, all tests passed
	logger.LogSuccessf("All Errorf tests PASSED")
}
