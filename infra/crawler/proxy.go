package crawler

import (
	"github.com/gabrielmoura/WebCrawler/config"
	"net/http"
	"net/url"
	"time"
)

func ProxyClient() *http.Client {
	urlProxy, _ := url.Parse(config.Conf.Proxy.ProxyURL)
	transport := &http.Transport{
		Proxy: http.ProxyURL(urlProxy),
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second, // Definir um timeout de 5 segundos
	}

	return client
}
