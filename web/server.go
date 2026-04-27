package web

import (
	"net/http"

	"fintracker/internal/db"
	"fintracker/internal/scraper"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type AppServer struct {
	dashboardHandler *DashboardHandler
	feedHandler      *FeedHandler
	chatHandler      *ChatHandler
	healthHandler    *HealthHandler
}

func NewAppServer(store *db.Store, ai ChatBot, fetcher *scraper.Fetcher, ollamaUrl string) *AppServer {
	return &AppServer{
		dashboardHandler: NewDashboardHandler(store),
		feedHandler:      NewFeedHandler(store, fetcher),
		chatHandler:      NewChatHandler(store, ai),
		healthHandler:    NewHealthHandler(store, ollamaUrl),
	}
}

func (s *AppServer) RegisterRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.dashboardHandler.HandleHome)
	mux.HandleFunc("/api/chat", s.chatHandler.HandleChatAPI)
	mux.HandleFunc("/feed/discover", s.feedHandler.HandleDiscover)
	mux.HandleFunc("/feed/add-source", s.feedHandler.HandleAddSource)
	mux.HandleFunc("/feed/delete-source", s.feedHandler.HandleDeleteSource)
	mux.HandleFunc("/healthz", s.healthHandler.HandleHealthz)
	mux.Handle("/metrics", promhttp.Handler())
	//DEBUG TODO REMOVE
	mux.HandleFunc("/admin/clean-queue", s.dashboardHandler.HandleCleanQueue)
	return mux
}
