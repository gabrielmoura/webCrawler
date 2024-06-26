package cache

import (
	"fmt"
	"github.com/dgraph-io/badger/v4"
	"github.com/gabrielmoura/WebCrawler/config"
	"strconv"
)

var queue *BadgerQueue

// Queue represents an interface for interacting with a queue data structure.
type Queue interface {
	Enqueue(url string, depth int) error
	Dequeue() (string, int, error)

	Read() ([]QueueType, error)
	Delete(url string) error
}
type QueueType struct {
	Url   string `json:"url"`
	Depth int    `json:"depth"`
}

// BadgerQueue is an implementation of the Queue interface using BadgerDB.
type BadgerQueue struct {
	db *badger.DB
}

// NewBadgerQueue creates a new BadgerQueue instance.
func NewBadgerQueue(db *badger.DB) *BadgerQueue {
	return &BadgerQueue{db: db}
}

// Enqueue adds a URL to the queue.
func (q *BadgerQueue) Enqueue(url string, depth int) error {
	blockWrite.RLock()
	defer blockWrite.RUnlock()
	key := []byte(fmt.Sprintf("%s:%s", config.QueueName, url))
	return q.db.Update(func(txn *badger.Txn) error {
		return txn.Set(key, []byte(strconv.Itoa(depth)))
	})
}

// Dequeue retrieves and removes a URL from the queue.
func (q *BadgerQueue) Dequeue() (string, int, error) {
	blockWrite.RLock()
	defer blockWrite.RUnlock()

	var url string
	var depth int

	err := q.db.Update(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false // We don't need the values
		it := txn.NewIterator(opts)
		defer it.Close()

		prefix := []byte(config.QueueName)
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			url = string(item.KeyCopy(nil)) // Copy the key to avoid issues
			item.Value(func(val []byte) error {
				depth, _ = strconv.Atoi(string(val))
				return nil
			})
			return txn.Delete(item.Key()) // Remove from queue after retrieval
		}
		return nil // No items in queue
	})

	if err != nil {
		return "", 0, fmt.Errorf("error dequeuing from queue: %v", err)
	}

	if url != "" {
		url = url[len(config.QueueName)+1:]
	}

	return url, depth, nil
}

// Read retrieves all URLs from the queue.
func (q *BadgerQueue) Read() ([]QueueType, error) {
	var urls []QueueType

	err := q.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = true
		it := txn.NewIterator(opts)
		defer it.Close()

		prefix := []byte(config.QueueName)
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			url := string(item.KeyCopy(nil))
			depth := 0
			item.Value(func(val []byte) error {
				depth, _ = strconv.Atoi(string(val))
				return nil
			})
			urls = append(urls, QueueType{Url: url[len(config.QueueName)+1:], Depth: depth})
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error reading from queue: %v", err)
	}

	return urls, nil
}

// Delete removes a URL from the queue.
func (q *BadgerQueue) Delete(url string) error {
	blockWrite.RLock()
	defer blockWrite.RUnlock()
	key := []byte(fmt.Sprintf("%s:%s", config.QueueName, url))
	return q.db.Update(func(txn *badger.Txn) error {
		return txn.Delete(key)
	})
}
