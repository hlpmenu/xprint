package main

import (
	"log"
	_ "net/http/pprof"
	"os"
)

func main() {
	args := os.Args[0:]

	if len(args) < 2 {
		log.Fatal("Need a arg")
	}

	// Use the arg variable from above, not directly os.Args
	arg := args[1]

	switch arg {
	case "bench":
		Runtest()
	case "quick":
		Quicktest()
	case "newbench":
		NewBenchmark()
	case "bigjson":
		BigJSONTest()
	case "mixedtype":
		MixedTypeTest()
	case "validate-printf":
		RunAll()
	case "validate-errorf":
		RunErrorfTests()
	case "nil-args-test":
		RunNilCheckTests()
	case "empty-args-bench":
		EmptyArgsBenchmark()
	case "types-bench":
		TypesBenchmark()
	case "test-floats":
		TestFloats()
	}

}
