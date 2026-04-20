package web

import (
	"database/sql"
	"net/http"

	"fintracker/internal/db"
)

type UIHandler struct {
	store *db.Store
}

func NewUiHandler(store *db.Store) *UIHandler {
	return &UIHandler{store: store}
}

func (h *UIHandler) HandleHome(w http.ResponseWriter, r *http.Request) {
	searchQuery := r.URL.Query().Get("q")
	var results any
	var err error
	if searchQuery != "" {
		res, dbErr := h.store.SearchAnalyses(r.Context(), db.SearchAnalysesParams{
			Column1: sql.NullString{String: searchQuery, Valid: true},
			Column2: sql.NullString{String: searchQuery, Valid: true},
		})
		results = res
		err = dbErr
	} else {
		res, dbErr := h.store.GetRecentAnalyses(r.Context(), 30)
		results = res
		err = dbErr
	}
	if err != nil {
		return
	}
	data := struct {
		SearchTerm string
		Results    any
	}{
		SearchTerm: searchQuery,
		Results:    results,
	}
	tmpl.ExecuteTemplate(w, "base", data)
}
