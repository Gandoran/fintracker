package web

import (
	"context"
	"encoding/json"
	"fintracker/internal/db"
	"net/http"
	"time"
)

type HealthStatus struct {
	Status    string `json:"status"`
	Database  string `json:"database"`
	Ollama    string `json:"ollama"`
	Timestamp string `json:"timestamp"`
}

type HealthHandler struct {
	store     *db.Store
	ollamaURL string
}

func NewHealthHandler(store *db.Store, ollamaURL string) *HealthHandler {
	return &HealthHandler{
		store:     store,
		ollamaURL: ollamaURL,
	}
}

func (h *HealthHandler) HandleHealthz(w http.ResponseWriter, r *http.Request) {
	status := HealthStatus{
		Status:    "ok",
		Timestamp: time.Now().Format(time.RFC3339),
	}
	status = h.HandleDB(r, status)
	status = h.HandleOllama(r, status)
	w.Header().Set("Content-Type", "application/json")
	if status.Status == "degraded" {
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		w.WriteHeader(http.StatusOK)
	}
	json.NewEncoder(w).Encode(status)
}

func (h *HealthHandler) HandleDB(r *http.Request, status HealthStatus) HealthStatus {
	dbCtx, dbCancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer dbCancel()
	if err := h.store.Ping(dbCtx); err != nil {
		status.Database = "down (" + err.Error() + ")"
		status.Status = "degraded"
	} else {
		status.Database = "ok"
	}
	return status
}

func (h *HealthHandler) HandleOllama(r *http.Request, status HealthStatus) HealthStatus {
	ollamaCtx, ollamaCancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer ollamaCancel()
	req, _ := http.NewRequestWithContext(ollamaCtx, "GET", h.ollamaURL, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != 200 {
		status.Ollama = "down"
		status.Status = "degraded"
	} else {
		status.Ollama = "ok"
		if resp != nil {
			resp.Body.Close()
		}
	}
	return status
}
