package main

import (
	"github.com/gabrielmoura/WebCrawler/config"
	"github.com/gabrielmoura/WebCrawler/infra/cache"
	"github.com/gabrielmoura/WebCrawler/infra/crawler"
	"github.com/gabrielmoura/WebCrawler/infra/db"
	"github.com/gabrielmoura/WebCrawler/infra/log"
)

func main() {
	log.InitLogger()
	config.LoadConfig()
	db.InitDB()
	cache.InitCache()

	startUrl := "http://wiki.i2p-projekt.i2p/wiki/index.php/Eepsite/Services"
	crawler.Wg.Add(1)
	go crawler.HandlePage(startUrl, 0)
	crawler.Wg.Wait()
}
