package main

import (
	"fmt"
	"github.com/jbenet/data"
	"path"
	"strings"
)

const IndexfileName = "Indexfile"

// Indexfile is the main index file that describes a dataset. It
// is merely a list of Refs (pointers to manifests), and a list of
// collaborators allowed to modify the package.
//
// Path is "datasets/<owner>/<name>/Indexfile"
type Indexfile struct {
	data.SerializedFile "-"

	Dataset string
	Refs    data.DatasetRefs

	// Other users allowed to modify the package.
	Collaborators map[string]string
}

func IndexfilePath(dataset string) string {
	return path.Join(data.DatasetDir, dataset, IndexfileName)
}

// Constructs a new Indexfile, based on its path: "<owner>/<name>"
func NewIndexfile(p string) (*Indexfile, error) {
	if !IndexfileNameRegexp.MatchString(p) {
		return nil, fmt.Errorf("invalid Indexfile path: %v", p)
	}

	f := &Indexfile{
		SerializedFile: data.SerializedFile{Path: p},
		Refs: data.DatasetRefs{
			Published: map[string]string{},
			Versions:  map[string]string{},
		},
	}
	f.SerializedFile.Format = f

	err := f.ReadFile()
	if err != nil {
		return f, err
	}

	return f, nil
}

func (f *Indexfile) Handle() *data.Handle {
	return data.NewHandle(f.Dataset)
}

func (f *Indexfile) Valid() bool {
	return f.Handle().Valid()
}

func (f *Indexfile) Name() string {
	return strings.Split(f.Path, "/")[2]
}

func (f *Indexfile) Owner() string {
	return strings.Split(f.Path, "/")[1]
}
