package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"fintracker/internal/config"
	"fintracker/internal/db"
	"fintracker/internal/ollama"
	"fintracker/internal/pipeline"
	"fintracker/internal/scraper"
	"fintracker/internal/web"
)

func main() {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("Configuration Error: %v", err)
	}
	store, err := db.NewStore("fintracker.db")
	if err != nil {
		log.Fatalf("Connection Db Error: %v", err)
	}
	fetcher := scraper.NewFetcher()
	ai := ollama.NewClient(cfg.LLM.URL, cfg.LLM.Model, cfg.LLM.Temperature)
	appServer := web.NewAppServer(store)
	worker := pipeline.NewWorker(cfg, fetcher, ai, store)
	http.HandleFunc("/", appServer.HandleHome)
	srv := &http.Server{Addr: ":8080"}
	go func() {
		log.Println("FinTracker Active on http://localhost:8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server Error: %v", err)
		}
	}()
	go worker.Start()
	//CTRL+C block
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop // if contrl+c is pressed
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Error during Shutdown: %v", err)
	}
}
