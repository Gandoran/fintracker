package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	fmt.Println("Starting...")
	go startBackgroundDaemon()
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/api/analyze", handleAnalyzeClick)
	port := ":8080"
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Errore on the server: %v", err)
	}
}

func startBackgroundDaemon() {
	ticker := time.NewTicker(30 * time.Minute)
	defer ticker.Stop()
	for {
		fmt.Println("Checking Thing...")
		<-ticker.C
	}
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Hello World</h1>")
}

func handleAnalyzeClick(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<p> Hello Gemma 4 </p>")
}
