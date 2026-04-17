package web

import (
	"database/sql"
	"net/http"

	"fintracker/internal/db"
)

type PageData struct {
	SearchTerm string
	Results    any
}

func (s *AppServer) HandleHome(w http.ResponseWriter, r *http.Request) {
	searchQuery := r.URL.Query().Get("q")
	var results any
	var err error
	if searchQuery != "" {
		res, dbErr := s.store.SearchAnalyses(r.Context(), db.SearchAnalysesParams{
			Column1: sql.NullString{String: searchQuery, Valid: true},
			Column2: sql.NullString{String: searchQuery, Valid: true},
		})
		results = res
		err = dbErr
	} else {
		res, dbErr := s.store.GetRecentAnalyses(r.Context(), 30)
		results = res
		err = dbErr
	}
	if err != nil {
		http.Error(w, "Errore interno del database", http.StatusInternalServerError)
		return
	}
	data := PageData{
		SearchTerm: searchQuery,
		Results:    results,
	}
	homeTmpl.Execute(w, data)
}
