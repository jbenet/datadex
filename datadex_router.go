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
	d.HandleFunc("/", dsHomeHandler)
	d.HandleFunc("/Indexfile", dsIndexfileHandler)
	d.HandleFunc("/Datafile", dsDatafileHandler)
	// d.HandleFunc("/refs", dsRefsHandler)
	// d.HandleFunc("/tree/{ref}/", dsTreeHandler)
	d.HandleFunc("/blob/{ref}/", dsBlobHandler)
	d.HandleFunc("/archive/", dsArchivesHandler)
	d.HandleFunc("/archives/", dsArchivesHandler)
	d.HandleFunc("/archive/{ref}.tar.gz", dsDownloadArchiveHandler)
	// d.HandleFunc("/archive/{ref}.zip", dsArchiveHandler)

	// publish
	r.HandleFunc("/publish", publishPostHandler).Methods("POST")

	r.HandleFunc("/version", versionHandler)
	r.HandleFunc("/", homeHandler)
	return r
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "datadex -- the DATAset inDEX\n")
}

func versionHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "datadex/%s\n", Version)
}
