package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

const Version = "0.0.0"

var DEBUG bool

func main() {

	port := flag.Int("port", 6000, "Listen port")
	vers := flag.Bool("version", false, "Show version")
	flag.BoolVar(&DEBUG, "debug", false, "Debug mode")
	flag.Parse()

	if DEBUG {
		DOut("debugging on\n")
	}

	if *vers {
		DOut("datadex version: %s\n", Version)
		os.Exit(0)
	}

	addr := fmt.Sprintf("localhost:%d", *port)
	DOut("listening on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func DErr(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
}

func DOut(format string, a ...interface{}) {
	fmt.Fprintf(os.Stdout, format, a...)
}
