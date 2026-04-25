package scraper

import (
	"context"
	"io"
	"net/http"
	"time"
)

func (f *Fetcher) fetchHTMLFast(ctx context.Context, targetURL string) (string, int) {
	reqCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(reqCtx, "GET", targetURL, nil)
	if err != nil {
		return "", 0
	}
	f.setHeader(req)
	resp, err := f.httpClient.Do(req)
	if err != nil {
		return "", 0
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", resp.StatusCode
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", resp.StatusCode
	}
	return string(bodyBytes), resp.StatusCode
}

func (f *Fetcher) setHeader(req *http.Request) {
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,it;q=0.8")
	req.Header.Set("Referer", "https://www.google.com/")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
}
