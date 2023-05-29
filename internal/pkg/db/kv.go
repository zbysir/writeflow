package db

import (
	"fmt"
	"github.com/docker/libkv"
	"github.com/docker/libkv/store"
	"github.com/docker/libkv/store/boltdb"
	"os"
	"path/filepath"
)

func init() {
	boltdb.Register()
}

type KvDb struct {
	dbDir string
}

func NewKvDb(dbDir string) *KvDb {
	return &KvDb{
		dbDir: dbDir,
	}
}

func (k *KvDb) Open(database, table string) (store.Store, error) {
	filename := filepath.Join(k.dbDir, fmt.Sprintf("%s.boltdb", database))
	dir := filepath.Dir(filename)
	if dir == "" {
		filename = "./" + filename
	}
	kv, err := libkv.NewStore(store.BOLTDB, []string{filename}, &store.Config{Bucket: table})
	if err != nil {
		return nil, fmt.Errorf("libkv.NewStore error: %w [%v]", err, filename)
	}
	return kv, nil
}

func (k *KvDb) Clean(database string) error {
	filename := filepath.Join(k.dbDir, fmt.Sprintf("%s.boltdb", database))
	err := os.Remove(filename)
	if err != nil {
		return err
	}
	return nil
}
