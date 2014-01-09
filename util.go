package main

import (
	"crypto/rand"
	"github.com/jbenet/data"
	"io"
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

func httpWriteMarshal(w http.ResponseWriter, out interface{}) {
	rdr, err := data.Marshal(out)
	if err != nil {
		http.Error(w, "error serializing", http.StatusInternalServerError)
		return
	}

	_, err = io.Copy(w, rdr)
	if err != nil {
		http.Error(w, "Error writing response.", http.StatusInternalServerError)
	}
}

func randString(length int) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return string(b[:]), nil
}
