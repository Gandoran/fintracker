package web

import (
	"net/http"

	"fintracker/internal/db"
	"fintracker/internal/scraper"
)

type AppServer struct {
	dashboardHandler *DashboardHandler
	feedHandler      *FeedHandler
	chatHandler      *ChatHandler
}

func NewAppServer(store *db.Store, ai ChatBot, fetcher *scraper.Fetcher) *AppServer {
	return &AppServer{
		dashboardHandler: NewDashboardHandler(store),
		feedHandler:      NewFeedHandler(store, fetcher),
		chatHandler:      NewChatHandler(store, ai),
	}
}

func (s *AppServer) RegisterRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.dashboardHandler.HandleHome)
	mux.HandleFunc("/api/chat", s.chatHandler.HandleChatAPI)
	mux.HandleFunc("/feed/discover", s.feedHandler.HandleDiscover)
	mux.HandleFunc("/feed/add-source", s.feedHandler.HandleAddSource)
	mux.HandleFunc("/feed/delete-source", s.feedHandler.HandleDeleteSource)
	//DEBUG TODO REMOVE
	mux.HandleFunc("/admin/clean-queue", s.dashboardHandler.HandleCleanQueue)
	return mux
}
