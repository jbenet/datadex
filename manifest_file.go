package datadex

import (
	"fmt"
	"github.com/jbenet/data"
	"io/ioutil"
	"strings"
)

type ManifestFile struct {
	Path   string
	Ref    string
	buffer []byte
}

var noReadmeFile = &ManifestFile{
	Path: "readme",
	Ref:  "",
}

func (f *ManifestFile) Read() error {
	if f.Ref == "" {
		return fmt.Errorf("No reference to read.")
	}

	i, err := data.NewMainDataIndex()
	if err != nil {
		return err
	}

	rdr, err := i.BlobStore.Get(data.BlobKey(f.Ref))
	if err != nil {
		return err
	}

	// read the blob
	f.buffer, err = ioutil.ReadAll(rdr)
	if err != nil {
		return err
	}

	return nil
}

func (f *ManifestFile) Bytes() []byte {
	return f.buffer
}

func (f *ManifestFile) RenderedBytes() []byte {
	// render markdown?
	if strings.HasSuffix(f.Path, ".md") {
		return RenderMarkdownSafe(f.buffer)
	}
	return f.buffer
}

func (f *ManifestFile) CanBeRendered() bool {
	return strings.HasSuffix(f.Path, ".md")
}

func FileForManifestRef(ref string) (*ManifestFile, error) {
	mf, err := data.NewManifestWithRef(ref)
	if err != nil {
		return nil, fmt.Errorf("Error loading manifest: %s", err.Error())
	}

	attempts := []string{
		"readme.md",
		"readme.txt",
		"readme",
	}

	for _, attempt := range attempts {
		attempt = strings.ToLower(attempt)

		for path, hash := range mf.Files {
			path = strings.ToLower(path)
			if path == attempt {
				return &ManifestFile{
					Path: path,
					Ref:  hash,
				}, nil
			}
		}
	}

	return nil, ErrNotFound
}
