package scraper

import (
	"math/rand"
	"net/http"
)

var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.3 Safari/605.1.15",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:123.0) Gecko/20100101 Firefox/123.0",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:123.0) Gecko/20100101 Firefox/123.0",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Edge/122.0.0.0 Safari/537.36",
}

type randomUATransport struct {
	baseTransport http.RoundTripper
}

func (t *randomUATransport) RoundTrip(req *http.Request) (*http.Response, error) {
	//go require a clone to avoid http request race condition
	clonedReq := req.Clone(req.Context())
	randomIndex := rand.Intn(len(userAgents))
	randomUA := userAgents[randomIndex]
	clonedReq.Header.Set("User-Agent", randomUA)
	clonedReq.Header.Set("Accept", "application/rss+xml, application/xml, application/json, text/xml, */*")
	return t.baseTransport.RoundTrip(clonedReq)
}
