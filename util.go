package main

import (
	"github.com/jbenet/data"
	"log"
	"net/http"
)

func httpWriteFile(w http.ResponseWriter, df *data.SerializedFile) {
	err := df.Write(w)
	if err != nil {
		log.Print("Error outputting SerializedFile: %s", err)
		http.Error(w, "Error in serialized file.", http.StatusInternalServerError)
		return
	}
}
