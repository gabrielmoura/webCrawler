package crawler

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gabrielmoura/WebCrawler/config"
	"github.com/gabrielmoura/WebCrawler/infra/cache"
	"github.com/gabrielmoura/WebCrawler/infra/log"
	"go.uber.org/zap"
	"golang.org/x/net/html"
	"io"
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

	log.Logger.Info(fmt.Sprintf("Visiting %s", pageUrl))
	plainText, htmlDoc, err := visitLink(pageUrl)
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
	words, _ := countWordsInText(plainText)
	dataPage.Words = words
	dataPage.Url = pageUrl
	dataPage.Links = links
	dataPage.Timestamp = time.Now()
	dataPage.Visited = true

	SetPage(pageUrl, dataPage)

	SetVisited(pageUrl)

	handleAddToQueue(links, depth+1)
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

func visitLink(pageUrl string) ([]byte, *html.Node, error) {
	// Improved error handling using the errors package
	client := httpClient()
	resp, err := client.Get(pageUrl)
	if err != nil {
		return nil, nil, fmt.Errorf("error fetching URL %s: %w", pageUrl, err)
	}
	defer resp.Body.Close()

	// Streamlined MIME type check and early return
	if !isAllowedMIME(resp.Header.Get("Content-Type"), config.AcceptableMimeTypes) {
		return nil, nil, mimeNotAllow
	}

	// Efficiently read the response body into a buffer
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("error reading response body: %w", err)
	}

	// Parse HTML from the buffered content
	htmlDoc, err := html.Parse(bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, nil, fmt.Errorf("error parsing HTML: %w", err)
	}

	return bodyBytes, htmlDoc, nil
}
