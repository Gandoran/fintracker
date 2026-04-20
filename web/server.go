package web

import (
	"net/http"

	"fintracker/internal/db"
)

type AppServer struct {
	uiHandler   *UIHandler
	chatHandler *ChatHandler
}

func NewAppServer(store *db.Store, ai ChatBot) *AppServer {
	return &AppServer{
		uiHandler:   NewUiHandler(store),
		chatHandler: NewChatHandler(store, ai),
	}
}

func (s *AppServer) RegisterRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.uiHandler.HandleHome)
	mux.HandleFunc("/api/chat", s.chatHandler.HandleChatAPI)
	return mux
}
