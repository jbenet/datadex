package datastore

import (
  "bytes"
  "fmt"
  "github.com/jbenet/data"
  ds "github.com/jbenet/datastore.go"
  "github.com/jbenet/datastore.go/leveldb"
  "github.com/jbenet/datastore.go/elastigo"
  esCore "github.com/mattbaird/elastigo/core"
  "io/ioutil"
)


type Config struct {
  DatabasePath string
  ElasticSearchAddress elastigo.Address
  Constructor func (ds.Key) Model
}

var defaultConfig = &Config{
  DatabasePath: ".datastore",
  ElasticSearchAddress: elastigo.Address{
    Host: "localhost",
    Port: 9300,
  },
}


type Datastore struct {
  db *leveldb.Datastore
  es *elastigo.Datastore
}

func NewDatastore(c *Config) (*Datastore, error) {
  if c == nil {
    c = defaultConfig
  }

  es, err := elastigo.NewDatastore(c.ElasticSearchAddress, "datadex")
  if err != nil {
    return nil, err
  }

  db, err := leveldb.NewDatastore(c.DatabasePath, nil)
  if err != nil {
    return nil, err
  }

  return &Datastore{db: db, es: es}, nil
}

func (d *Datastore) Put(key ds.Key, value Model) (err error) {
  buf, err := MarshalledBytes(value)
  if err != nil {
    return err
  }

  if err = d.db.Put(key, buf); err != nil {
    return err
  }

  ifds := value.IndexFields()
  if err = d.es.Put(key, ifds); err != nil {
    return err
  }

  return nil
}

func (d *Datastore) Get(key ds.Key) (value Model, err error) {
  // setup return type based on key type


  // get data from leveldb
  val, err := d.db.Get(key)
  if err != nil {
    return nil, err
  }

  buf, ok := val.([]byte)
  if !ok {
    return nil, ds.ErrInvalidType
  }

  // unmarshal it into the model
  r := bytes.NewReader(buf)
  err = data.Unmarshal(r, value)
  if err != nil {
    return nil, err
  }

  return value, nil
}

func (d *Datastore) Has(key ds.Key) (exists bool, err error) {
  return d.db.Has(key)
}

func (d *Datastore) Delete(key ds.Key) (err error) {
  if err := d.es.Delete(key); err != nil {
    return err
  }

  return d.db.Delete(key)
}

func (d *Datastore) Search(key ds.Key, query string) (*[]Model, error) {
  out, err := esCore.SearchRequest(true, d.es.Index(key), key.Name(), query, "", 0)
  if err != nil {
    return nil, err
  }

  models := []Model{}
  for _, h := range out.Hits.Hits {
    model, err := d.Get(key.Instance(h.Id))
    if err != nil {
      return nil, fmt.Errorf("%v (ElasticSearch inconsistent?)", err)
    }
    models = append(models, model)
  }
  return &models, nil
}

// Model

type Model interface {
  IndexFields() *map[string]interface{}
  Key() ds.Key
}


func MarshalledBytes(m Model) ([]byte, error) {
  r, err := data.Marshal(m)
  if err != nil {
    return nil, err
  }

  return ioutil.ReadAll(r)
}
