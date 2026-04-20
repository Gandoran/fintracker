package web

import (
	"database/sql"
	"net/http"

	"fintracker/internal/db"
	"fintracker/web/views"
)

type UIHandler struct {
	store *db.Store
}

func NewUiHandler(store *db.Store) *UIHandler {
	return &UIHandler{store: store}
}

func (h *UIHandler) HandleHome(w http.ResponseWriter, r *http.Request) {
	searchQuery := r.URL.Query().Get("q")
	dateQuery := r.URL.Query().Get("date")
	var uiArticles []views.UIArticle
	ctx := r.Context()
	if dateQuery != "" {
		rows, err := h.store.GetAnalysesByDate(ctx, dateQuery)
		if err == nil {
			for _, row := range rows {
				uiArticles = append(uiArticles, mapToUI(row.Article, row.Analysis))
			}
		}
	} else if searchQuery != "" {
		rows, err := h.store.SearchAnalyses(ctx, db.SearchAnalysesParams{
			Column1: sql.NullString{String: searchQuery, Valid: true},
			Column2: sql.NullString{String: searchQuery, Valid: true},
		})
		if err == nil {
			for _, row := range rows {
				uiArticles = append(uiArticles, mapToUI(row.Article, row.Analysis))
			}
		}
	} else {
		rows, err := h.store.GetRecentAnalyses(ctx, 30)
		if err == nil {
			for _, row := range rows {
				uiArticles = append(uiArticles, mapToUI(row.Article, row.Analysis))
			}
		}
	}
	data := views.DashboardData{
		SearchTerm:   searchQuery,
		SelectedDate: dateQuery,
		Analyses:     uiArticles,
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	views.Dashboard(data).Render(r.Context(), w)
}

func mapToUI(art db.Article, an db.Analysis) views.UIArticle {
	return views.UIArticle{
		ID:               art.ID,
		Title:            art.Title,
		Tickers:          an.Tickers,
		Sentiment:        an.Sentiment,
		ReliabilityScore: an.ReliabilityScore,
		Summary:          an.Summary,
		Impact:           an.Impact,
		AnalyzedAt:       an.AnalyzedAt.Time.Format("02/01/2006 15:04"),
	}
}
