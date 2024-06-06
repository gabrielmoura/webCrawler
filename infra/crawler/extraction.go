package crawler

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gabrielmoura/WebCrawler/config"
	"github.com/gabrielmoura/WebCrawler/infra/data"
	"github.com/gabrielmoura/WebCrawler/infra/log"
	"go.uber.org/zap"
	"golang.org/x/net/html"
	"regexp"
)

// countWordsInText Extrai e conta a frequência de palavras do conteúdo HTML, ignorando palavras irrelevantes comuns.
func countWordsInText(data []byte) (map[string]int, error) {
	log.Logger.Debug("Word Count")
	// Etapa 1: Ignorar determinadas tags HTML
	htmlRegex := regexp.MustCompile("(?s)<(script|style|noscript|link|meta)[^>]*?>.*?</(script|style|noscript|link|meta)>")
	parcialPlainText := htmlRegex.ReplaceAll(data, []byte(""))

	// Etapa 2: remover tags HTML
	tagsRegex := regexp.MustCompile("<([^>]*)>")
	plainText := tagsRegex.ReplaceAll(parcialPlainText, []byte(""))

	// Etapa 3: Normalizar texto
	normalizedText := bytes.ToLower(plainText)

	// Etapa 4: Remova caracteres especiais e divida em palavras
	wordRegex := regexp.MustCompile("[^\\pL\\pN\\pZ'-]+")
	noSpecialCh := wordRegex.ReplaceAll(normalizedText, []byte(" "))
	words := bytes.Split(noSpecialCh, []byte(" "))

	// Etapa 5: Conte a frequência das palavras (ignorando palavras comuns)
	wordCounts := make(map[string]int)
	for _, wordBytes := range words {
		word := string(bytes.TrimSpace(wordBytes))

		// Pule palavras curtas e palavras de parada comuns
		if len(word) < 2 || containsMap(config.CommonStopWords, word) {
			continue
		}
		wordCounts[word]++
		log.Logger.Debug("Word: ", zap.Int(word, wordCounts[word]))
	}

	return wordCounts, nil
}

// ContainsMap Verifica se uma palavra está em uma lista de stop words comuns,
func containsMap(wordMap map[string][]string, item string) bool {
	for key, slice := range wordMap {
		// Ignora a primeira string do mapa (chave vazia ou primeira chave lexicograficamente)
		if key == "" {
			continue
		}

		for _, a := range slice {
			if a == item {
				return true
			}
		}
	}
	return false
}

func extractData(n *html.Node) (*data.Page, error) {
	var dataPage data.Page

	var extract func(*html.Node)
	extract = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "title":
				extractTitle(n, &dataPage)
			case "meta":
				extractMeta(n, &dataPage)
			case "script":
				extractJSONLD(n, &dataPage)
			}

		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extract(c)
		}
	}

	extract(n)
	return &dataPage, nil
}

func extractDescription(n *html.Node) string {
	for _, a := range n.Attr {
		if a.Key == "name" && a.Val == "description" {
			for _, a := range n.Attr {
				if a.Key == "content" {
					return a.Val
				}
			}
		}
	}
	return ""
}

func extractOG(n *html.Node) map[string]string {
	ogData := make(map[string]string)
	for _, a := range n.Attr {
		if a.Key == "property" && len(a.Val) > 3 && a.Val[:3] == "og:" {
			for _, a := range n.Attr {
				if a.Key == "content" {
					ogData[a.Val] = a.Val
				}
			}
		}
	}
	return ogData
}

func extractKeywords(n *html.Node) []string {
	for _, a := range n.Attr {
		if a.Key == "name" && a.Val == "keywords" {
			for _, a := range n.Attr {
				if a.Key == "content" {
					return append([]string{}, a.Val)
				}
			}
		}
	}
	return nil
}

func extractManifest(n *html.Node) string {
	if n.Data == "link" {
		var isManifest bool
		for _, a := range n.Attr {
			if a.Key == "rel" && a.Val == "manifest" {
				isManifest = true
			}
			if a.Key == "href" && isManifest {
				return a.Val
			}
		}
	}
	return ""
}
func extractTitle(n *html.Node, dataPage *data.Page) {
	if n.Data == "title" && n.FirstChild != nil {
		dataPage.Title = n.FirstChild.Data
	}
}
func extractMeta(n *html.Node, dataPage *data.Page) {

	if n.Data == "meta" {
		if dataPage.Meta == nil {
			dataPage.Meta = metaNull
		}
		description := extractDescription(n)
		if description != "" {
			dataPage.Description = description
		}

		ogData := extractOG(n)
		for k, v := range ogData {
			dataPage.Meta.OG[k] = v
		}

		keywords := extractKeywords(n)
		if keywords != nil {
			dataPage.Meta.Keywords = append(dataPage.Meta.Keywords, keywords...)
		}
	}

	manifest := extractManifest(n)
	if manifest != "" {
		dataPage.Meta.Manifest = manifest
	}
}

// extractLinks Extrai links de um documento HTML.
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
func extractJSONLD(n *html.Node, dataPage *data.Page) {
	for _, a := range n.Attr {
		if a.Key == "type" && a.Val == "application/ld+json" {
			if dataPage.Meta == nil {
				dataPage.Meta = metaNull
			}
			if n.FirstChild != nil {
				content := n.FirstChild.Data
				dataPage.Meta.Ld = content
			}
		}
	}
}

var metaNull = &data.MetaData{
	OG:       make(map[string]string),
	Keywords: []string{},
	Manifest: "",
	Ld:       "",
}
