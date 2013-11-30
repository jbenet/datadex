package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jbenet/data"
	"net/http"
	"path"
)

func dsHomeHandler(w http.ResponseWriter, r *http.Request) {
	dataset := requestDataset(r)
	fmt.Fprintf(w, "%s\n", dataset)
}

func dsDatafileHandler(w http.ResponseWriter, r *http.Request) {
	ds := requestDataset(r)
	df, err := data.NewDatafile(data.DatafilePath(ds))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	buf, err := df.Marshal()
	if err != nil {
		http.Error(w, "Error fetching data.", 500)
		return
	}

	fmt.Fprintf(w, "%s\n", buf)
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
	dataset := path.Join(vars["author"], vars["dataset"])
	return dataset
}
