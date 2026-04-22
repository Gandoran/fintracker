package web

import (
	"context"
	"database/sql"
	"net/http"

	"fintracker/internal/db"
	"fintracker/web/views"
)

type DashboardHandler struct {
	store *db.Store
}

func NewDashboardHandler(store *db.Store) *DashboardHandler {
	return &DashboardHandler{store: store}
}

func (h *DashboardHandler) HandleHome(w http.ResponseWriter, r *http.Request) {
	searchQuery := r.URL.Query().Get("q")
	dateQuery := r.URL.Query().Get("date")
	ctx := r.Context()
	uiArticles := h.fetchAnalyses(ctx, dateQuery, searchQuery)
	sources, _ := h.store.GetAllSources(ctx)
	data := views.DashboardData{
		SearchTerm:   searchQuery,
		SelectedDate: dateQuery,
		Analyses:     uiArticles,
		Sources:      sources,
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	views.Dashboard(data).Render(ctx, w)
}

func (h *DashboardHandler) fetchAnalyses(ctx context.Context, dateQuery, searchQuery string) []views.UIArticle {
	var uiArticles []views.UIArticle
	//date search
	if dateQuery != "" {
		rows, err := h.store.GetAnalysesByDate(ctx, dateQuery)
		if err == nil {
			for _, row := range rows {
				uiArticles = append(uiArticles, mapToUI(row.Article, row.Analysis))
			}
		}
		return uiArticles
	}
	//text search
	if searchQuery != "" {
		rows, err := h.store.SearchAnalyses(ctx, db.SearchAnalysesParams{
			Column1: sql.NullString{String: searchQuery, Valid: true},
			Column2: sql.NullString{String: searchQuery, Valid: true},
		})
		if err == nil {
			for _, row := range rows {
				uiArticles = append(uiArticles, mapToUI(row.Article, row.Analysis))
			}
		}
		return uiArticles
	}
	rows, err := h.store.GetRecentAnalyses(ctx, 30)
	if err == nil {
		for _, row := range rows {
			uiArticles = append(uiArticles, mapToUI(row.Article, row.Analysis))
		}
	}
	return uiArticles
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
		Link:             art.Link,
	}
}

// TODO REMOVE - debugg
func (h *DashboardHandler) HandleCleanQueue(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Metodo non consentito", http.StatusMethodNotAllowed)
		return
	}
	err := h.store.DeleteAllPendingArticles(r.Context())
	if err != nil {
		http.Error(w, "Errore durante la pulizia della coda: "+err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
