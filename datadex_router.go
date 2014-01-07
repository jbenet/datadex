package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func NewDatadexRouter() *mux.Router {
	r := mux.NewRouter()

	// dataset
	r.HandleFunc("/{author}/{dataset}", dsHomeHandler)
	d := r.PathPrefix("/{author}/{dataset}").Subrouter()
	d.StrictSlash(true)

	dget := d.Methods("GET").Subrouter()
	// dpost := d.Methods("POST").Subrouter()

	dget.HandleFunc("/", dsHomeHandler)
	dget.HandleFunc("/Indexfile", dsIndexfileHandler)
	dget.HandleFunc("/Datafile", dsDatafileHandler)
	dget.HandleFunc("/refs", dsRefsHandler)
	d.HandleFunc("/refs/{ref}", dsRefHandler).Methods("GET", "POST")
	// dget.HandleFunc("/tree/{ref}/", dsTreeHandler)
	dget.HandleFunc("/blob/{ref}/", dsBlobHandler)
	dget.HandleFunc("/archive/", dsArchivesHandler)
	dget.HandleFunc("/archives/", dsArchivesHandler)
	dget.HandleFunc("/archive/{ref}.tar.gz", dsDownloadArchiveHandler)
	// dget.HandleFunc("/archive/{ref}.zip", dsArchiveHandler)

	// publish
	rget := d.Methods("GET").Subrouter()
	rpost := d.Methods("POST").Subrouter()
	rpost.HandleFunc("/publish", publishPostHandler).Methods("POST")

	rget.HandleFunc("/version", versionHandler)
	rget.HandleFunc("/", homeHandler)
	return r
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "datadex -- the DATAset inDEX\n")
}

func versionHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "datadex/%s\n", Version)
}
