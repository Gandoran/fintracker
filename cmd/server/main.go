package main

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"fintracker/internal/config"
	"fintracker/internal/db"
	"fintracker/internal/ollama"
	"fintracker/internal/scraper"
	"fintracker/internal/web"
)

func main() {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("Errore configurazione: %v", err)
	}
	store, err := db.NewStore("fintracker.db")
	if err != nil {
		log.Fatalf("Impossibile connettersi al DB: %v", err)
	}
	fetcher := scraper.NewFetcher()
	ai := ollama.NewClient(cfg.LLM.URL, cfg.LLM.Model, cfg.LLM.Temperature)
	appServer := web.NewAppServer(store)
	go runDaemon(cfg, fetcher, ai, store)
	http.HandleFunc("/", appServer.HandleHome)
	http.ListenAndServe(":8080", nil)
}

func runDaemon(cfg *config.Config, f *scraper.Fetcher, ai *ollama.Client, store *db.Store) {
	ticker := time.NewTicker(time.Duration(cfg.SCRAPER.IntervalMinutes) * time.Minute)
	defer ticker.Stop()
	for {
		log.Println("[DEMONE] Find New Article...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		articles, _ := f.Fetch(ctx, cfg.SCRAPER.Feeds)
		for _, art := range articles {
			savedArticle, err := store.CreateArticle(ctx, db.CreateArticleParams{
				Title:       art.Title,
				Link:        art.Link,
				Content:     art.Content,
				Source:      art.Source,
				PublishedAt: art.Published,
			})
			if err != nil {
				continue
			}
			analysis, err := ai.AnalyzeArticle(ctx, art)
			if err == nil {
				_, err = store.CreateAnalysis(ctx, db.CreateAnalysisParams{
					ArticleID: savedArticle.ID,
					Summary:   analysis.Summary,
					Sentiment: analysis.Sentiment,
					Impact:    analysis.Impact,
					Tickers:   strings.Join(analysis.Ticker, ", "),
				})

				if err != nil {
					log.Printf("Error on db Save: %v", err)
				}
			} else {
				log.Printf("Ollama Error '%s': %v", art.Title, err)
			}
		}
		cancel()
		<-ticker.C
	}
}
