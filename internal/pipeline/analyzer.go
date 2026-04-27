package pipeline

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"fintracker/internal/db"
	"fintracker/internal/models"
)

// TODO REFACTOR
func (w *Worker) processNextPendingArticle() bool {
	dbArt, err := w.store.GetNextPendingArticle(context.Background())
	fmt.Printf("Token usati dall'articolo: %v\n", float64(len(dbArt.Content))/4)
	if err != nil {
		return false
	}
	artModel := models.Article{
		Title:   dbArt.Title,
		Link:    dbArt.Link,
		Content: dbArt.Content,
		Source:  dbArt.Source,
	}
	ollamaCtx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	analysis, err := w.ai.AnalyzeArticle(ollamaCtx, artModel)
	if err != nil {
		w.store.UpdateArticleStatus(context.Background(), db.UpdateArticleStatusParams{
			Status: "FAILED",
			ID:     dbArt.ID,
		})
		return true
	}
	_, err = w.store.CreateAnalysis(context.Background(), db.CreateAnalysisParams{
		ArticleID:        dbArt.ID,
		Summary:          analysis.Summary,
		Sentiment:        analysis.Sentiment,
		Impact:           analysis.Impact,
		Tickers:          strings.Join(analysis.Ticker, ","),
		ReferenceLinks:   strings.Join(analysis.ReferenceLinks, ","),
		ReliabilityScore: int64(analysis.Reliability),
	})
	if err != nil {
		return false
	}
	w.store.UpdateArticleStatus(context.Background(), db.UpdateArticleStatusParams{
		Status: "COMPLETED",
		ID:     dbArt.ID,
	})
	w.SendTelegramNotify(analysis, &artModel)
	return true
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
