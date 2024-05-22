package crawler

import (
	"fmt"
	"github.com/gabrielmoura/WebCrawler/config"
	"github.com/gabrielmoura/WebCrawler/infra/data"
	"github.com/gabrielmoura/WebCrawler/infra/i2p"
	"golang.org/x/net/html"
	"net/http"
	"net/url"
	"sync"
	"time"
)

var Wg sync.WaitGroup
var semaphore = make(chan struct{}, *config.MaxConcurrency)

func HandlePage(pageUrl string, depth int) {
	defer Wg.Done()

	if depth > *config.MaxDepth {
		return
	}

	if GetVisited(pageUrl) {
		return
	}
	SetVisited(pageUrl)

	semaphore <- struct{}{}
	defer func() { <-semaphore }()

	fmt.Println("Visiting", pageUrl)
	htmlDoc, err := CheckLink(pageUrl)
	if err != nil {
		fmt.Println("Error checking link:", err)
		return
	}

	links, err := extractLinks(htmlDoc)
	if err != nil {
		fmt.Println("Error extracting links:", err)
		return
	}

	dataPage, err := extractData(htmlDoc)
	if err != nil {
		fmt.Println("Error extracting data:", err)
		return
	}
	dataPage.Url = pageUrl
	dataPage.Links = links
	dataPage.Timestamp = time.Now()
	dataPage.Visited = true

	SetPage(pageUrl, dataPage)

	fmt.Println("Total links", len(links))

	for _, link := range links {
		Wg.Add(1)
		go HandlePage(link, depth+1)
	}
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

func extractLinks(n *html.Node) ([]string, error) {
	var links []string

	var extract func(*html.Node)
	extract = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					urlE, err := url.Parse(a.Val)
					if err != nil || urlE.Scheme == "" {
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
	return html.Parse(resp.Body)
}

func HttpClient() *http.Client {
	if config.Conf.I2PCfg.Enabled {
		return i2p.I2PClient()
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
