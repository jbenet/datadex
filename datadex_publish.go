package main

import (
	"fmt"
	"github.com/jbenet/data"
	"time"
)

func publishRef(f *Indexfile, ref string) error {
	// valid ref?
	if !data.IsHash(ref) {
		return fmt.Errorf("Invalid ref: %s", ref)
	}

	// already there?
	_, found := f.Refs.Published[ref]
	if found {
		return nil
	}

	df, err := DatafileForManifestRef(ref)
	if err != nil {
		return fmt.Errorf("Error loading datafile. %s", err.Error())
	}

	// no dataset? must be entirely new package.
	if len(f.Dataset) == 0 {
		f.Dataset = df.Handle().Path()
	}

	if f.Dataset != df.Handle().Path() {
		return fmt.Errorf("Attempt to publish ref (%.7s, %s) to"+
			" another dataset (%s) forbidden.\n", ref, df.Dataset, f.Dataset)
	}

	// ok, update it now :)
	f.Refs.Published[ref] = time.Now().UTC().String()
	f.Refs.Versions[df.Handle().Version] = ref
	err = f.WriteFile()
	if err != nil {
		return err
	}

	pOut("Published %s (%.7s)\n", df.Dataset, ref)
	return nil
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
