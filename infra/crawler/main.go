package crawler

import (
	"errors"
	"fmt"
	"github.com/gabrielmoura/WebCrawler/config"
	"github.com/gabrielmoura/WebCrawler/infra/cache"
	"github.com/gabrielmoura/WebCrawler/infra/data"
	"github.com/gabrielmoura/WebCrawler/infra/log"
	"go.uber.org/zap"
	"golang.org/x/net/html"
	"net/http"
	"sync"
	"time"
)

var Wg sync.WaitGroup
var semaphore = make(chan struct{}, *config.MaxConcurrency)
var mimeNotAllow = errors.New("mime: not allowed")

func HandlePageV2(pageUrl string) {
	// Só processa uma página se ela ainda não foi visitada
	defer Wg.Done()

	if GetVisited(pageUrl) {
		return
	}
	SetVisited(pageUrl)

	log.Logger.Info(fmt.Sprintf("Visiting %s", pageUrl))
	htmlDoc, err := CheckLink(pageUrl)
	if err != nil {
		if err == mimeNotAllow {
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

	handleAddToCache(links)

	log.Logger.Info(fmt.Sprintf("Total links %d", len(links)))
}
func HandleQueue(initialURL string) {
	// Só processa a fila se ela não estiver vazia
	log.Logger.Info("Handling queue")
	ok, err := cache.GetFromQueue()

	if err != nil { // Check if queue is empty
		log.Logger.Info("Queue is empty", zap.Error(err))
	}
	if ok == "" {
		cache.AddToQueue(initialURL)
	}
	loopQueue(0)
}

func loopQueue(depth int) {
	if depth > *config.MaxDepth {
		return
	}
	for {
		links, _ := cache.GetFromQueueV2(*config.MaxConcurrency) // Get a batch of links
		if len(links) == 0 {
			break
		}

		for _, link := range links {
			Wg.Add(1) // Para
			go HandlePageV2(link)
		}
		Wg.Wait()
	}
	//err := cache.OptimizeCache()
	//if err != nil {
	//	log.Logger.Error(fmt.Sprintf("Error optimizing cache: %s", err))
	//}
	loopQueue(depth + 1)
}

func extractData(n *html.Node) (*data.Page, error) {
	var dataPage data.Page

	var extract func(*html.Node)
	extract = func(n *html.Node) {
		if n.Type == html.ElementNode {
			if n.Data == "title" && n.FirstChild != nil {
				dataPage.Title = n.FirstChild.Data
			} else if n.Data == "meta" {
				var isDescription bool
				for _, a := range n.Attr {
					if a.Key == "name" && a.Val == "description" {
						isDescription = true
					}
					if a.Key == "content" {
						if isDescription {
							dataPage.Description = a.Val
						} else {
							dataPage.Meta = append(dataPage.Meta, a.Val)
						}
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extract(c)
		}
	}

	extract(n)
	return &dataPage, nil
}

func extractLinks(parentLink string, n *html.Node) ([]string, error) {
	var links []string

	var extract func(*html.Node)
	extract = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					urlE, err := prepareLink(a.Val)
					if err != nil {
						if errors.Is(invalidSchemaErr, err) {
							preparedLink, err := prepareParentLink(parentLink, a.Val)
							if err != nil {
								continue
							}
							urlE = preparedLink
						}
						log.Logger.Debug(fmt.Sprintf("Error preparing link: %s", err))
						continue
					}
					links = append(links, urlE.String())
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extract(c)
		}
	}

	extract(n)
	return links, nil
}

func CheckLink(pageUrl string) (*html.Node, error) {
	client := HttpClient()
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

func HttpClient() *http.Client {
	if config.Conf.Proxy.Enabled {
		return ProxyClient()
	} else {
		return &http.Client{
			Timeout: 5 * time.Second, // Definir um timeout de 5 segundos
			Transport: &http.Transport{
				MaxIdleConns:        10,
				MaxIdleConnsPerHost: 10,
			},
		}
	}
}
