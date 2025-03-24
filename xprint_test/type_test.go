package xprint_test

import (
	"fmt"
	"strconv"
	"testing"

	xprint "gopkg.hlmpn.dev/pkg/xprint"
)

func TestInterface(t *testing.T) {
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

func TestInt32(t *testing.T) {
	s := int32(42)
	xf := xprint.Printf("%d", s)
	ff := fmt.Sprintf("%d", s)
	if xf != ff {
		t.Errorf("[Int32] ERROR: Mismatch between fmt.Sprintf and xprint.Printf")
	} else {
		t.Logf("[Int32] Success: Test passed")
	}
}

func TestInt64(t *testing.T) {
	s := int64(42)
	xf := xprint.Printf("%d", s)
	ff := fmt.Sprintf("%d", s)
	if xf != ff {
		t.Errorf("[Int64] ERROR: Mismatch between fmt.Sprintf and xprint.Printf")
	} else {
		t.Logf("[Int64] Success: Test passed")
	}
}

func TestInt(t *testing.T) {
	s := int(42)
	xf := xprint.Printf("%d", s)
	ff := fmt.Sprintf("%d", s)
	if xf != ff {
		t.Errorf("[Int] ERROR: Mismatch between fmt.Sprintf and xprint.Printf")
	} else {
		t.Logf("[Int] Success: Test passed")
	}
}

func TestMapStringInterface(t *testing.T) {
	m := map[string]interface{}{
		"hello": "world",
		"foo":   "bar",
	}
	xf := xprint.Printf("%v", m)
	ff := fmt.Sprintf("%v", m)
	t.Logf("--- Map:[Map[string]interface{}] ---")
	t.Logf("fmt output: %s", ff)
	t.Logf("xprint output: %s", xf)
	t.Logf("-------------------------------")
}

func TestMapStringString(t *testing.T) {
	m := map[string]string{
		"hello": "world",
		"foo":   "bar",
	}
	xf := xprint.Printf("%v", m)
	ff := fmt.Sprintf("%v", m)
	t.Logf("--- Map:[Map[string]string]] ---")
	t.Logf("fmt output: %s", ff)
	t.Logf("xprint output: %s", xf)
	t.Logf("-------------------------------")
}

func TestMapStringInt(t *testing.T) {
	m := map[string]int{
		"hello": 42,
		"foo":   43,
	}
	xf := xprint.Printf("%v", m)
	ff := fmt.Sprintf("%v", m)
	t.Logf("--- Map:[Map[string]int]] ---")
	t.Logf("fmt output: %s", ff)
	t.Logf("xprint output: %s", xf)
	t.Logf("-------------------------------")
}

func TestNilPointer(t *testing.T) {
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

func TestPointer(t *testing.T) {
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

func TestNilStruct(t *testing.T) {
	var s *struct{} = nil
	xf := xprint.Printf("%p", s)
	ff := fmt.Sprintf("%p", s)
	if xf != ff {
		t.Errorf("[NilStruct] ERROR: Mismatch between fmt.Sprintf and xprint.Printf")
	} else {
		t.Logf("[NilStruct] Success: Test passed")
	}
}

func TestStruct(t *testing.T) {
	s := struct{}{}
	xf := xprint.Printf("%p", &s)
	ff := fmt.Sprintf("%p", &s)
	if xf != ff {
		t.Errorf("[Struct] ERROR: Mismatch between fmt.Sprintf and xprint.Printf")
	} else {
		t.Logf("[Struct] Success: Test passed")
	}
}

func TestStructPointer(t *testing.T) {
	s := struct{}{}
	xf := xprint.Printf("%p", &s)
	ff := fmt.Sprintf("%p", &s)
	if xf != ff {
		t.Errorf("[StructPointer] ERROR: Mismatch between fmt.Sprintf and xprint.Printf")
	} else {
		t.Logf("[StructPointer] Success: Test passed")
	}
}

func TestStructPointer2(t *testing.T) {
	s := struct{}{}
	xf := xprint.Printf("%p", &s)
	ff := fmt.Sprintf("%p", &s)
	if xf != ff {
		t.Errorf("[StructPointer2] ERROR: Mismatch between fmt.Sprintf and xprint.Printf")
	} else {
		t.Logf("[StructPointer2] Success: Test passed")
	}
}

func TestBool(t *testing.T) {
	b := true
	xf := xprint.Printf("%t", b)
	ff := fmt.Sprintf("%t", b)
	t.Logf("strconv.FormatBool output: %s", strconv.FormatBool(b))
	t.Logf("xprint.Printf output: %s", xf)
	t.Logf("-------------------------------")
	t.Log("bool as s")
	xfs := xprint.Printf("%s", b)
	ffs := fmt.Sprintf("%s", b)
	if xfs != ffs {
		t.Logf("fmt output: %s", ffs)
		t.Logf("xprint output: %s", xfs)
		t.Errorf("[Bool] ERROR: Mismatch between fmt.Sprintf and xprint.Printf")
	}
	strconv.FormatBool(b)
	if xf != ff {
		t.Logf("fmt output: %s", ff)
		t.Logf("xprint output: %s", xf)
		t.Errorf("[Bool] ERROR: Mismatch between fmt.Sprintf and xprint.Printf")
	} else {
		t.Logf("[Bool] Success: Test passed")
	}
}

func TestAll(t *testing.T) {
	t.Run("Interface", TestInterface)
	t.Run("Int32", TestInt32)
	t.Run("Int64", TestInt64)
	t.Run("Int", TestInt)
	t.Run("MapStringInterface", TestMapStringInterface)
	t.Run("MapStringString", TestMapStringString)
	t.Run("MapStringInt", TestMapStringInt)
	t.Run("NilPointer", TestNilPointer)
	t.Run("NilStruct", TestNilStruct)
	t.Run("Pointer", TestPointer)
	t.Run("Struct", TestStruct)
	t.Run("StructPointer", TestStructPointer)
	t.Run("StructPointer2", TestStructPointer2)
	t.Run("Bool", TestBool)
}
