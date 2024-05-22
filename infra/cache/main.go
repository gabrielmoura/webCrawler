package cache

import (
	"encoding/json"
	"fmt"
	"github.com/dgraph-io/badger/v4"
	"github.com/gabrielmoura/WebCrawler/infra/db"
	"github.com/gabrielmoura/WebCrawler/infra/log"
	"time"
)

var cdb *badger.DB

func InitCache() error {
	opts := badger.DefaultOptions("").WithInMemory(true)
	opts.Logger = nil
	open, err := badger.Open(opts)
	if err != nil {
		return err
	}
	cdb = open

	err = SyncCache()
	if err != nil {
		return err
	}
	return nil
}
func SyncCache() error {
	time.Sleep(1 * time.Second)
	log.Logger.Info("Syncing cache")
	visited, err := db.AllVisited()
	if err != nil {
		return fmt.Errorf("error getting all visited: %v", err)
	}
	jsonBytes, err := json.Marshal(visited)
	if err != nil {
		fmt.Println("Erro ao converter para JSON:", err)
		return err
	}
	err = cdb.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte("visitedIndex"), jsonBytes)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error setting visitedIndex: %v", err)
	}
	return nil
}
func IsVisited(url string) bool {
	cdb.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("visitedIndex"))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			var visited []string
			err := json.Unmarshal(val, &visited)
			if err != nil {
				return err
			}
			for _, v := range visited {
				if v == url {
					return nil
				}
			}
			return fmt.Errorf("not found")
		})
	})
	return false
}
func SetVisited(url string) error {
	err := cdb.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("visitedIndex"))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			var visited []string
			err := json.Unmarshal(val, &visited)
			if err != nil {
				return err
			}
			visited = append(visited, url)
			jsonBytes, err := json.Marshal(visited)
			if err != nil {
				return err
			}
			return txn.Set([]byte("visitedIndex"), jsonBytes)
		})
	})
	if err != nil {
		return err
	}
	return nil
}
