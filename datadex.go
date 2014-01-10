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

	port := flag.Int("port", 8080, "Listen port")
	vers := flag.Bool("version", false, "Show version")
	flag.BoolVar(&DEBUG, "debug", false, "Debug mode")
	flag.Parse()

	dOut("debugging on\n")

	if *vers {
		pOut("datadex version: %s\n", Version)
		os.Exit(0)
	}

	r := NewDatadexRouter()
	http.Handle("/", r)

	addr := fmt.Sprintf("localhost:%d", *port)
	pOut("listening on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func pErr(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
}

func pOut(format string, a ...interface{}) {
	fmt.Fprintf(os.Stdout, format, a...)
}

func dErr(format string, a ...interface{}) {
	if DEBUG {
		pErr(format, a...)
	}
}

func dOut(format string, a ...interface{}) {
	if DEBUG {
		pOut(format, a...)
	}
}
