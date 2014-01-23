package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jbenet/data"
	"net/http"
	"path"
	"strings"
)

func dsHomeHandler(w http.ResponseWriter, r *http.Request) {
	dataset := requestDataset(r)
	fmt.Fprintf(w, "%s\n", dataset)
}

func dsIndexfileHandler(w http.ResponseWriter, r *http.Request) {
	ds := requestDataset(r)
	f, err := NewIndexfile(IndexfilePath(ds))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	httpWriteFile(w, &f.SerializedFile)
}

func dsDatafileHandler(w http.ResponseWriter, r *http.Request) {
	ds := requestDataset(r)
	f, err := data.NewDatafile(data.DatafilePath(ds))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	httpWriteFile(w, &f.SerializedFile)
}

func dsRefsHandler(w http.ResponseWriter, r *http.Request) {
	ds := requestDataset(r)
	f, err := NewIndexfile(IndexfilePath(ds))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	httpWriteMarshal(w, f.Refs)
}

func dsRefHandler(w http.ResponseWriter, r *http.Request) {
	ds := requestDataset(r)
	ref := mux.Vars(r)["ref"]

	if r.Method == "POST" {
		f, _ := NewIndexfile(IndexfilePath(ds))
		u, err := authenticatedUser(r)
		if err != nil {
			pErr("attempt to publish ref forbidden (%s)\n", err)
			http.Error(w, "Publishing forbidden.", http.StatusForbidden)
			return
		}

		if u.Disabled {
			pErr("attempt to publish ref forbidden (disabled %s to %s)\n", u.User, ds)
			http.Error(w, "Publishing forbidden.", http.StatusForbidden)
			return
		}

		if !f.UserCanModify(u.User()) {
			pErr("attempt to publish ref forbidden (%s to %s)\n", u.User, ds)
			http.Error(w, "Publishing forbidden.", http.StatusForbidden)
			return
		}

		err = publishRef(f, ref)
		if err != nil {
			pErr("%s\n", err)
			switch {
			case strings.Contains(err.Error(), "forbidden"):
				http.Error(w, "Publishing forbidden.", http.StatusForbidden)
			default:
				http.Error(w, "Error publishing ref.", http.StatusInternalServerError)
			}
			return
		}
	}

	f, err := NewIndexfile(IndexfilePath(ds))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	time, found := f.Refs.Published[ref]
	if !found {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "%s\n", time)
}

func dsArchivesHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "master\n")
}

func dsDownloadArchiveHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nam := vars["dataset"]
	ref := vars["ref"]
	ext := ".tar.gz" //vars["ext"]

	ds := requestDataset(r)
	path := path.Join(data.DatasetDir, ds, ext)

	disp := fmt.Sprintf("attachment; filename=%s-%s%s", nam, ref, ext)
	w.Header().Set("Content-Disposition", disp)
	w.Header().Set("Content-Type", "application/x-gzip")
	http.ServeFile(w, r, path)
}

func requestDataset(r *http.Request) string {
	vars := mux.Vars(r)
	dataset := path.Join(vars["user"], vars["dataset"])
	return dataset
}
