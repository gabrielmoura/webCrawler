package crawler

import (
	"errors"
	"fmt"
	"github.com/gabrielmoura/WebCrawler/config"
	"github.com/gabrielmoura/WebCrawler/infra/cache"
	"github.com/gabrielmoura/WebCrawler/infra/log"
	"go.uber.org/zap"
	"golang.org/x/net/html"
	"net/http"
	"sync"
	"time"
)

var Wg sync.WaitGroup
var mimeNotAllow = errors.New("mime: not allowed")

// processPage processa uma página, extrai links e dados
func processPage(pageUrl string, depth int) {

	log.Logger.Debug(fmt.Sprintf("Looping queue, depth: %d", depth))
	if depth > *config.MaxDepth {
		log.Logger.Info(fmt.Sprintf("Reached max depth of %d, %d", *config.MaxDepth, depth))
		return
	}
	// Só processa uma página se ela ainda não foi visitada

	if GetVisited(pageUrl) {
		return
	}
	SetVisited(pageUrl)

	log.Logger.Info(fmt.Sprintf("Visiting %s", pageUrl))
	htmlDoc, err := visitLink(pageUrl)
	if err != nil {
		if errors.Is(err, mimeNotAllow) {
			//log.Logger.Info(fmt.Sprintf("MIME not allowed: %s", pageUrl))
			return
		}
		log.Logger.Debug(fmt.Sprintf("Error checking link: %s", err))
		return
	}

	links, err := extractLinks(pageUrl, htmlDoc)
	if err != nil {
		log.Logger.Error(fmt.Sprintf("Error extracting links: %s", err))
		return
	}

	dataPage, err := extractData(htmlDoc)
	if err != nil {
		log.Logger.Error(fmt.Sprintf("Error extracting data: %s", err))
		return
	}
	dataPage.Url = pageUrl
	dataPage.Links = links
	dataPage.Timestamp = time.Now()
	dataPage.Visited = true

	SetPage(pageUrl, dataPage)

	handleAddToQueue(links, depth+1)

	log.Logger.Info(fmt.Sprintf("Total links %d", len(links)))
}
func HandleQueue(initialURL string) {
	// Só processa a fila se ela não estiver vazia
	log.Logger.Info("Handling queue")
	ok, _, err := cache.GetFromQueue()

	if err != nil { // Check if queue is empty
		log.Logger.Info("Queue is empty", zap.Error(err))
	}
	if ok == "" {
		cache.AddToQueue(initialURL, 0)
	}
	loopQueue()
}

func loopQueue() {
	for {
		links, _ := cache.GetFromQueueV2(*config.MaxConcurrency) // Get a batch of links
		if len(links) == 0 {
			break
		}

		for _, link := range links {
			if link.Url == "" {
				continue
			}
			Wg.Add(1) // Para cada link, incrementa o WaitGroup

			go func(link cache.QueueType) {
				defer Wg.Done()
				processPage(link.Url, link.Depth)
			}(link)
		}
		Wg.Wait()
	}
}

func visitLink(pageUrl string) (*html.Node, error) {
	client := httpClient()
	resp, err := client.Get(pageUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}
	contentType := resp.Header.Get("Content-Type")
	if !isAllowedMIME(contentType, acceptableMimeTypes) {
		resp.Body.Close() // Fechar o corpo da resposta
		return nil, mimeNotAllow
	}
	return html.Parse(resp.Body)
}
