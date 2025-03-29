package xprint_test

import (
	"bytes"
	"fmt"
	"testing"

	"gopkg.hlmpn.dev/pkg/xprint"
)

func TestAppendString(t *testing.T) {
	o := xprint.Append([]byte("Hello "), "World!")
	fo := fmt.Append([]byte("Hello "), "World!")
	if !bytes.Equal(o, fo) {
		t.Errorf("Expected %s, got %s", fo, o)
	}
}

func TestAppendMultipleArgs(t *testing.T) {
	o := xprint.Append([]byte("Results: "), "Text", 42, true)
	fo := fmt.Append([]byte("Results: "), "Text", 42, true)
	if !bytes.Equal(o, fo) {
		t.Errorf("Expected %s, got %s", fo, o)
	}
}

func TestAppendNumericTypes(t *testing.T) {
	testCases := []struct {
		name string
		arg  any
	}{
		{"int", 42},
		{"int8", int8(8)},
		{"int16", int16(16)},
		{"int32", int32(32)},
		{"int64", int64(64)},
		{"uint", uint(42)},
		{"uint8", uint8(8)},
		{"uint16", uint16(16)},
		{"uint32", uint32(32)},
		{"uint64", uint64(64)},
		{"uintptr", uintptr(0xABCD)},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			o := xprint.Append([]byte("Number: "), tc.arg)
			fo := fmt.Append([]byte("Number: "), tc.arg)
			if !bytes.Equal(o, fo) {
				t.Errorf("Expected %s, got %s", fo, o)
			}
		})
	}
}

func TestAppendFloats(t *testing.T) {
	testCases := []struct {
		name string
		arg  any
	}{
		{"float32", float32(3.14159)},
		{"float64", float64(2.71828)},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			o := xprint.Append([]byte("Float: "), tc.arg)
			fo := fmt.Append([]byte("Float: "), tc.arg)
			if !bytes.Equal(o, fo) {
				t.Errorf("Expected %s, got %s", fo, o)
			}
		})
	}
}

func TestAppendComplex(t *testing.T) {
	testCases := []struct {
		name string
		arg  any
	}{
		{"complex64", complex64(complex(1, 2))},
		{"complex128", complex128(complex(3, 4))},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			o := xprint.Append([]byte("Complex: "), tc.arg)
			fo := fmt.Append([]byte("Complex: "), tc.arg)
			if !bytes.Equal(o, fo) {
				t.Errorf("Expected %s, got %s", fo, o)
			}
		})
	}
}

func TestAppendBool(t *testing.T) {
	testCases := []struct {
		name string
		arg  bool
	}{
		{"true", true},
		{"false", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			o := xprint.Append([]byte("Bool: "), tc.arg)
			fo := fmt.Append([]byte("Bool: "), tc.arg)
			if !bytes.Equal(o, fo) {
				t.Errorf("Expected %s, got %s", fo, o)
			}
		})
	}
}

func TestAppendNil(t *testing.T) {
	o := xprint.Append([]byte("Nil: "), nil)
	fo := fmt.Append([]byte("Nil: "), nil)
	if !bytes.Equal(o, fo) {
		t.Errorf("Expected %s, got %s", fo, o)
	}
}

func TestAppendBytes(t *testing.T) {
	data := []byte("byte slice")
	o := xprint.Append([]byte("Bytes: "), data)
	fo := fmt.Append([]byte("Bytes: "), data)
	if !bytes.Equal(o, fo) {
		t.Errorf("Expected %s, got %s", fo, o)
	}
}

func TestAppendPointers(t *testing.T) {
	// Regular pointer
	str := "pointer test"
	strPtr := &str

	// Nil pointer
	var nilPtr *string

	t.Run("regular pointer", func(t *testing.T) {
		o := xprint.Append([]byte("Pointer: "), strPtr)
		fo := fmt.Append([]byte("Pointer: "), strPtr)
		if !bytes.Equal(o, fo) {
			t.Errorf("Expected %s, got %s", fo, o)
		}
	})

	t.Run("nil pointer", func(t *testing.T) {
		o := xprint.Append([]byte("Nil pointer: "), nilPtr)
		fo := fmt.Append([]byte("Nil pointer: "), nilPtr)
		if !bytes.Equal(o, fo) {
			t.Errorf("Expected %s, got %s", fo, o)
		}
	})
}

func TestAppendSlices(t *testing.T) {
	intSlice := []int{1, 2, 3, 4, 5}
	stringSlice := []string{"a", "b", "c"}
	emptySlice := []float64{}

	testCases := []struct {
		name string
		arg  any
	}{
		{"int slice", intSlice},
		{"string slice", stringSlice},
		{"empty slice", emptySlice},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			o := xprint.Append([]byte("Slice: "), tc.arg)
			fo := fmt.Append([]byte("Slice: "), tc.arg)
			if !bytes.Equal(o, fo) {
				t.Errorf("Expected %s, got %s", fo, o)
			}
		})
	}
}

func TestAppendMaps(t *testing.T) {
	stringMap := map[string]string{"key1": "value1", "key2": "value2"}
	intMap := map[string]int{"one": 1, "two": 2}
	emptyMap := map[int]bool{}

	t.Run("string map", func(t *testing.T) {
		o := xprint.Append([]byte("Map: "), stringMap)
		fo := fmt.Append([]byte("Map: "), stringMap)
		// For maps, we don't compare exact output as ordering may differ
		if len(o) != len(fo) {
			t.Errorf("Length mismatch for string map. Expected %d, got %d", len(fo), len(o))
		}
	})

	t.Run("int map", func(t *testing.T) {
		o := xprint.Append([]byte("Map: "), intMap)
		fo := fmt.Append([]byte("Map: "), intMap)
		if len(o) != len(fo) {
			t.Errorf("Length mismatch for int map. Expected %d, got %d", len(fo), len(o))
		}
	})

	t.Run("empty map", func(t *testing.T) {
		o := xprint.Append([]byte("Map: "), emptyMap)
		fo := fmt.Append([]byte("Map: "), emptyMap)
		if !bytes.Equal(o, fo) {
			t.Errorf("Expected %s, got %s", fo, o)
		}
	})
}

func TestAppendStruct(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	person := Person{Name: "John", Age: 30}

	o := xprint.Append([]byte("Struct: "), person)
	fo := fmt.Append([]byte("Struct: "), person)
	if !bytes.Equal(o, fo) {
		t.Errorf("Expected %s, got %s", fo, o)
	}
}

func TestAppendInterface(t *testing.T) {
	var i any = "interface value"

	o := xprint.Append([]byte("Interface: "), i)
	fo := fmt.Append([]byte("Interface: "), i)
	if !bytes.Equal(o, fo) {
		t.Errorf("Expected %s, got %s", fo, o)
	}
}

func TestAppendError(t *testing.T) {
	err := fmt.Errorf("test error") //nolint:all //

	o := xprint.Append([]byte("Error: "), err)
	fo := fmt.Append([]byte("Error: "), err)
	if !bytes.Equal(o, fo) {
		t.Errorf("Expected %s, got %s", fo, o)
	}
}

func TestAppendEmptyArgs(t *testing.T) {
	o := xprint.Append([]byte("No args: "))
	fo := fmt.Append([]byte("No args: "))
	if !bytes.Equal(o, fo) {
		t.Errorf("Expected %s, got %s", fo, o)
	}
}

func TestAppendMultipleTypes(t *testing.T) {
	o := xprint.Append([]byte("Mixed: "), "string", 42, true, 3.14, []int{1, 2, 3})
	fo := fmt.Append([]byte("Mixed: "), "string", 42, true, 3.14, []int{1, 2, 3})
	if !bytes.Equal(o, fo) {
		t.Errorf("Expected %s, got %s", fo, o)
	}
}

// Appendf tests

func TestAppendfString(t *testing.T) {
	o := xprint.Appendf([]byte("Hello "), "%s", "World!")
	fo := fmt.Appendf([]byte("Hello "), "%s", "World!")
	if !bytes.Equal(o, fo) {
		t.Errorf("Expected %s, got %s", fo, o)
	}
}

func TestAppendfMultipleArgs(t *testing.T) {
	o := xprint.Appendf([]byte("Results: "), "%s %d %t", "Text", 42, true)
	fo := fmt.Appendf([]byte("Results: "), "%s %d %t", "Text", 42, true)
	if !bytes.Equal(o, fo) {
		t.Errorf("Expected %s, got %s", fo, o)
	}
}

func TestAppendfNumericTypes(t *testing.T) {
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
			o := xprint.Appendf([]byte("Number: "), tc.format, tc.arg)
			fo := fmt.Appendf([]byte("Number: "), tc.format, tc.arg)
			if !bytes.Equal(o, fo) {
				t.Errorf("Expected %s, got %s", fo, o)
			}
		})
	}
}

func TestAppendfFloats(t *testing.T) {
	testCases := []struct {
		name   string
		arg    any
		format string
	}{
		{"float32", float32(3.14159), "%.5f"},
		{"float64", float64(2.71828), "%.5f"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			o := xprint.Appendf([]byte("Float: "), tc.format, tc.arg)
			fo := fmt.Appendf([]byte("Float: "), tc.format, tc.arg)
			if !bytes.Equal(o, fo) {
				t.Errorf("Expected %s, got %s", fo, o)
			}
		})
	}
}

func TestAppendfComplex(t *testing.T) {
	testCases := []struct {
		name   string
		arg    any
		format string
	}{
		{"complex64", complex64(complex(1, 2)), "%v"},
		{"complex128", complex128(complex(3, 4)), "%v"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			o := xprint.Appendf([]byte("Complex: "), tc.format, tc.arg)
			fo := fmt.Appendf([]byte("Complex: "), tc.format, tc.arg)
			if !bytes.Equal(o, fo) {
				t.Errorf("Expected %s, got %s", fo, o)
			}
		})
	}
}

func TestAppendfBool(t *testing.T) {
	testCases := []struct {
		name string
		arg  bool
	}{
		{"true", true},
		{"false", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			o := xprint.Appendf([]byte("Bool: "), "%t", tc.arg)
			fo := fmt.Appendf([]byte("Bool: "), "%t", tc.arg)
			if !bytes.Equal(o, fo) {
				t.Errorf("Expected %s, got %s", fo, o)
			}
		})
	}
}

func TestAppendfNil(t *testing.T) {
	o := xprint.Appendf([]byte("Nil: "), "%v", nil)
	fo := fmt.Appendf([]byte("Nil: "), "%v", nil)
	if !bytes.Equal(o, fo) {
		t.Errorf("Expected %s, got %s", fo, o)
	}
}

func TestAppendfBytes(t *testing.T) {
	data := []byte("byte slice")
	o := xprint.Appendf([]byte("Bytes: "), "%s", data)
	fo := fmt.Appendf([]byte("Bytes: "), "%s", data)
	if !bytes.Equal(o, fo) {
		t.Errorf("Expected %s, got %s", fo, o)
	}
}

func TestAppendfPointers(t *testing.T) {
	// Regular pointer
	str := "pointer test"
	strPtr := &str

	// Nil pointer
	var nilPtr *string

	t.Run("regular pointer", func(t *testing.T) {
		o := xprint.Appendf([]byte("Pointer: "), "%p", strPtr)
		fo := fmt.Appendf([]byte("Pointer: "), "%p", strPtr)
		if !bytes.Equal(o, fo) {
			t.Errorf("Expected %s, got %s", fo, o)
		}
	})

	t.Run("nil pointer", func(t *testing.T) {
		o := xprint.Appendf([]byte("Nil pointer: "), "%v", nilPtr)
		fo := fmt.Appendf([]byte("Nil pointer: "), "%v", nilPtr)
		if !bytes.Equal(o, fo) {
			t.Errorf("Expected %s, got %s", fo, o)
		}
	})
}

func TestAppendfSlices(t *testing.T) {
	intSlice := []int{1, 2, 3, 4, 5}
	stringSlice := []string{"a", "b", "c"}
	emptySlice := []float64{}

	testCases := []struct {
		name string
		arg  any
	}{
		{"int slice", intSlice},
		{"string slice", stringSlice},
		{"empty slice", emptySlice},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			o := xprint.Appendf([]byte("Slice: "), "%v", tc.arg)
			fo := fmt.Appendf([]byte("Slice: "), "%v", tc.arg)
			if !bytes.Equal(o, fo) {
				t.Errorf("Expected %s, got %s", fo, o)
			}
		})
	}
}

func TestAppendfMaps(t *testing.T) {
	stringMap := map[string]string{"key1": "value1", "key2": "value2"}
	intMap := map[string]int{"one": 1, "two": 2}
	emptyMap := map[int]bool{}

	t.Run("string map", func(t *testing.T) {
		o := xprint.Appendf([]byte("Map: "), "%v", stringMap)
		fo := fmt.Appendf([]byte("Map: "), "%v", stringMap)
		// For maps, we don't compare exact output as ordering may differ
		if len(o) != len(fo) {
			t.Errorf("Length mismatch for string map. Expected %d, got %d", len(fo), len(o))
		}
	})

	t.Run("int map", func(t *testing.T) {
		o := xprint.Appendf([]byte("Map: "), "%v", intMap)
		fo := fmt.Appendf([]byte("Map: "), "%v", intMap)
		if len(o) != len(fo) {
			t.Errorf("Length mismatch for int map. Expected %d, got %d", len(fo), len(o))
		}
	})

	t.Run("empty map", func(t *testing.T) {
		o := xprint.Appendf([]byte("Map: "), "%v", emptyMap)
		fo := fmt.Appendf([]byte("Map: "), "%v", emptyMap)
		if !bytes.Equal(o, fo) {
			t.Errorf("Expected %s, got %s", fo, o)
		}
	})
}

func TestAppendfStruct(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	person := Person{Name: "John", Age: 30}

	o := xprint.Appendf([]byte("Struct: "), "%+v", person)
	fo := fmt.Appendf([]byte("Struct: "), "%+v", person)
	if !bytes.Equal(o, fo) {
		t.Errorf("Expected %s, got %s", fo, o)
	}
}

func TestAppendfInterface(t *testing.T) {
	var i any = "interface value"

	o := xprint.Appendf([]byte("Interface: "), "%v", i)
	fo := fmt.Appendf([]byte("Interface: "), "%v", i)
	if !bytes.Equal(o, fo) {
		t.Errorf("Expected %s, got %s", fo, o)
	}
}

func TestAppendfError(t *testing.T) {
	err := fmt.Errorf("test error") //nolint:all //

	o := xprint.Appendf([]byte("Error: "), "%v", err)
	fo := fmt.Appendf([]byte("Error: "), "%v", err)
	if !bytes.Equal(o, fo) {
		t.Errorf("Expected %s, got %s", fo, o)
	}
}

func TestAppendfEmptyArgs(t *testing.T) {
	o := xprint.Appendf([]byte("No args: "), "")
	fo := fmt.Appendf([]byte("No args: "), "")
	if !bytes.Equal(o, fo) {
		t.Errorf("Expected %s, got %s", fo, o)
	}
}

func TestAppendfMultipleTypes(t *testing.T) {
	o := xprint.Appendf([]byte("Mixed: "), "%s %d %t %f %v", "string", 42, true, 3.14, []int{1, 2, 3})
	fo := fmt.Appendf([]byte("Mixed: "), "%s %d %t %f %v", "string", 42, true, 3.14, []int{1, 2, 3})
	if !bytes.Equal(o, fo) {
		t.Errorf("Expected %s, got %s", fo, o)
	}
}

func TestAppendfFormatSpecifiers(t *testing.T) {
	testCases := []struct {
		name   string
		format string
		arg    any
	}{
		{"decimal", "%d", 42},
		{"binary", "%b", 42},
		{"octal", "%o", 42},
		{"hex lowercase", "%x", 42},
		{"hex uppercase", "%X", 42},
		{"float", "%.2f", 3.14159},
		{"scientific notation", "%e", 3.14159},
		{"string width", "%10s", "test"},
		{"quoted string", "%q", "test"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			o := xprint.Appendf([]byte("Format: "), tc.format, tc.arg)
			fo := fmt.Appendf([]byte("Format: "), tc.format, tc.arg)
			if !bytes.Equal(o, fo) {
				t.Errorf("Expected %s, got %s", fo, o)
			}
		})
	}
}
