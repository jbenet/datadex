package main

import (
	"fmt"
	"github.com/jbenet/data"
	"log"
	"net/http"
	"os"
	"path"
	"time"
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

	dsWriteFile(w, &ndf.SerializedFile)
}

func publishValidation(ndf *data.Datafile, odf *data.Datafile) error {

	// need to validate things here.
	// - ownership (username match)
	// - path match
	// - version increment (should be able to force old versions...?)

	return nil
}

func publishRef(f *Indexfile, ref string) (bool, error) {
	// valid ref?
	if !data.IsHash(ref) {
		return false, fmt.Errorf("Invalid ref: %s", ref)
	}

	// already there?
	_, found := f.Refs.Published[ref]
	if found {
		return true, nil
	}

	df, err := DatafileForManifestRef(ref)
	if err != nil {
		return false, fmt.Errorf("Error loading datafile. %s", err.Error())
	}

	if f.Dataset != df.Handle().Path() {
		pErr("Attempt to publish ref (%.7s, %s) to another dataset (%s).\n",
			ref, df.Dataset, f.Dataset)
		return false, nil
	}

	// ok, update it now :)
	f.Refs.Published[ref] = time.Now().UTC().String()
	f.Refs.Versions[df.Handle().Version] = ref
	err = f.WriteFile()
	if err != nil {
		return false, err
	}

	pOut("Published %s (%.7s)\n", df.Dataset, ref)
	return true, nil
}

func DatafileForManifestRef(ref string) (*data.Datafile, error) {
	mf, err := data.NewManifestWithRef(ref)
	if err != nil {
		return nil, fmt.Errorf("Error loading manifest: %s", err.Error())
	}

	dref := mf.HashForPath(data.DatafileName)
	df, err := data.NewDatafileWithRef(dref)
	if err != nil {
		return nil, fmt.Errorf("Error loading datafile: %s", err.Error())
	}

	return df, nil
}
