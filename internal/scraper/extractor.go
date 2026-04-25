package scraper

import (
	"context"
	"net/url"
	"strings"

	"github.com/go-shiori/go-readability"
)

func (f *Fetcher) extractFullArticle(ctx context.Context, articleURL string) string {
	html, status := f.fetchHTMLFast(ctx, articleURL)
	testo := f.parseHTMLToArticle(html, articleURL)
	if status == 200 && testo != "" {
		return testo
	}
	html = f.fetchHTMLSlow(ctx, articleURL)
	testo = f.parseHTMLToArticle(html, articleURL)
	return testo
}

func (f *Fetcher) parseHTMLToArticle(html string, articleURL string) string {
	if html == "" {
		return ""
	}
	parsedURL, _ := url.Parse(articleURL)
	article, err := readability.FromReader(strings.NewReader(html), parsedURL)
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
