package meetup

import (
  "github.com/hashicorp/go-memdb"
  "errors"
  "sort"
)

const (
  TABLE_NAME        = "talk"
  ID_INDEX_NAME     = "id"
)

type Cache struct {
  db *memdb.MemDB
}

func (cache *Cache) GetTalk(key string) (*Talk, error) {
  transaction := cache.db.Txn(false)

  result, err := transaction.First(TABLE_NAME, ID_INDEX_NAME, key)

  if err != nil {
    transaction.Abort()
    return nil, err
  }

  talk, ok := result.(*Talk)

  if !ok {
    return nil, errors.New("Could not cast DB object to Talk")
  }

  transaction.Commit()
  return talk, nil
}

func (cache *Cache) GetTalks() (*Talks, error) {
  transaction := cache.db.Txn(false)

  results, err := transaction.Get(TABLE_NAME, ID_INDEX_NAME)

  if err != nil {
    transaction.Abort()
    return nil, err
  }

  talks := Talks{}

  for {
    item := results.Next()
    if item == nil {
      break
    }

    talk, ok := item.(*Talk)

    if !ok {
      return nil, errors.New("Could not cast DB object to Talk")
    }

    talks.Talks = append(talks.Talks, talk)
  }
  transaction.Commit()
  sort.Sort(ByDate{talks})
  return &talks, nil
}

func (cache *Cache) Fill(talks *Talks) error {
  transaction := cache.db.Txn(true)

  _, err := transaction.DeleteAll(TABLE_NAME, ID_INDEX_NAME)

  if err != nil {
    transaction.Abort()
    return err
  }

  for _, talk := range talks.Talks {
    transaction.Insert(TABLE_NAME, talk)
  }

  transaction.Commit()
  return nil
}

func (cache *Cache) Update(talk *Talk) error {
  transaction := cache.db.Txn(true)

  _, err := transaction.DeleteAll(TABLE_NAME, ID_INDEX_NAME, talk.Key)

  if err != nil {
    transaction.Abort()
    return err
  }

  err = transaction.Insert(TABLE_NAME, talk)
  if err != nil {
    transaction.Abort()
    return err
  }

  transaction.Commit()
  return nil
}

func Init() *Cache {
  schema := &memdb.DBSchema{
    Tables: map[string]*memdb.TableSchema{
      TABLE_NAME: {
        Name: TABLE_NAME,
        Indexes: map[string]*memdb.IndexSchema{
          ID_INDEX_NAME: {
            Name:    ID_INDEX_NAME,
            Unique:  true,
            Indexer: &memdb.StringFieldIndex{Field: "Key"},
          },
        },
      },
    },
  }

  db, err := memdb.NewMemDB(schema)

  if err != nil {
    panic(err)
  }

  return &Cache{db}
}
