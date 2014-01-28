package datadex

import (
	"github.com/gorilla/mux"
	"github.com/jbenet/data"
	"log"
	"net/http"
)

var mainDataIndex *data.DataIndex

func init() {
	var err error
	mainDataIndex, err = data.NewMainDataIndex()
	if err != nil {
		log.Fatal(err)
	}
}

func dsBlobHandler(w http.ResponseWriter, r *http.Request) {
	ref := mux.Vars(r)["ref"]
	url := blobUrl(ref)
	pOut("302 %v -> %v\n", ref, url)
	http.Redirect(w, r, url, 302)
}

func blobUrl(ref string) string {
	return mainDataIndex.BlobStore.Url(data.BlobKey(ref))
}
