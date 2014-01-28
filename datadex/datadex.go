package main

import (
  "flag"
  "fmt"
  "github.com/jbenet/datadex"
  "log"
  "net/http"
  "os"
)

func main() {

  port := flag.Int("port", 8080, "Listen port")
  vers := flag.Bool("version", false, "Show version")
  flag.Parse()

  if *vers {
    fmt.Printf("datadex version: %s\n", datadex.Version)
    os.Exit(0)
  }

  r := datadex.NewDatadexRouter()
  http.Handle("/", r)

  addr := fmt.Sprintf("0.0.0.0:%d", *port)
  fmt.Printf("listening on %s\n", addr)
  log.Fatal(http.ListenAndServe(addr, nil))
}
