package datadex

import (
	"fmt"
	"github.com/jbenet/data"
	"time"
)

func publishRef(f *Dataset, ref string) error {
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
	if len(f.Path) == 0 {
		f.Path = df.Handle().Path()
	}

	if f.Path != df.Handle().Path() {
		return fmt.Errorf("Attempt to publish ref (%.7s, %s) to"+
			" another dataset (%s) forbidden.\n", ref, df.Dataset, f.Path)
	}

	// add vesion object too
	ver := NewDatasetVersion(df.Handle())
	ver.Ref = ref
	ver.DatePublished = time.Now().UTC().String()
	ver.Put()

	// ok, update dataset now :)
	f.Tagline = df.Tagline
	f.Refs.Published[ref] = ver.DatePublished
	f.Refs.Versions[df.Handle().Version] = ref
	err = f.Put()

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
