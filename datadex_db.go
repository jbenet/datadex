package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jbenet/data"
	"github.com/jbenet/tiedot/db"
	"math/rand"
	"os"
	"strings"
	"time"
)

var indexDB *IndexDB

type IndexDB struct {
	db    *db.DB
	users *db.Col
}

const UsersCollection = "Users"

var ErrNotFound = errors.New("Object not found.")
var ErrTooManyFound = errors.New("More than one object found.")

func init() {
	// It is very important to initialize random number generator seed!
	rand.Seed(time.Now().UTC().UnixNano())

	dbpath, err := data.ConfigGet("db.path")
	if err != nil {
		if len(err.Error()) == 0 {
			dbpath = ".indexdb"
		} else {
			pErr("Error loading db.path config: %v", err)
			os.Exit(-1)
		}
	}

	// Open database
	d, err := db.OpenDB(dbpath)
	if err != nil {
		pErr("Error opening db: %v", err)
		os.Exit(-1)
	}

	indexDB, err = NewIndexDB(d)
	if err != nil {
		pErr("Error creating db: %v", err)
		os.Exit(-1)
	}
}

func NewIndexDB(d *db.DB) (*IndexDB, error) {
	i := &IndexDB{db: d}
	var err error

	// collections to create { name : [index, ...] }
	collections := map[string][]string{
		UsersCollection: []string{
			"Username",
		},
	}

	for colname, indexes := range collections {
		if i.users, err = i.CreateCollection(colname, indexes); err != nil {
			return nil, err
		}
	}

	return i, nil
}

func (i *IndexDB) ColFindId(col *db.Col, q string) (uint64, error) {
	pOut("FIND: %s %s\n", col.BaseDir, q)
	var query interface{}
	json.Unmarshal([]byte(q), &query)

	results := make(map[uint64]struct{})
	if err := db.EvalQuery(query, col, &results); err != nil {
		return 0, err
	}

	switch len(results) {
	case 0:
		return 0, ErrNotFound
	case 1:
		for id := range results {
			return id, nil
		}
	}
	return 0, ErrTooManyFound
}

func (i *IndexDB) ColPutQuery(col *db.Col, q string, in interface{}) error {
	id, err := i.ColFindId(col, q)
	switch err {
	case nil:
		return i.ColPutId(col, id, in)
	case ErrNotFound:
		return i.ColPutId(col, 0, in)
	}
	return err
}

func (i *IndexDB) ColPutId(col *db.Col, id uint64, in interface{}) error {
	var wrap map[string]interface{}
	err := JsonMarshalUnmarshal(in, &wrap)
	if err != nil {
		return err
	}

	pOut("PUT %v %d\n", col.BaseDir, id)
	if id == 0 {
		_, err := col.Insert(wrap)
		return err
	}
	return col.Update(id, wrap)
}

func (i *IndexDB) ColGetId(col *db.Col, id uint64, out interface{}) error {
	var wrap map[string]interface{}
	if _, err := col.Read(id, &wrap); err != nil {
		return err
	}

	pOut("GET %v %d\n", col.BaseDir, id)
	return JsonMarshalUnmarshal(&wrap, out)
}

func (i *IndexDB) GetUser(name string) (*User, error) {
	if len(name) < 1 {
		return nil, fmt.Errorf("Cannot get user without username.")
	}

	u := &User{}
	q := fmt.Sprintf(`{"eq": "%s", "in": ["Username"]}`, name)
	id, err := i.ColFindId(i.users, q)
	if err != nil {
		return nil, err
	}
	err = i.ColGetId(i.users, id, &u)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (i *IndexDB) PutUser(user *User) error {
	if len(user.Username) < 1 {
		return fmt.Errorf("Cannot put user without username.")
	}

	q := fmt.Sprintf(`{"eq": "%s", "in": ["Username"]}`, user.Username)
	return i.ColPutQuery(i.users, q, user)
}

func (i *IndexDB) CreateCollection(name string, idx []string) (*db.Col, error) {
	// Create tables, if needed.
	col := i.db.Use(name)
	for col == nil {
		if err := i.db.Create(name, 1); err != nil {
			return nil, fmt.Errorf("Error creating table %s: %s", name, err)
		}
		col = i.db.Use(name)

	}

	// apply indexes
	for _, index := range idx {
		index_r := strings.Split(index, ",")
		pOut("Creating db index %s %v\n", name, index_r)
		if err := col.Index(index_r); err != nil {

			// silence already-indexed errors
			if strings.Contains(err.Error(), "already indexed in") {
				continue
			}

			return nil, fmt.Errorf("Error applying db index: %s, %v, %v",
				name, index, err)
		}
	}

	return col, nil
}

func JsonMarshalUnmarshal(in interface{}, out interface{}) error {
	data, err := json.Marshal(in)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, out)
}