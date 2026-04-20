package pipeline

import (
	"context"
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
	ticker := time.NewTicker(time.Duration(w.cfg.SCRAPER.IntervalMinutes) * time.Minute)
	defer ticker.Stop()
	for {
		log.Println("[PIPELINE] Finding new articles...")
		w.processFeeds()
		<-ticker.C
	}
}

func (w *Worker) processFeeds() {
	fetchCtx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	articles, err := w.fetcher.Fetch(fetchCtx, w.cfg.SCRAPER.Feeds)
	if err != nil {
		log.Printf("fetch Error: %v", err)
		return
	}
	for _, art := range articles {
		w.processSingleArticle(art)
		time.Sleep(2 * time.Second) // Cooldown GPU
	}
}

func (w *Worker) processSingleArticle(art models.Article) {
	savedArticle, err := w.store.CreateArticle(context.Background(), db.CreateArticleParams{
		Title:       art.Title,
		Link:        art.Link,
		Content:     art.Content,
		Source:      art.Source,
		PublishedAt: art.Published,
	})
	if err != nil {
		return
	}
	ollamaCtx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	analysis, err := w.ai.AnalyzeArticle(ollamaCtx, art)
	if err != nil {
		return
	}
	_, err = w.store.CreateAnalysis(context.Background(), db.CreateAnalysisParams{
		ArticleID:        savedArticle.ID,
		Summary:          analysis.Summary,
		Sentiment:        analysis.Sentiment,
		Impact:           analysis.Impact,
		Tickers:          strings.Join(analysis.Ticker, ", "),
		ReferenceLinks:   strings.Join(analysis.ReferenceLinks, ","),
		ReliabilityScore: int64(analysis.Reliability),
	})
	if err != nil {
		log.Printf("Error saving on DB: %v", err)
	} else {
		log.Println("Analys Saved")
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
