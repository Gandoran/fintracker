package scraper

import (
	"context"
	"fmt"
	"net/http"

	"fintracker/internal/models"

	"github.com/mmcdole/gofeed"
)

type Fetcher struct {
	parser     *gofeed.Parser
	httpClient *http.Client
}

func NewFetcher(client *http.Client) *Fetcher {
	fp := gofeed.NewParser()
	fp.Client = client
	return &Fetcher{
		parser:     fp,
		httpClient: client,
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
		allArticles = append(allArticles, f.parseFeed(ctx, feed)...)
	}
	return allArticles, nil
}
