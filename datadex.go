package main

import (
  "flag"
  "fmt"
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

  DOut("listening on localhost:%d\n", *port)
}

func DErr(format string, a ...interface{}) {
  fmt.Fprintf(os.Stderr, format, a...)
}

func DOut(format string, a ...interface{}) {
  fmt.Fprintf(os.Stdout, format, a...)
}
