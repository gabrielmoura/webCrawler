package extract

import (
	"github.com/gabrielmoura/WebCrawler/infra/data"
	"golang.org/x/net/html"
)

func extractTitle(n *html.Node, dataPage *data.Page) {
	if n.Type == html.ElementNode && n.Data == "title" && n.FirstChild != nil {
		dataPage.Title = n.FirstChild.Data
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		extractTitle(c, dataPage)
	}
}

func extractDescription(n *html.Node, dataPage *data.Page) {
	if n.Type == html.ElementNode && n.Data == "meta" {
		var isDescription bool
		var content string
		for _, a := range n.Attr {
			if a.Key == "name" && a.Val == "description" {
				isDescription = true
			}
			if a.Key == "content" {
				content = a.Val
			}
		}
		if isDescription {
			dataPage.Description = content
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		extractDescription(c, dataPage)
	}
}

func extractMeta(n *html.Node, dataPage *data.Page) {
	if n.Type == html.ElementNode && n.Data == "meta" {
		var content string
		var isDescription bool
		for _, a := range n.Attr {
			if a.Key == "name" && a.Val == "description" {
				isDescription = true
				break
			}
			if a.Key == "content" {
				content = a.Val
			}
		}
		if !isDescription && content != "" {
			dataPage.Meta = append(dataPage.Meta, content)
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		extractMeta(c, dataPage)
	}
}

func ExtractData(n *html.Node) (*data.Page, error) {
	var dataPage data.Page

	extractTitle(n, &dataPage)
	extractDescription(n, &dataPage)
	extractMeta(n, &dataPage)

	return &dataPage, nil
}
