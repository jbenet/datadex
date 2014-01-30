package datadex

import (
	"fmt"
	"github.com/jbenet/data"
	"github.com/jbenet/datadex/datastore"
	ds "github.com/jbenet/datastore.go"
	"github.com/jbenet/datastore.go/elastigo"
	"os"
	"strings"
)

var indexDB *IndexDB

type IndexDB struct {
	ds *datastore.Datastore
}

var kUser = ds.NewKey("/User")
var kDataset = ds.NewKey("/Dataset")
var kDatasetVersion = ds.NewKey("/DatasetVersion")

var ErrNotFound = ds.ErrNotFound

func init() {

	dbpath, err := data.ConfigGet("db.path")
	if err != nil {
		if len(err.Error()) == 0 {
			dbpath = ".indexdb"
		} else {
			pErr("Error loading db.path config: %v", err)
			os.Exit(-1)
		}
	}

	// Setup Datastore
	d, err := datastore.NewDatastore(&datastore.Config{
		DatabasePath: dbpath,
		ElasticSearchAddress: elastigo.Address{
			Host: "localhost",
			Port: 9200,
		},
		Constructor: NewInstanceForKey,
	})

	if err != nil {
		pErr("Error on datastore init: %v", err)
		os.Exit(-1)
	}

	indexDB, err = NewIndexDB(d)
	if err != nil {
		pErr("Error on db init: %v", err)
		os.Exit(-1)
	}
}

func NewIndexDB(d *datastore.Datastore) (*IndexDB, error) {
	i := &IndexDB{ds: d}
	return i, nil
}

// Datadex specific stuff:

// User

func (i *IndexDB) GetUsers() ([]*User, error) {
	hits, err := i.ds.Search(kUser, "")
	if err != nil {
		return nil, err
	}

	users := make([]*User, len(*hits))
	for k, v := range *hits {
		users[k] = v.(*User)
	}
	return users, nil
}

func (i *IndexDB) GetUser(name string) (*User, error) {
	if len(name) < 1 {
		return nil, fmt.Errorf("Cannot get user without username.")
	}

	r, err := i.ds.Get(UserKey(name))
	out, _ := r.(*User)
	return out, err
}

func (i *IndexDB) PutUser(user *User) error {
	if len(user.Username) < 1 {
		return fmt.Errorf("Cannot put user without username.")
	}

	return i.ds.Put(user.Key(), user)
}

// Dataset

func (i *IndexDB) GetUserDatasets(username string) ([]*Dataset, error) {
	if len(username) < 1 {
		return nil, fmt.Errorf("No username provided")
	}

	query := `{"query" : {"term" : { "owner" : "` + username + `" }}}`
	hits, err := i.ds.Search(kDataset, query)
	if err != nil {
		return nil, err
	}
	rets := make([]*Dataset, len(*hits))
	for k, v := range *hits {
		rets[k] = v.(*Dataset)
	}
	return rets, nil
}

func (i *IndexDB) GetDatasets() ([]*Dataset, error) {
	hits, err := i.ds.Search(kDataset, "")
	if err != nil {
		return nil, err
	}

	rets := make([]*Dataset, len(*hits))
	for k, v := range *hits {
		rets[k] = v.(*Dataset)
	}
	return rets, nil
}

func (i *IndexDB) GetDataset(path string) (*Dataset, error) {
	parts := strings.Split(path, "/")
	if len(parts) != 2 && len(parts[0]) > 0 && len(parts[1]) > 0 {
		return nil, fmt.Errorf("Invalid dataset path: '%s'.", path)
	}

	// /User:<username>/Dataset:<dataset>
	ret, err := i.ds.Get(DatasetKey(parts[0], parts[1]))
	out, _ := ret.(*Dataset)
	return out, err
}

func (i *IndexDB) PutDataset(d *Dataset) error {
	if !d.Valid() {
		return fmt.Errorf("Invalid dataset: %v", d)
	}

	return i.ds.Put(d.Key(), d)
}

// DatasetVersion

func (i *IndexDB) GetDatasetVersions(path string) ([]*DatasetVersion, error) {
	q := ""
	if len(path) > 0 {
		q = "path:" + path
	}

	hits, err := i.ds.Search(kDatasetVersion, q)
	if err != nil {
		return nil, err
	}

	rets := make([]*DatasetVersion, len(*hits))
	for k, v := range *hits {
		rets[k] = v.(*DatasetVersion)
	}
	return rets, nil
}

func (i *IndexDB) GetDatasetVersion(h *data.Handle) (*DatasetVersion, error) {
	if h == nil || !h.Valid() {
		return nil, fmt.Errorf("Invalid dataset handle: %v.", h)
	}

	ret, err := i.ds.Get(HandleKey(h))
	out, _ := ret.(*DatasetVersion)
	return out, err
}

func (i *IndexDB) PutDatasetVersion(dv *DatasetVersion) error {
	if !dv.Valid() {
		return fmt.Errorf("Invalid dataset version: %v", dv)
	}

	return i.ds.Put(dv.Key(), dv)
}

// Model key specifics

func NewInstanceForKey(key ds.Key) datastore.Model {
	switch key.Type() {
	case "User":
		return &User{}
	case "Dataset":
		return &Dataset{}
	}
	return nil
}

func UserKey(username string) ds.Key {
	return kUser.Instance(username)
}

func DatasetKey(owner, name string) ds.Key {
	return kUser.Instance(owner).Child(kDataset.Name()).Instance(name)
}

func DatasetVersionKey(owner, name, version string) ds.Key {
	return kUser.Instance(owner).
		Child(kDataset.Name()).Instance(name).
		Child(kDatasetVersion.Name()).Instance(version)
}

func HandleKey(f *data.Handle) ds.Key {
	return DatasetVersionKey(f.Author, f.Name, f.Version)
}

func (f *User) Key() ds.Key {
	return UserKey(f.Username)
}

func (f *Dataset) Key() ds.Key {
	return DatasetKey(f.Owner, f.Name)
}

func (f *DatasetVersion) Key() ds.Key {
	return HandleKey(f.Handle())
}

// IndexFields implementations

func (f *User) IndexFields() *map[string]interface{} {
	m := &map[string]interface{}{
		"key":             f.Key().String(),
		"link":            "/" + f.Username,
		"username":        f.Username,
		"email":           f.Profile.Email,
		"name":            f.Profile.Name,
		"date_registered": f.DateRegistered,
	}
	return m
}

func (f *Dataset) IndexFields() *map[string]interface{} {
	m := &map[string]interface{}{
		"key":          f.Key().String(),
		"link":	        "/" + f.Path,
		"path":         f.Path,
		"owner":        f.Owner,
		"name":         f.Name,
		"tagline":      f.Tagline,
		"last_updated": f.Refs.LastUpdated(),
		"num_versions": len(f.Refs.Versions),
	}
	return m
}

func (f *DatasetVersion) IndexFields() *map[string]interface{} {
	m := &map[string]interface{}{
		"key":            f.Key().String(),
		"link":	          "/" + f.Dataset,
		"dataset":        f.Dataset,
		"path":           f.Path,
		"version":        f.Version,
		"ref":            f.Ref,
		"date_published": f.DatePublished,
		"num_views":      f.NumViews,
		"num_downloads":  f.NumDownloads,
	}
	return m
}
