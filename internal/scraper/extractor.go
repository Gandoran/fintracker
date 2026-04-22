package scraper

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/go-shiori/go-readability"
)

func (f *Fetcher) extractFullArticle(ctx context.Context, articleURL string) string {
	reqCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(reqCtx, "GET", articleURL, nil)
	if err != nil {
		return ""
	}
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,it;q=0.8")
	req.Header.Set("Referer", "https://www.google.com/")
	req.Header.Set("Connection", "keep-alive")
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
	return f.CheckText(article)
}

func (f *Fetcher) CheckText(article readability.Article) string {
	extractedText := strings.TrimSpace(article.TextContent)
	lowerText := strings.ToLower(extractedText)
	if len(extractedText) < 250 {
		return ""
	}
	badKeywords := []string{
		"enable javascript", "please wait while we verify", "are you a robot", "accept cookies", "javascript is disabled", "turn on javascript", "bloomberg.com needs to review the security",
	}
	for _, badWord := range badKeywords {
		if strings.Contains(lowerText, badWord) {
			return ""
		}
	}
	return extractedText
}
