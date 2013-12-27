package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	// "github.com/gorilla/mux"
	"github.com/jbenet/data"
)

func publishPostHandler(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	ndf := &data.Datafile{}
	err := ndf.Read(r.Body)
	if err != nil {
		errs := fmt.Sprintf("File error: %s", err)
		http.Error(w, errs, http.StatusBadRequest)
		return
	}

	// point to the correct Datafile location
	ndf.Path = data.DatafilePath(ndf.Handle().Path())

	odf, err := data.NewDatafile(ndf.Path)
	if err == nil {
		// old version exists
		err = publishValidation(ndf, odf)
		if err != nil {
			errs := fmt.Sprintf("Validation error: %s", err)
			http.Error(w, errs, http.StatusBadRequest)
			return
		}

	} else {
		// new version. ensure dirs exist.
		err = os.MkdirAll(path.Dir(ndf.Path), 0777)
		if err != nil {
			log.Print("Datafile writing error: ", err)
			http.Error(w, "Datafile writing error.", http.StatusInternalServerError)
			return
		}
	}

	err = ndf.WriteFile()
	if err != nil {
		log.Print("Datafile writing error: ", err)
		http.Error(w, "Datafile writing error.", http.StatusInternalServerError)
		return
	}

	dsWriteDatafile(w, ndf)
}

func publishValidation(ndf *data.Datafile, odf *data.Datafile) error {

	// need to validate things here.
	// - ownership (username match)
	// - path match
	// - version increment (should be able to force old versions...?)

	return nil
}
