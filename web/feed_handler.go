package web

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"fintracker/internal/db"
	"fintracker/internal/scraper"
)

type FeedHandler struct {
	store   *db.Store
	fetcher *scraper.Fetcher
}

func NewFeedHandler(store *db.Store, fetcher *scraper.Fetcher) *FeedHandler {
	return &FeedHandler{
		store:   store,
		fetcher: fetcher,
	}
}

func (h *FeedHandler) HandleDiscover(w http.ResponseWriter, r *http.Request) {
	targetURL := r.URL.Query().Get("url")
	if targetURL == "" {
		http.Error(w, "URL mancante", http.StatusBadRequest)
		return
	}
	if !strings.HasPrefix(targetURL, "http://") && !strings.HasPrefix(targetURL, "https://") {
		targetURL = "https://" + targetURL
	}
	feeds, err := h.fetcher.DiscoverRSS(r.Context(), targetURL)
	if err != nil {
		fmt.Printf("Discovery Error %s: %v\n", targetURL, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(feeds)
}

func (h *FeedHandler) HandleAddSource(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		return
	}
	var req struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Dati non validi", http.StatusBadRequest)
		return
	}
	_, err := h.store.CreateSource(r.Context(), db.CreateSourceParams{
		Name:     req.Name,
		Url:      req.URL,
		Category: sql.NullString{String: "Finanza", Valid: true},
	})
	if err != nil {
		http.Error(w, "Impossibile salvare nel database: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *FeedHandler) HandleDeleteSource(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Metodo non consentito", http.StatusMethodNotAllowed)
		return
	}
	idStr := r.FormValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "ID non valido", http.StatusBadRequest)
		return
	}
	err = h.store.DeleteSource(r.Context(), id)
	if err != nil {
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
