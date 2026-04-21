package scraper

import (
	"context"
	"net/http"
	"time"

	"github.com/go-shiori/go-readability"
)

func (f *Fetcher) extractFullArticle(ctx context.Context, articleURL string) string {
	reqCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(reqCtx, "GET", articleURL, nil)
	if err != nil {
		return ""
	}
	resp, err := f.httpClient.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return ""
	}
	article, err := readability.FromReader(resp.Body, req.URL)
	if err != nil {
		return ""
	}
	return article.TextContent
}
