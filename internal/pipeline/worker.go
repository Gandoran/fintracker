package pipeline

import (
	"context"
	"fintracker/internal/config"
	"fintracker/internal/db"
	"fintracker/internal/ollama"
	"fintracker/internal/scraper"
	"fintracker/notifier"
	"time"
)

type Worker struct {
	cfg     *config.Config
	fetcher *scraper.Fetcher
	ai      *ollama.AnalyzerClient
	store   *db.Store
	bot     *notifier.TelegramBot
}

func NewWorker(cfg *config.Config, f *scraper.Fetcher, ai *ollama.AnalyzerClient, store *db.Store, bot *notifier.TelegramBot) *Worker {
	return &Worker{cfg: cfg, fetcher: f, ai: ai, store: store, bot: bot}
}

func (w *Worker) Start() {
	go w.runScraperLoop()
	w.runAnalyzerLoop()
}

func (w *Worker) runScraperLoop() {
	ticker := time.NewTicker(time.Duration(w.cfg.SCRAPER.IntervalMinutes) * time.Minute)
	defer ticker.Stop()
	w.fetchAndSave()
	for {
		<-ticker.C
		w.fetchAndSave()
	}
}

func (w *Worker) runAnalyzerLoop() {
	for {
		processed := w.processNextPendingArticle()
		if !processed {
			time.Sleep(10 * time.Second)
		} else {
			time.Sleep(3 * time.Second)
		}
	}
}

// DEBUG TODO REMOVE
func (w *Worker) CleanupQueue() {
	err := w.store.DeleteAllPendingArticles(context.Background())
	if err != nil {
	}
}
