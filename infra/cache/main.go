package cache

import (
	"fmt"
	"github.com/dgraph-io/badger/v4"
	"github.com/gabrielmoura/WebCrawler/config"
	"github.com/gabrielmoura/WebCrawler/infra/db"
	"github.com/gabrielmoura/WebCrawler/infra/log"
	"time"
)

var cdb *badger.DB
var indexName = "visitedIndex"

func getBadgerMode() badger.Options {
	if config.Conf.Cache.Mode == "mem" {
		return badger.DefaultOptions("").WithInMemory(true)
	} else {
		return badger.DefaultOptions(config.Conf.Cache.DBDir)
	}
}

func InitCache() error {
	opts := getBadgerMode()
	opts.Logger = nil
	open, err := badger.Open(opts)
	if err != nil {
		return err
	}
	cdb = open

	que := NewBadgerQueue(cdb)
	queue = que
	return nil
}
func SyncCache() error {
	time.Sleep(1 * time.Second)
	log.Logger.Info("Syncing cache")
	visited, err := db.AllVisited()

	err = cdb.Update(func(txn *badger.Txn) error {
		for _, v := range visited {
			key := []byte(fmt.Sprintf("%s:%s", indexName, v))
			err := txn.Set(key, []byte{})
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error setting visitedIndex: %v", err)
	}
	return nil
}
func OptimizeCache() error {
	log.Logger.Info("Optimizing cache")
	err := cdb.RunValueLogGC(0.5)
	if err != nil {
		return fmt.Errorf("error optimizing cache: %v", err)
	}
	return nil
}
func IsVisited(url string) bool {
	key := []byte(fmt.Sprintf("%s:%s", indexName, url))
	err := cdb.View(func(txn *badger.Txn) error {
		_, err := txn.Get(key)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return false
	}
	return true
}
func SetVisited(url string) error {
	key := []byte(fmt.Sprintf("%s:%s", indexName, url))
	err := cdb.Update(func(txn *badger.Txn) error {
		err := txn.Set(key, []byte{})
		if err != nil {
			return err
		}
		return nil // return nil to commit the transaction
	})
	if err != nil {
		return fmt.Errorf("error setting visited: %v", err)
	}
	return nil
}

func AddToQueue(url string) error {

	err := queue.Enqueue(url)
	if err != nil {
		return fmt.Errorf("error adding to queue: %v", err)
	}
	return nil
}
func GetFromQueue() (string, error) {
	url, err := queue.Dequeue()
	if err != nil {
		return "", fmt.Errorf("error getting from queue: %v", err)
	}
	return url, nil
}
func GetFromQueueV2(getNumber int) ([]string, error) {
	var urls []string
	for i := 0; i < getNumber; i++ {
		url, err := queue.Dequeue()
		if err != nil {
			return nil, fmt.Errorf("error getting from queue: %v", err)
		}
		if url != "" {
			urls = append(urls, url)
		}
	}
	return urls, nil
}
