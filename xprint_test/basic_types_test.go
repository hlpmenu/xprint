package xprint_test

import (
	"fmt"
	"strconv"
	"testing"

	xprint "gopkg.hlmpn.dev/pkg/xprint"
)

func TestInterfaceFormatting(t *testing.T) {
	s := "hello world!"
	i := interface{}(s)
	xf := xprint.Printf("%s", i)
	ff := fmt.Sprintf("%s", i)
	if xf != ff {
		t.Errorf("[Interface] ERROR: Mismatch between fmt.Sprintf and xprint.Printf")
	} else {
		t.Logf("[Interface] Success: Test passed")
	}
}

func TestInt32Formatting(t *testing.T) {
	s := int32(42)
	xf := xprint.Printf("%d", s)
	ff := fmt.Sprintf("%d", s)
	if xf != ff {
		t.Errorf("[Int32] ERROR: Mismatch between fmt.Sprintf and xprint.Printf")
	} else {
		t.Logf("[Int32] Success: Test passed")
	}
}

func TestInt64Formatting(t *testing.T) {
	s := int64(42)
	xf := xprint.Printf("%d", s)
	ff := fmt.Sprintf("%d", s)
	if xf != ff {
		t.Errorf("[Int64] ERROR: Mismatch between fmt.Sprintf and xprint.Printf")
	} else {
		t.Logf("[Int64] Success: Test passed")
	}
}

func TestIntFormatting(t *testing.T) {
	s := int(42)
	xf := xprint.Printf("%d", s)
	ff := fmt.Sprintf("%d", s)
	if xf != ff {
		t.Errorf("[Int] ERROR: Mismatch between fmt.Sprintf and xprint.Printf")
	} else {
		t.Logf("[Int] Success: Test passed")
	}
}

func TestMapStringInterfaceFormatting(t *testing.T) {
	m := map[string]interface{}{
		"hello": "world",
		"foo":   "bar",
	}
	xf := xprint.Printf("%v", m)
	ff := fmt.Sprintf("%v", m)
	t.Log("----------------------------------------")
	t.Logf("Map:[Map[string]interface{}]")
	t.Logf("fmt output: %s", ff)
	t.Logf("xprint output: %s", xf)
	t.Log("----------------------------------------")
}

func TestMapStringStringFormatting(t *testing.T) {
	m := map[string]string{
		"hello": "world",
		"foo":   "bar",
	}
	xf := xprint.Printf("%v", m)
	ff := fmt.Sprintf("%v", m)
	t.Log("----------------------------------------")
	t.Logf("Map:[Map[string]string]]")
	t.Logf("fmt output: %s", ff)
	t.Logf("xprint output: %s", xf)
	t.Log("----------------------------------------")
}

func TestMapStringIntFormatting(t *testing.T) {
	m := map[string]int{
		"hello": 42,
		"foo":   43,
	}
	xf := xprint.Printf("%v", m)
	ff := fmt.Sprintf("%v", m)
	t.Log("----------------------------------------")
	t.Logf("Map:[Map[string]int]]")
	t.Logf("fmt output: %s", ff)
	t.Logf("xprint output: %s", xf)
	t.Log("----------------------------------------")
}

func TestNilPointerFormatting(t *testing.T) {
	var s *string = nil
	xf := xprint.Printf("%p", s)
	ff := fmt.Sprintf("%p", s)
	if xf != ff {
		t.Logf("fmt output: %s", ff)
		t.Logf("xprint output: %s", xf)
		t.Errorf("[NilPointer] ERROR: Mismatch between fmt.Sprintf and xprint.Printf")
	} else {
		t.Logf("[NilPointer] Success: Test passed")
	}
}

func TestPointerFormatting(t *testing.T) {
	s := "Hello"
	xf := xprint.Printf("%p", &s)
	ff := fmt.Sprintf("%p", &s)
	if xf != ff {
		t.Logf("fmt output: %s", ff)
		t.Logf("xprint output: %s", xf)
		t.Errorf("[Pointer] ERROR: Mismatch between fmt.Sprintf and xprint.Printf")
	} else {
		t.Logf("[Pointer] Success: Test passed")
	}
}

func TestNilStructFormatting(t *testing.T) {
	var s *struct{} = nil
	xf := xprint.Printf("%p", s)
	ff := fmt.Sprintf("%p", s)
	if xf != ff {
		t.Errorf("[NilStruct] ERROR: Mismatch between fmt.Sprintf and xprint.Printf")
	} else {
		t.Logf("[NilStruct] Success: Test passed")
	}
}

func TestStructFormatting(t *testing.T) {
	s := struct{}{}
	xf := xprint.Printf("%p", &s)
	ff := fmt.Sprintf("%p", &s)
	if xf != ff {
		t.Errorf("[Struct] ERROR: Mismatch between fmt.Sprintf and xprint.Printf")
	} else {
		t.Logf("[Struct] Success: Test passed")
	}
}

func TestStructPointerFormatting(t *testing.T) {
	s := struct{}{}
	xf := xprint.Printf("%p", &s)
	ff := fmt.Sprintf("%p", &s)
	if xf != ff {
		t.Errorf("[StructPointer] ERROR: Mismatch between fmt.Sprintf and xprint.Printf")
	} else {
		t.Logf("[StructPointer] Success: Test passed")
	}
}

func TestStructPointer2Formatting(t *testing.T) {
	s := struct{}{}
	xf := xprint.Printf("%p", &s)
	ff := fmt.Sprintf("%p", &s)
	if xf != ff {
		t.Errorf("[StructPointer2] ERROR: Mismatch between fmt.Sprintf and xprint.Printf")
	} else {
		t.Logf("[StructPointer2] Success: Test passed")
	}
}

func TestBoolFormatting(t *testing.T) {
	b := true
	xf := xprint.Printf("%t", b)
	ff := fmt.Sprintf("%t", b)
	t.Logf("strconv.FormatBool output: %s", strconv.FormatBool(b))
	t.Logf("xprint.Printf output: %s", xf)
	t.Log("----------------------------------------")
	t.Log("bool as s")
	xfs := xprint.Printf("%s", b)
	ffs := fmt.Sprintf("%s", b)
	if xfs != ffs {
		t.Logf("fmt output: %s", ffs)
		t.Logf("xprint output: %s", xfs)
		t.Errorf("[Bool] ERROR: Mismatch between fmt.Sprintf and xprint.Printf for %%s")
	}

	if xf != ff {
		t.Logf("fmt output: %s", ff)
		t.Logf("xprint output: %s", xf)
		t.Errorf("[Bool] ERROR: Mismatch between fmt.Sprintf and xprint.Printf for %%t")
	} else {
		t.Logf("[Bool] Success: Test passed")
	}
}
func TestAllBasicTypes(t *testing.T) {
	// Run all tests as subtests
	t.Run("Interface", TestInterfaceFormatting)
	t.Run("Int32", TestInt32Formatting)
	t.Run("Int64", TestInt64Formatting)
	t.Run("Int", TestIntFormatting)
	t.Run("MapStringInterface", TestMapStringInterfaceFormatting)
	t.Run("MapStringString", TestMapStringStringFormatting)
	t.Run("MapStringInt", TestMapStringIntFormatting)
	t.Run("NilPointer", TestNilPointerFormatting)
	t.Run("NilStruct", TestNilStructFormatting)
	t.Run("Pointer", TestPointerFormatting)
	t.Run("Struct", TestStructFormatting)
	t.Run("StructPointer", TestStructPointerFormatting)
	t.Run("StructPointer2", TestStructPointer2Formatting)
	t.Run("Bool", TestBoolFormatting)
}
