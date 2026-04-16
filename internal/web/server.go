package web

import (
	"html/template"
	"net/http"

	"fintracker/internal/db"
)

type AppServer struct {
	store *db.Store
}

func NewAppServer(store *db.Store) *AppServer {
	return &AppServer{
		store: store,
	}
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
				<h3>{{.Article.Title}}</h3>
				<p><strong>Ticker:</strong> {{.Analysis.Tickers}} | <strong>Sentiment:</strong> {{.Analysis.Sentiment}}</p>
				<p><em>Riassunto: {{.Analysis.Summary}}</em></p>
				<p style="color: #d9534f;"><strong>Impatto previsto:</strong> {{.Analysis.Impact}}</p>
				<small>Analizzato il: {{.Analysis.AnalyzedAt.Format "15:04:05"}}</small>
			</div>
		{{else}}
			<p>no analyze available...</p>
		{{end}}
	</body>
	</html>`
	tmpl := template.Must(template.New("home").Parse(htmlPage))
	recentAnalyses, err := s.store.GetRecentAnalyses(r.Context(), 30)
	if err != nil {
		http.Error(w, "Error inside the server", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, recentAnalyses)
}
