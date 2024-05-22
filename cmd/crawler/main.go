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

	crawler.HandleQueue(config.Conf.InicialURL)
}
