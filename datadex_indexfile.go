package main

import (
	"github.com/jbenet/data"
	"strings"
)

// This is the object that describes a dataset.
// It is merely a list of Refs (pointers to manifests), and a list of
// collaborators allowed to modify the package.
type Dataset struct {
	Path    string
	Name    string
	Owner   string
	Tagline string // replicated for convenience. use latest published.
	Refs    data.DatasetRefs

	// Other users allowed to modify the package.
	Collaborators map[string]string
}

func NewDataset(dataset string) *Dataset {
	parts := strings.Split(dataset, "/")
	return &Dataset{
		Path:  dataset,
		Name:  parts[0],
		Owner: parts[1],
	}
}

func (f *Dataset) Handle() *data.Handle {
	return data.NewHandle(f.Path)
}

func (f *Dataset) Valid() bool {
	return f.Handle().Valid()
}

func (f *Dataset) UserCanModify(user string) bool {
	if len(user) == 0 {
		return false
	}

	if user == f.Owner {
		return true
	}

	_, exists := f.Collaborators[user]
	return exists
}
