package datastore

import (
  "bytes"
  "encoding/json"
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
  ElasticSearchUrl string
  Constructor func (ds.Key) Model
}

var defaultConfig = &Config{
  DatabasePath: ".datastore",
  ElasticSearchUrl: "http://localhost:9200/datadex",
}


type Datastore struct {
  config *Config
  dbr *leveldb.Datastore
  esr *elastigo.Datastore
  db *ds.LogDatastore
  es *ds.LogDatastore
}

func NewDatastore(c *Config) (*Datastore, error) {
  if c == nil {
    c = defaultConfig
  }

  es, err := elastigo.NewDatastore(c.ElasticSearchUrl)
  if err != nil {
    return nil, err
  }

  db, err := leveldb.NewDatastore(c.DatabasePath, nil)
  if err != nil {
    return nil, err
  }

  return &Datastore{
    config: c,
    dbr: db,
    esr: es,
    db: ds.NewLogDatastore(db, "leveldb"),
    es: ds.NewLogDatastore(es, "elastigo"),
  }, nil
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
  value = d.config.Constructor(key)

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
  out, err := esCore.SearchRequest(true, d.esr.Index(key), key.Name(), query, "", 0)
  if err != nil {
    return nil, err
  }

  models := []Model{}
  for _, h := range out.Hits.Hits {

    // unmarshal response
    var data map[string]interface{}
    err := json.Unmarshal(h.Source, &data)
    if err != nil {
      return nil, fmt.Errorf("Unmarshal error: %v", err)
    }

    // get key
    k, ok := data["key"].(string)
    if !ok {
      return nil, fmt.Errorf("Unmarshal error. Key not a string: %v", data["key"])
    }

    // retrieve object
    model, err := d.Get(ds.NewKey(k))
    if err != nil {
      return nil, fmt.Errorf("%v (ElasticSearch inconsistent?)", err)
    }

    // ok! we have it.
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
