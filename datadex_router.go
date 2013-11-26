package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func NewDatadexRouter() *mux.Router {
	r := mux.NewRouter()

	// dataset
	r.HandleFunc("/{author}/{dataset}", DSHomeHandler)
	d := r.PathPrefix("/{author}/{dataset}").Subrouter()
	d.StrictSlash(true)
	d.HandleFunc("/", DSHomeHandler)
	d.HandleFunc("/Datafile", DSDatafileHandler)
	// d.HandleFunc("/refs", DSRefsHandler)
	// d.HandleFunc("/tree/{ref}/", DSTreeHandler)
	// d.HandleFunc("/blob/{ref}/", DSBlobHandler)
	d.HandleFunc("/archive/", DSArchivesHandler)
	d.HandleFunc("/archives/", DSArchivesHandler)
	d.HandleFunc("/archive/{ref}.tar.gz", DSDownloadArchiveHandler)
	// d.HandleFunc("/archive/{ref}.zip", DSArchiveHandler)

	r.HandleFunc("/", HomeHandler)
	return r
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "datadex/1.0\n")
}
