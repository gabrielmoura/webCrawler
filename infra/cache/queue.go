package cache

import (
	"fmt"
	"github.com/dgraph-io/badger/v4"
)

var queueName = "queueIndex"
var queue *BadgerQueue

// Queue represents an interface for interacting with a queue data structure.
type Queue interface {
	Enqueue(url string) error
	Dequeue() (string, error)
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
func (q *BadgerQueue) Enqueue(url string) error {
	key := []byte(fmt.Sprintf("%s:%s", queueName, url))

	return q.db.Update(func(txn *badger.Txn) error {
		return txn.Set(key, nil) // Empty value, we only care about keys
	})
}

// Dequeue retrieves and removes a URL from the queue.
func (q *BadgerQueue) Dequeue() (string, error) {
	var url string

	err := q.db.Update(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false // We don't need the values
		it := txn.NewIterator(opts)
		defer it.Close()

		prefix := []byte(queueName)
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			url = string(item.KeyCopy(nil)) // Copy the key to avoid issues
			return txn.Delete(item.Key())   // Remove from queue after retrieval
		}
		return nil // No items in queue
	})

	if err != nil {
		return "", fmt.Errorf("error dequeuing from queue: %v", err)
	}

	if url != "" {
		url = url[len(queueName)+1:]
	}
	return url, nil
}
