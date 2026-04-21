package pipeline

import (
	"context"
	"fmt"
	"time"

	"fintracker/internal/db"
)

func (w *Worker) fetchAndSave() {
	fetchCtx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	sources, err := w.store.GetActiveSources(context.Background())
	if err != nil || len(sources) == 0 {
		fmt.Print("error!")
		return
	}
	var urls []string
	for _, s := range sources {
		urls = append(urls, s.Url)
	}
	articles, _ := w.fetcher.Fetch(fetchCtx, urls)
	for _, art := range articles {
		w.store.CreateArticle(context.Background(), db.CreateArticleParams{
			Title:       art.Title,
			Link:        art.Link,
			Content:     art.Content,
			Source:      art.Source,
			PublishedAt: art.Published,
		})
	}
}
