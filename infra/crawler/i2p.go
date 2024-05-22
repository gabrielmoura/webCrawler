package crawler

import (
	"net/http"
	"net/url"
	"time"
)

func I2PClient() *http.Client {
	transport := &http.Transport{
		Proxy: http.ProxyURL(&url.URL{
			Scheme: "http",
			Host:   "localhost:4444",
		}),
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second, // Definir um timeout de 5 segundos
	}

	return client
}
