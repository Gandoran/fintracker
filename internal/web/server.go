package web

import (
	"html/template"
	"net/http"
	"strings"

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
	<body style="font-family: sans-serif; padding: 20px; background-color: #f4f4f9;">
		<h1 style="color: #333;">FinTracker: Analisi Finanziaria</h1>
		<hr>
		{{range .}}
			<div style="background: white; border: 1px solid #ccc; padding: 20px; margin-bottom: 15px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.05);">
				<h3 style="margin-top: 0; color: #0056b3;">{{.Article.Title}}</h3>
				<p><strong>Ticker:</strong> <span style="background: #eef; padding: 3px 8px; border-radius: 4px;">{{.Analysis.Tickers}}</span> | <strong>Sentiment:</strong> {{.Analysis.Sentiment}}</p>
				<p style="font-size: 1.05em; line-height: 1.4;"><em>Riassunto: {{.Analysis.Summary}}</em></p>
				<p style="color: #d9534f; border-left: 3px solid #d9534f; padding-left: 10px;"><strong>Impatto previsto:</strong> {{.Analysis.Impact}}</p>
				{{if .Analysis.ReferenceLinks}}
					<div style="margin-top: 15px; padding-top: 10px; border-top: 1px dashed #eee;">
						<strong style="font-size: 0.9em; color: #555;">🔗 Fonti di verifica (Tavli Search):</strong>
						<ul style="margin: 5px 0; padding-left: 20px; font-size: 0.85em;">
						{{range split .Analysis.ReferenceLinks ","}}
							{{if .}} <li><a href="{{.}}" target="_blank" style="color: #0066cc; text-decoration: none;">{{.}}</a></li>
							{{end}}
						{{end}}
						</ul>
					</div>
				{{end}}

				<div style="margin-top: 15px; text-align: right;">
					<small style="color: #888;">Analizzato il: {{.Analysis.AnalyzedAt.Format "02/01/2006 15:04:05"}}</small>
				</div>
			</div>
		{{else}}
			<p>Nessuna analisi disponibile. Il demone è in attesa di nuove notizie...</p>
		{{end}}
	</body>
	</html>`
	funcMap := template.FuncMap{
		"split": strings.Split,
	}
	tmpl := template.Must(template.New("home").Funcs(funcMap).Parse(htmlPage))
	recentAnalyses, err := s.store.GetRecentAnalyses(r.Context(), 30)
	if err != nil {
		http.Error(w, "Errore interno del server DB", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, recentAnalyses)
}
