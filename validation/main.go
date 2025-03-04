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
	if os.Args[1] == "bench" {
		Runtest()
	} else if os.Args[1] == "quick" {
		Quicktest()
	} else if os.Args[1] == "newbench" {
		NewBenchmark()
	} else if os.Args[1] == "bigjson" {
		BigJSONTest()
	} else if os.Args[1] == "mixedtype" {
		MixedTypeTest()
	} else if os.Args[1] == "simple" {
		RunAll()
	}
}
