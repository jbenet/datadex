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

	port := flag.Int("port", 9000, "Listen port")
	vers := flag.Bool("version", false, "Show version")
	flag.BoolVar(&DEBUG, "debug", false, "Debug mode")
	flag.Parse()

	DOut("debugging on\n")

	if *vers {
		Out("datadex version: %s\n", Version)
		os.Exit(0)
	}

	r := NewDatadexRouter()
	http.Handle("/", r)

	addr := fmt.Sprintf("localhost:%d", *port)
	Out("listening on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func Err(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
}

func Out(format string, a ...interface{}) {
	fmt.Fprintf(os.Stdout, format, a...)
}

func DErr(format string, a ...interface{}) {
	if DEBUG {
		Err(format, a...)
	}
}

func DOut(format string, a ...interface{}) {
	if DEBUG {
		Out(format, a...)
	}
}
