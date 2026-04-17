package web

import (
	"net/http"

	"fintracker/internal/db"
)

type AppServer struct {
	store *db.Store
}

func NewAppServer(store *db.Store) *AppServer {
	return &AppServer{store: store}
}

func (s *AppServer) RegisterRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.HandleHome)
	return mux
}
