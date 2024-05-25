package crawler

import (
	"errors"
	"fmt"
	"github.com/gabrielmoura/WebCrawler/infra/data"
	"github.com/gabrielmoura/WebCrawler/infra/log"
	"golang.org/x/net/html"
)

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
