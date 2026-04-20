package web

import (
	"context"
	"encoding/json"
	"fintracker/internal/db"
	"net/http"
)

type ChatBot interface {
	ChatWithArticle(ctx context.Context, articleContent string, userQuestion string) (string, error)
}

type ChatHandler struct {
	store *db.Store
	ai    ChatBot
}

func NewChatHandler(store *db.Store, ai ChatBot) *ChatHandler {
	return &ChatHandler{store: store, ai: ai}
}

type ChatRequestPayload struct {
	ArticleID int64  `json:"article_id"`
	Question  string `json:"question"`
}

func (h *ChatHandler) HandleChatAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		return
	}
	var payload ChatRequestPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	art, err := h.store.GetArticleByID(r.Context(), payload.ArticleID)
	if err != nil {
		http.Error(w, "Article not found", http.StatusNotFound)
		return
	}
	answer, err := h.ai.ChatWithArticle(r.Context(), art.Content, payload.Question)
	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"answer": answer})
}
