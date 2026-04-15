package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"fintracker/internal/config"
	"fintracker/internal/ollama"
	"fintracker/internal/scraper"
	"fintracker/internal/web"
)

func main() {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("Errore configurazione: %v", err)
	}
	fetcher := scraper.NewFetcher()
	ai := ollama.NewClient(cfg.LLM.URL, cfg.LLM.Model, cfg.LLM.Temperature)
	appServer := web.NewAppServer()
	go runDaemon(cfg, fetcher, ai, appServer)
	http.HandleFunc("/", appServer.HandleHome)
	http.ListenAndServe(":8080", nil)
}

func runDaemon(cfg *config.Config, f *scraper.Fetcher, ai *ollama.Client, webServer *web.AppServer) {
	ticker := time.NewTicker(time.Duration(cfg.SCRAPER.IntervalMinutes) * time.Minute)
	defer ticker.Stop()
	for {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		articles, _ := f.Fetch(ctx, cfg.SCRAPER.Feeds)
		for _, art := range articles {
			log.Printf("Analyzing: %s", art.Title)
			analysis, err := ai.AnalyzeArticle(ctx, art)
			if err == nil {
				webServer.AddResult(*analysis)
			} else {
				log.Printf("Ollama error on '%s': %v", art.Title, err)
			}
		}
		cancel()
		<-ticker.C
	}
}
