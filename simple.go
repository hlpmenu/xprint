package main

import (
	"fmt"
	"strconv"

	"gopkg.hlmpn.dev/pkg/go-logger"
	xprint "gopkg.hlmpn.dev/pkg/xprint/t"
)

func TestInterface() {
	s := "hello world!"
	i := interface{}(s)
	xf := xprint.Printf("%s", i)
	ff := fmt.Sprintf("%s", i)
	if xf != ff {
		logger.LogErrorf("[Interface] ERROR: Mismatch between fmt.Sprintf and xprint.Printf")
	} else {
		logger.LogSuccessf("[Interface] Success: Test passed")
	}

}

func TestInt32() {
	s := int32(42)
	xf := xprint.Printf("%d", s)
	ff := fmt.Sprintf("%d", s)
	if xf != ff {
		logger.LogErrorf("[Int32] ERROR: Mismatch between fmt.Sprintf and xprint.Printf")
	} else {
		logger.LogSuccessf("[Int32] Success: Test passed")
	}
}

func TestInt64() {
	s := int64(42)
	xf := xprint.Printf("%d", s)
	ff := fmt.Sprintf("%d", s)
	if xf != ff {
		logger.LogErrorf("[Int64] ERROR: Mismatch between fmt.Sprintf and xprint.Printf")
	} else {
		logger.LogSuccessf("[Int64] Success: Test passed")
	}
}

func TestInt() {
	s := int(42)
	xf := xprint.Printf("%d", s)
	ff := fmt.Sprintf("%d", s)
	if xf != ff {
		logger.LogErrorf("[Int] ERROR: Mismatch between fmt.Sprintf and xprint.Printf")
	} else {
		logger.LogSuccessf("[Int] Success: Test passed")
	}
}

func TestMapStringInterface() {
	m := map[string]interface{}{
		"hello": "world",
		"foo":   "bar",
	}
	xf := xprint.Printf("%v", m)
	ff := fmt.Sprintf("%v", m)
	LogLine()
	logger.LogOrangef("Map:[Map[string]interface{}]")
	logger.LogPurplef("fmt output: %s", ff)
	logger.LogOrangef("xprint output: %s", xf)
	LogLine()

}

func TestMapStringString() {
	m := map[string]string{
		"hello": "world",
		"foo":   "bar",
	}
	xf := xprint.Printf("%v", m)
	ff := fmt.Sprintf("%v", m)
	LogLine()
	logger.LogOrangef("Map:[Map[string]string]]")
	logger.LogPurplef("fmt output: %s", ff)
	logger.LogOrangef("xprint output: %s", xf)
	LogLine()

}

func TestMapStringInt() {
	m := map[string]int{
		"hello": 42,
		"foo":   43,
	}
	xf := xprint.Printf("%v", m)
	ff := fmt.Sprintf("%v", m)
	LogLine()
	logger.LogOrangef("Map:[Map[string]int]]")
	logger.LogPurplef("fmt output: %s", ff)
	logger.LogOrangef("xprint output: %s", xf)
	LogLine()
}

func TestNilPointer() {
	var s *string = nil
	xf := xprint.Printf("%p", s)
	ff := fmt.Sprintf("%p", s)
	if xf != ff {
		logger.LogPurplef("fmt output: %s", ff)
		logger.LogOrangef("xprint output: %s", xf)
		logger.LogErrorf("[NilPointer] ERROR: Mismatch between fmt.Sprintf and xprint.Printf")
	} else {
		logger.LogSuccessf("[NilPointer] Success: Test passed")
	}

}

func TestPointer() {
	s := "Hello"
	xf := xprint.Printf("%p", &s)
	ff := fmt.Sprintf("%p", &s)
	if xf != ff {
		logger.LogPurplef("fmt output: %s", ff)
		logger.LogOrangef("xprint output: %s", xf)
		logger.LogErrorf("[Pointer] ERROR: Mismatch between fmt.Sprintf and xprint.Printf")
	} else {
		logger.LogSuccessf("[Pointer] Success: Test passed")
	}
}
func TestNilStruct() {
	var s *struct{} = nil
	xf := xprint.Printf("%p", s)
	ff := fmt.Sprintf("%p", s)
	if xf != ff {
		logger.LogErrorf("[NilStruct] ERROR: Mismatch between fmt.Sprintf and xprint.Printf")
	} else {
		logger.LogSuccessf("[NilStruct] Success: Test passed")
	}

}

func TestStruct() {
	s := struct{}{}
	xf := xprint.Printf("%p", &s)
	ff := fmt.Sprintf("%p", &s)
	if xf != ff {
		logger.LogErrorf("[Struct] ERROR: Mismatch between fmt.Sprintf and xprint.Printf")
	} else {
		logger.LogSuccessf("[Struct] Success: Test passed")
	}
}

func TestStructPointer() {
	s := struct{}{}
	xf := xprint.Printf("%p", &s)
	ff := fmt.Sprintf("%p", &s)
	if xf != ff {
		logger.LogErrorf("[StructPointer] ERROR: Mismatch between fmt.Sprintf and xprint.Printf")
	} else {
		logger.LogSuccessf("[StructPointer] Success: Test passed")
	}
}
func TestStructPointer2() {
	s := struct{}{}
	xf := xprint.Printf("%p", &s)
	ff := fmt.Sprintf("%p", &s)
	if xf != ff {
		logger.LogErrorf("[StructPointer2] ERROR: Mismatch between fmt.Sprintf and xprint.Printf")
	} else {
		logger.LogSuccessf("[StructPointer2] Success: Test passed")
	}
}

func TestBool() {
	b := true
	xf := xprint.Printf("%t", b)
	ff := fmt.Sprintf("%t", b)
	logger.LogOrangef("strconv.FormatBool output: %s", strconv.FormatBool(b))
	logger.LogOrangef("xprint.Printf output: %s", xf)
	LogLine()
	logger.Log("bool as s")
	xfs := xprint.Printf("%s", b)
	ffs := fmt.Sprintf("%s", b)
	if xfs != ffs {
		logger.LogPurplef("fmt output: %s", ffs)
		logger.LogOrangef("xprint output: %s", xfs)
		logger.LogErrorf("[Bool] ERROR: Mismatch between fmt.Sprintf and xprint.Printf")
	}
	strconv.FormatBool(b)
	if xf != ff {
		logger.LogPurplef("fmt output: %s", ff)
		logger.LogOrangef("xprint output: %s", xf)
		logger.LogErrorf("[Bool] ERROR: Mismatch between fmt.Sprintf and xprint.Printf")
	} else {
		logger.LogSuccessf("[Bool] Success: Test passed")
	}
}

func RunAll() {
	TestInterface()
	TestInt32()
	TestInt64()
	TestInt()
	TestMapStringInterface()
	TestMapStringString()
	TestMapStringInt()
	TestNilPointer()
	TestNilStruct()
	TestPointer()
	TestStruct()
	TestStructPointer()
	TestStructPointer2()
	TestBool()
}
