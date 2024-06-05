package crawler

import (
	"fmt"
	"github.com/gabrielmoura/WebCrawler/config"
	"github.com/gabrielmoura/WebCrawler/infra/cache"
	"github.com/gabrielmoura/WebCrawler/infra/log"
	"net/url"
	"strings"
)

// isDenyPostfix checks if the link has a deny postfix
func isDenyPostfix(url string, denySuffixes []string) bool {
	for _, denySuffix := range denySuffixes {
		if strings.HasSuffix(strings.ToLower(url), denySuffix) {
			return true
		}
	}
	return false
}

// isAllowedSchema checks if the link has an acceptable schema
func isAllowedSchema(link string, acceptableSchema []string) bool {
	nLink, err := url.Parse(link)
	if err != nil {
		log.Logger.Debug(fmt.Sprintf("Error parsing link in checking schema: %s", err))
		return false
	}
	for _, schema := range acceptableSchema {
		if nLink.Scheme == schema {
			return true
		}
	}
	return false
}

// isAllowedMIME checks if the link has an acceptable MIME type
func isAllowedMIME(contentType string, allowedMIMEs []string) bool {
	for _, allowedMIME := range allowedMIMEs {
		if strings.Contains(contentType, allowedMIME) {
			return true
		}
	}
	return false
}

// checkTLD checks if the link has an acceptable TLD
func checkTLD(link string) bool {
	if len(config.Conf.Filter.Tlds) > 0 {
		linkUrl, err := url.Parse(link)
		if err != nil {
			return false
		}
		for _, tld := range config.Conf.Filter.Tlds {
			if strings.HasSuffix(linkUrl.Hostname(), tld) {
				return true
			}
		}
		return false
	}
	return true
}

func handleAddToQueue(links []string, depth int) {
	for _, link := range links {
		if checkTLD(link) && isAllowedSchema(link, config.AcceptableSchema) {
			err := cache.AddToQueue(link, depth)
			if err != nil {
				log.Logger.Error(fmt.Sprintf("Error adding link to queue: %s", err))
				return
			}
		}

	}
}

var invalidSchemaErr = fmt.Errorf("invalid schema")

func prepareLink(link string) (*url.URL, error) {
	linkUrl, err := url.Parse(link)
	if err != nil {
		return nil, err
	}

	if linkUrl.Scheme == "" {
		return nil, invalidSchemaErr
	}
	q, _ := url.ParseQuery(linkUrl.RawQuery)
	q.Del("utm_source")
	q.Del("utm_medium")
	q.Del("utm_campaign")
	q.Del("utm_term")
	q.Del("utm_content")
	q.Del("#")
	linkUrl.RawQuery = q.Encode()

	if isDenyPostfix(linkUrl.Path, config.DenySuffixes) {
		return nil, fmt.Errorf("deny postfix")
	}

	return linkUrl, nil
}
func prepareParentLink(parentLink, link string) (*url.URL, error) {

	// Remove o primeiro caractere se for uma barra ou ponto
	if strings.HasPrefix(link, "/") || strings.HasPrefix(link, ".") {
		link = link[1:]
	}

	nURL, err := url.Parse(link)
	if err != nil {
		return nil, err
	}
	if nURL.Path == "" {
		return nil, fmt.Errorf("empty path")
	}

	pURL, err := url.Parse(parentLink)
	if err != nil {
		return nil, err
	}

	nURL.Host = pURL.Host
	nURL.Scheme = pURL.Scheme
	log.Logger.Debug(fmt.Sprintf("New URL: %v\n", nURL))

	return nURL, nil
}
