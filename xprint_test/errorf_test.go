package xprint_test

import (
	"errors"
	"fmt"
	"testing"

	"gopkg.hlmpn.dev/pkg/xprint"
)

func TestSimpleError(t *testing.T) {
	msg := "This is a simple error"
	xf := xprint.Errorf("%s", msg)
	ff := fmt.Errorf("%s", msg)

	if xf.Error() != ff.Error() {
		t.Logf("fmt output: %s", ff.Error())
		t.Logf("xprint output: %s", xf.Error())
		t.Errorf("[SimpleError] ERROR: Mismatch between fmt.Errorf and xprint.Errorf")
	} else {
		t.Logf("[SimpleError] Success: Test passed")
	}
}

func TestWrappedError(t *testing.T) {
	err := errors.New("Original error")
	xf := xprint.Errorf("Wrapped: %w", err)
	ff := fmt.Errorf("Wrapped: %w", err)

	if xf.Error() != ff.Error() {
		t.Logf("fmt output: %s", ff.Error())
		t.Logf("xprint output: %s", xf.Error())
		t.Errorf("[WrappedError] ERROR: Mismatch between fmt.Errorf and xprint.Errorf")
	} else {
		t.Logf("[WrappedError] Success: Error message test passed")
	}

	// Test unwrapping
	xUnwrapped := errors.Unwrap(xf)
	fUnwrapped := errors.Unwrap(ff)

	if xUnwrapped.Error() != fUnwrapped.Error() {
		t.Logf("fmt unwrapped: %s", fUnwrapped.Error())
		t.Logf("xprint unwrapped: %s", xUnwrapped.Error())
		t.Errorf("[WrappedError] ERROR: Mismatch after unwrapping")
	} else {
		t.Logf("[WrappedError] Success: Unwrap test passed")
	}
}

func TestMultipleWrappedErrors(t *testing.T) {
	err1 := errors.New("First error")
	err2 := errors.New("Second error")
	xf := xprint.Errorf("Multiple errors: %w and %w", err1, err2)
	ff := fmt.Errorf("Multiple errors: %w and %w", err1, err2)

	if xf.Error() != ff.Error() {
		t.Logf("fmt output: %s", ff.Error())
		t.Logf("xprint output: %s", xf.Error())
		t.Errorf("[MultipleWrappedErrors] ERROR: Mismatch between fmt.Errorf and xprint.Errorf")
	} else {
		t.Logf("[MultipleWrappedErrors] Success: Error message test passed")
	}

	// Test unwrapping multiple errors (Go 1.20+ feature)
	// Access the underlying errors using type assertions on the concrete types
	xfWrapErrors, ok1 := xf.(interface{ Unwrap() []error })
	ffWrapErrors, ok2 := ff.(interface{ Unwrap() []error })

	if !ok1 || !ok2 {
		t.Errorf("[MultipleWrappedErrors] ERROR: Failed to get multiple error unwrapper")
	} else {
		xUnwrapped := xfWrapErrors.Unwrap()
		fUnwrapped := ffWrapErrors.Unwrap()

		if len(xUnwrapped) != len(fUnwrapped) {
			t.Logf("fmt unwrapped count: %d", len(fUnwrapped))
			t.Logf("xprint unwrapped count: %d", len(xUnwrapped))
			t.Errorf("[MultipleWrappedErrors] ERROR: Unwrapped error count mismatch")
		} else {
			allMatch := true
			for i := range xUnwrapped {
				if xUnwrapped[i].Error() != fUnwrapped[i].Error() {
					allMatch = false
					t.Errorf("[MultipleWrappedErrors] ERROR: Unwrapped error %d mismatch", i)
				}
			}
			if allMatch {
				t.Logf("[MultipleWrappedErrors] Success: Unwrap multiple test passed")
			}
		}
	}
}

func TestMixedFormatting(t *testing.T) {
	err := errors.New("error part")
	str := "string part"
	num := 42

	xf := xprint.Errorf("Mixed: %s, %d, %w", str, num, err)
	ff := fmt.Errorf("Mixed: %s, %d, %w", str, num, err)

	if xf.Error() != ff.Error() {
		t.Logf("fmt output: %s", ff.Error())
		t.Logf("xprint output: %s", xf.Error())
		t.Errorf("[MixedFormatting] ERROR: Mismatch between fmt.Errorf and xprint.Errorf")
	} else {
		t.Logf("[MixedFormatting] Success: Error message test passed")
	}

	// Test unwrapping
	xUnwrapped := errors.Unwrap(xf)
	fUnwrapped := errors.Unwrap(ff)

	if xUnwrapped == nil || fUnwrapped == nil || xUnwrapped.Error() != fUnwrapped.Error() {
		t.Logf("fmt unwrapped: %v", fUnwrapped)
		t.Logf("xprint unwrapped: %v", xUnwrapped)
		t.Errorf("[MixedFormatting] ERROR: Mismatch after unwrapping")
	} else {
		t.Logf("[MixedFormatting] Success: Unwrap test passed")
	}
}

func TestFastPath(t *testing.T) {
	// Test fast path for simple string
	simpleStr := "Simple string test"
	xf1 := xprint.Errorf("%s", simpleStr)
	ff1 := fmt.Errorf("%s", simpleStr)

	if xf1.Error() != ff1.Error() {
		t.Errorf("[FastPath] ERROR: Simple string test failed")
	} else {
		t.Logf("[FastPath] Success: Simple string test passed")
	}

	// Test fast path for simple error wrapping
	err := errors.New("Original error for fast path")
	xf2 := xprint.Errorf("Wrapped fast path: %w", err)
	ff2 := fmt.Errorf("Wrapped fast path: %w", err)

	if xf2.Error() != ff2.Error() {
		t.Errorf("[FastPath] ERROR: Simple wrap test failed")
	} else {
		t.Logf("[FastPath] Success: Simple wrap test passed")
	}

	// Validate unwrapping works with fast path
	xUnwrapped := errors.Unwrap(xf2)
	fUnwrapped := errors.Unwrap(ff2)

	if xUnwrapped == nil || fUnwrapped == nil || xUnwrapped.Error() != fUnwrapped.Error() {
		t.Errorf("[FastPath] ERROR: Unwrap test failed")
	} else {
		t.Logf("[FastPath] Success: Unwrap test passed")
	}
}

func TestAllErrorfFunctions(t *testing.T) {
	t.Run("SimpleError", TestSimpleError)
	t.Run("WrappedError", TestWrappedError)
	t.Run("MultipleWrappedErrors", TestMultipleWrappedErrors)
	t.Run("MixedFormatting", TestMixedFormatting)
	t.Run("FastPath", TestFastPath)
}
