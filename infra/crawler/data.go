package crawler

import (
	"github.com/gabrielmoura/WebCrawler/infra/cache"
	"github.com/gabrielmoura/WebCrawler/infra/data"
	"github.com/gabrielmoura/WebCrawler/infra/db"
	"sync"
)

var visitedMutex sync.Mutex
var pagesMutex sync.Mutex

func SetVisited(url string) {
	visitedMutex.Lock()
	cache.SetVisited(url)
	visitedMutex.Unlock()
}

func GetVisited(url string) bool {
	visitedMutex.Lock()
	defer visitedMutex.Unlock()
	return cache.IsVisited(url)
}

func SetPage(url string, page *data.Page) {
	pagesMutex.Lock()
	err := db.WritePage(page)
	if err != nil {
		return
	}
	//pages[url] = page
	pagesMutex.Unlock()
}

func GetPage(url string) *data.Page {
	pagesMutex.Lock()
	defer pagesMutex.Unlock()
	p, err := db.ReadPage(url)
	if err != nil {
		return nil
	}
	return p
}
