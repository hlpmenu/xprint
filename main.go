package main

import (
	"log"
	"net/http"
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
	}
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	select {}

}
