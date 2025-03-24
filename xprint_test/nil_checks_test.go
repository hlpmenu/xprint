package xprint_test

import (
	"fmt"
	"testing"

	"gopkg.hlmpn.dev/pkg/xprint"
)

func TestNilChecks(t *testing.T) {
	t.Log("Running nil check tests...")

	s := xprint.Printf("%v", nil)
	ff := fmt.Sprintf("%v", nil)
	if s != ff {
		t.Errorf("[NilCheck] ERROR: Mismatch between fmt.Sprintf and xprint.Printf")
	} else {
		t.Logf("[NilCheck] Success: Test passed")
	}

	a := xprint.Printf("Hello world")
	b := fmt.Sprintf("Hello world")
	if a != b {
		t.Errorf("[NilCheck] ERROR: Mismatch between fmt.Sprintf and xprint.Printf")
	} else {
		t.Logf("[NilCheck] Success: Test passed")
	}
}
