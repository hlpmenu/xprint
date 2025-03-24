package xprint_test

import (
	"fmt"
	"testing"

	"gopkg.hlmpn.dev/pkg/xprint"
)

func TestAppend(t *testing.T) {
	TestAppendString(t)
	TestAppendMultipleArgs(t)
	TestAppendNumericTypes(t)
	TestAppendFloats(t)
	TestAppendComplex(t)
	TestAppendBool(t)
	TestAppendNil(t)
	TestAppendBytes(t)
	TestAppendPointers(t)
	TestAppendSlices(t)
	TestAppendMaps(t)
	TestAppendStruct(t)
	TestAppendInterface(t)
	TestAppendError(t)
	TestAppendEmptyArgs(t)
	TestAppendMultipleTypes(t)
}

func TestAppendf(t *testing.T) {
	TestAppendfString(t)
	TestAppendfMultipleArgs(t)
	TestAppendfNumericTypes(t)
	TestAppendfFloats(t)
	TestAppendfComplex(t)
	TestAppendfBool(t)
	TestAppendfNil(t)
	TestAppendfBytes(t)
	TestAppendfPointers(t)
	TestAppendfSlices(t)
	TestAppendfMaps(t)
	TestAppendfStruct(t)
	TestAppendfInterface(t)
	TestAppendfError(t)
	TestAppendfEmptyArgs(t)
	TestAppendfMultipleTypes(t)
	TestAppendfFormatSpecifiers(t)
}

func ComparePrintfAppend(t *testing.T) {

	x1 := xprint.Printf("%s %d %t %f %v", "string", 42, true, 3.14, []int{1, 2, 3})
	f1 := fmt.Sprintf("%s %d %t %f %v", "string", 42, true, 3.14, []int{1, 2, 3})

	if x1 != f1 {
		t.Errorf("Expected %s, got %s", f1, x1)
	}

	x2 := xprint.Printf("%s %d %t %f %v", "string", 42, true, 3.14, []int{1, 2, 3})
	f2 := fmt.Sprintf("%s %d %t %f %v", "string", 42, true, 3.14, []int{1, 2, 3})

	if x2 != f2 {
		t.Errorf("Expected %s, got %s", f2, x2)
	}

	testCases := []struct {
		name   string
		arg    any
		format string
	}{
		{"int", 42, "%d"},
		{"int8", int8(8), "%d"},
		{"int16", int16(16), "%d"},
		{"int32", int32(32), "%d"},
		{"int64", int64(64), "%d"},
		{"uint", uint(42), "%d"},
		{"uint8", uint8(8), "%d"},
		{"uint16", uint16(16), "%d"},
		{"uint32", uint32(32), "%d"},
		{"uint64", uint64(64), "%d"},
		{"uintptr", uintptr(0xABCD), "%x"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			o := xprint.Printf("Number: ", tc.format, tc.arg)
			fo := fmt.Sprintf("Number: ", tc.format, tc.arg)
			if o != fo {
				t.Errorf("Expected %s, got %s", fo, o)
			}
		})
	}

}
