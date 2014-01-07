package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jbenet/data"
	"log"
	"net/http"
	"path"
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

	dsWriteFile(w, &f.SerializedFile)
}

func dsDatafileHandler(w http.ResponseWriter, r *http.Request) {
	ds := requestDataset(r)
	f, err := data.NewDatafile(data.DatafilePath(ds))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	dsWriteFile(w, &f.SerializedFile)
}

func dsWriteFile(w http.ResponseWriter, df *data.SerializedFile) {
	err := df.Write(w)
	if err != nil {
		log.Print("Error outputting SerializedFile: %s", err)
		http.Error(w, "Error in serialized file.", http.StatusInternalServerError)
		return
	}
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
