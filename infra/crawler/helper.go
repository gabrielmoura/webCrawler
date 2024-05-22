package crawler

import (
	"fmt"
	"net/url"
	"strings"
)

var denySuffixes = []string{
	".css",
	".js",
	".png",
	".jpg",
	".jpeg",
	".gif",
	".svg",
	".ico",
	".mp4",
	".mp3",
	".avi",
	".flv",
	".mpeg",
	".webp",
	".webm",
	".woff",
	".woff2",
	".ttf",
	".eot",
	".otf",
	".pdf",
	".zip",
	".tar",
	".gz",
	".bz2",
	".xz",
	".7z",
	".rar",
	".apk",
	".exe",
	".dmg",
	".img",
}

func isDenyPostfix(url string, denySuffixes []string) bool {
	for _, denySuffix := range denySuffixes {
		if strings.HasSuffix(url, denySuffix) {
			return true
		}
	}
	return false
}

var acceptableMimeTypes = []string{
	"text/html",
	"text/plain",
	"text/xml",
	"application/xml",
	"application/xhtml+xml",
	"application/rss+xml",
	"application/atom+xml",
	"application/rdf+xml",
	"application/json",
	"application/ld+json",
	"application/vnd.geo+json",
	"application/xml-dtd",
	"application/rss+xml",
	"application/atom+xml",
	"application/rdf+xml",
	"application/json",
	"application/ld+json",
	"application/vnd.geo+json",
}

func isAllowedMIME(contentType string, allowedMIMEs []string) bool {
	for _, allowedMIME := range allowedMIMEs {
		if strings.Contains(contentType, allowedMIME) {
			return true
		}
	}
	return false
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

	if isDenyPostfix(linkUrl.Path, denySuffixes) {
		return nil, fmt.Errorf("deny postfix")
	}

	return linkUrl, nil
}
func prepareParentLink(parentLink, link string) (string, error) {
	pURL, err := url.Parse(parentLink)
	if err != nil {
		return "", err
	}
	nURL, err := url.Parse(link)
	if err != nil {
		return "", err
	}

	nURL.Host = pURL.Host
	nURL.Scheme = pURL.Scheme

	return nURL.String(), nil
}
