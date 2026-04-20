package pipeline

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"fintracker/internal/config"
	"fintracker/internal/db"
	"fintracker/internal/models"
	"fintracker/internal/ollama"
	"fintracker/internal/scraper"
	"fintracker/notifier"
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

func (w *Worker) fetchAndSave() {
	fetchCtx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	articles, err := w.fetcher.Fetch(fetchCtx, w.cfg.SCRAPER.Feeds)
	if err != nil {
		return
	}
	savedCount := 0
	for _, art := range articles {
		_, err := w.store.CreateArticle(context.Background(), db.CreateArticleParams{
			Title:       art.Title,
			Link:        art.Link,
			Content:     art.Content,
			Source:      art.Source,
			PublishedAt: art.Published,
		})

		if err == nil {
			savedCount++
		}
	}
}

func (w *Worker) runAnalyzerLoop() {
	for {
		dbArt, err := w.store.GetNextPendingArticle(context.Background())
		if err != nil {
			if err == sql.ErrNoRows {
				time.Sleep(10 * time.Second)
				continue
			}
			time.Sleep(10 * time.Second)
			continue
		}
		artModel := models.Article{
			Title:   dbArt.Title,
			Link:    dbArt.Link,
			Content: dbArt.Content,
			Source:  dbArt.Source,
		}
		ollamaCtx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		analysis, err := w.ai.AnalyzeArticle(ollamaCtx, artModel)
		cancel()
		if err != nil {
			w.store.UpdateArticleStatus(context.Background(), db.UpdateArticleStatusParams{
				Status: "FAILED",
				ID:     dbArt.ID,
			})
			continue
		}
		_, err = w.store.CreateAnalysis(context.Background(), db.CreateAnalysisParams{
			ArticleID:        dbArt.ID,
			Summary:          analysis.Summary,
			Sentiment:        analysis.Sentiment,
			Impact:           analysis.Impact,
			Tickers:          strings.Join(analysis.Ticker, ", "),
			ReferenceLinks:   strings.Join(analysis.ReferenceLinks, ","),
			ReliabilityScore: int64(analysis.Reliability),
		})
		if err != nil {
			continue
		}
		w.store.UpdateArticleStatus(context.Background(), db.UpdateArticleStatusParams{
			Status: "COMPLETED",
			ID:     dbArt.ID,
		})
		w.SendTelegramNotify(analysis, &artModel)
		time.Sleep(3 * time.Second)
	}
}

func (w *Worker) SendTelegramNotify(analysis *models.Analysis, art *models.Article) {
	if analysis.Sentiment != "Neutral" && analysis.Reliability >= 0 {
		icon := "📈"
		if analysis.Sentiment == "Bearish" {
			icon = "📉"
		}
		msg := fmt.Sprintf(
			"🚨 <b>Lumina AI Alert</b> %s\n\n"+
				"<b>Aziende:</b> %s\n"+
				"<b>Sentiment:</b> %s\n"+
				"<b>Affidabilità:</b> %d/10\n\n"+
				"<b>Riassunto:</b>\n%s\n\n"+
				"<a href=\"%s\">📰 Leggi articolo originale</a>",
			icon, strings.Join(analysis.Ticker, ", "), analysis.Sentiment, analysis.Reliability, analysis.Summary, art.Link,
		)
		err := w.bot.SendAlert(msg)
		if err != nil {
			log.Printf("Error on Telegram: %v", err)
		}
	}
}
