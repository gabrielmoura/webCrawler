package crawler

import (
	"github.com/gabrielmoura/WebCrawler/infra/cache"
	"github.com/gabrielmoura/WebCrawler/infra/data"
	"github.com/gabrielmoura/WebCrawler/infra/db"
	"sync"
)

var visitedMutex sync.Mutex
var pagesMutex sync.Mutex

// SetVisited adds a URL to the cache to mark it as visited.
func SetVisited(url string) {
	visitedMutex.Lock()
	cache.SetVisited(url)
	visitedMutex.Unlock()
}

// GetVisited retrieves a URL from the cache to check if it has been visited.
func GetVisited(url string) bool {
	visitedMutex.Lock()
	defer visitedMutex.Unlock()
	return cache.IsVisited(url)
}

// SetPage adds a page to the database.
func SetPage(url string, page *data.Page) {
	//pagesMutex.Lock()
	err := db.WritePage(page)
	if err != nil {
		return
	}
	//pages[url] = page
	//pagesMutex.Unlock()
}

// GetPage retrieves a page from the database.
func GetPage(url string) *data.Page {
	//pagesMutex.Lock()
	//defer pagesMutex.Unlock()
	p, err := db.ReadPage(url)
	if err != nil {
		return nil
	}
	return p
}
