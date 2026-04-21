package scraper

import (
	"context"
	"regexp"
	"strings"
	"time"

	"fintracker/internal/models"

	"github.com/mmcdole/gofeed"
)

func (f *Fetcher) parseFeed(ctx context.Context, feed *gofeed.Feed) []models.Article {
	var articles []models.Article
	for _, item := range feed.Items {
		if ctx.Err() != nil {
			break
		}
		pubDate := time.Now()
		if item.PublishedParsed != nil {
			pubDate = *item.PublishedParsed
		}
		articleContent := item.Description
		if item.Content != "" {
			articleContent = item.Content
		}
		fullText := f.extractFullArticle(ctx, item.Link)
		if fullText != "" {
			articleContent = fullText
		}
		articles = append(articles, models.Article{
			Title:     sanitizeText(item.Title),
			Link:      item.Link,
			Content:   sanitizeText(articleContent),
			Published: pubDate,
			Source:    feed.Title,
		})
	}
	return articles
}

func sanitizeText(htmlStr string) string {
	re := regexp.MustCompile(`<[^>]*>`)
	cleaned := re.ReplaceAllString(htmlStr, " ")
	cleaned = strings.Join(strings.Fields(cleaned), " ")
	return cleaned
}
