package web

import (
	"html/template"
	"net/http"
	"sync"

	"fintracker/internal/models"
)

type AppServer struct {
	results []models.Analysis
	mu      sync.RWMutex
}

func NewAppServer() *AppServer {
	return &AppServer{
		results: make([]models.Analysis, 0),
	}
}

func (s *AppServer) AddResult(a models.Analysis) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.results = append([]models.Analysis{a}, s.results...)
}

func (s *AppServer) HandleHome(w http.ResponseWriter, r *http.Request) {
	const htmlPage = `
	<html>
	<head><title>FinTracker Dashboard</title></head>
	<body style="font-family: sans-serif; padding: 20px;">
		<h1>FinTracker: Analyze...</h1>
		<hr>
		{{range .}}
			<div style="border: 1px solid #ccc; padding: 15px; margin-bottom: 10px; border-radius: 8px;">
				<h3>{{.Original.Title}}</h3>
				<p><strong>Ticker:</strong> {{range .Ticker}}{{.}} {{end}} | <strong>Sentiment:</strong> {{.Sentiment}}</p>
				<p><em>Riassunto: {{.Summary}}</em></p>
				<p style="color: #d9534f;"><strong>Impatto previsto:</strong> {{.Impact}}</p>
				<small>Analizzato il: {{.AnalysisAt.Format "15:04:05"}}</small>
			</div>
		{{else}}
			<p>no analyze available...</p>
		{{end}}
	</body>
	</html>`
	tmpl := template.Must(template.New("home").Parse(htmlPage))
	s.mu.RLock() //one article write per time
	defer s.mu.RUnlock()
	tmpl.Execute(w, s.results)
}
