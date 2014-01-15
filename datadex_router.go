package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jbenet/data"
	"net/http"
)

func NewDatadexRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler).Methods("GET")
	r.HandleFunc("/version", versionHandler).Methods("GET")

	// api
	a := r.PathPrefix(data.ApiUrlSuffix).Subrouter()

	// user
	a.HandleFunc("/{author}", userHandler)
	u := a.PathPrefix("/{author}").Subrouter()
	u.StrictSlash(true)

	u.HandleFunc("/", userHandler).Methods("GET")
	u.HandleFunc("/user/add", userAddHandler).Methods("POST")
	u.HandleFunc("/user/info", userInfoHandler).Methods("GET", "POST")
	u.HandleFunc("/user/pass", userPassHandler).Methods("POST")
	u.HandleFunc("/user/auth", userAuthHandler).Methods("POST")
	u.HandleFunc("/user/awscred", userAwsCredHandler).Methods("GET")

	// user/dataset
	u.HandleFunc("/{dataset}", dsHomeHandler)
	d := u.PathPrefix("/{dataset}").Subrouter()
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
	return r
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "datadex -- the DATAset inDEX\n")
}

func versionHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "datadex/%s\n", Version)
}
