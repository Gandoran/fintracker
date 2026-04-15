package scraper

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"fintracker/internal/models"

	"github.com/mmcdole/gofeed"
)

type Fetcher struct {
	parser *gofeed.Parser
}

func NewFetcher() *Fetcher {
	return &Fetcher{
		parser: gofeed.NewParser(),
	}
}

func (f *Fetcher) Fetch(ctx context.Context, urls []string) ([]models.Article, error) {
	var allArticles []models.Article
	for _, url := range urls {
		feed, err := f.parser.ParseURLWithContext(url, ctx)
		if err != nil {
			fmt.Printf("Error on %s: %v\n", url, err)
			continue
		}
		allArticles = append(allArticles, f.parseFeed(feed)...)
	}
	return allArticles, nil
}

func (f *Fetcher) parseFeed(feed *gofeed.Feed) []models.Article {
	var articles []models.Article
	for _, item := range feed.Items {
		rawText := item.Description
		if item.Content != "" {
			rawText = item.Content
		}
		pubDate := time.Now()
		if item.PublishedParsed != nil {
			pubDate = *item.PublishedParsed
		}
		articles = append(articles, models.Article{
			Title:     sanitizeText(item.Title),
			Link:      item.Link,
			Content:   sanitizeText(rawText),
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
