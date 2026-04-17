package web

import (
	"html/template"
	"strings"
)

const homeHTML = `
<!DOCTYPE html>
<html lang="it">
<head>
	<meta charset="UTF-8">
	<title>Lumina Dashboard</title>
	<script src="https://cdn.tailwindcss.com"></script>
</head>
<body class="bg-slate-50 text-slate-800 font-sans p-8">
	<div class="max-w-5xl mx-auto">
		<div class="flex justify-between items-center mb-8">
			<h1 class="text-3xl font-bold text-slate-900 tracking-tight">Lumina <span class="text-blue-600">AI</span></h1>
			<form action="/" method="GET" class="flex items-center space-x-2">
				<input type="text" name="q" value="{{.SearchTerm}}" placeholder="Cerca ticker o parola..." 
					class="px-4 py-2 border border-slate-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 w-64">
				<button type="submit" class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg transition">Cerca</button>
				{{if .SearchTerm}}<a href="/" class="text-slate-500 hover:text-slate-800 px-2">X</a>{{end}}
			</form>
		</div>

		<div class="grid grid-cols-1 md:grid-cols-2 gap-6">
			{{range .Results}}
			<div class="bg-white border border-slate-200 rounded-xl p-6 shadow-sm hover:shadow-md transition">
				<h3 class="text-xl font-semibold text-slate-900 mb-2">{{.Article.Title}}</h3>
				<div class="flex items-center space-x-3 mb-4 text-sm">
					<span class="bg-blue-100 text-blue-800 px-2 py-1 rounded font-mono font-medium">{{.Analysis.Tickers}}</span>
					<span class="font-semibold {{if eq .Analysis.Sentiment "Bullish"}}text-green-600{{else if eq .Analysis.Sentiment "Bearish"}}text-red-600{{else}}text-slate-500{{end}}">
						{{.Analysis.Sentiment}}
					</span>
					{{if ge .Analysis.ReliabilityScore 8}}
                        <span class="flex items-center gap-1 bg-emerald-50 text-emerald-700 border border-emerald-200 px-2 py-0.5 rounded-md font-bold text-xs shadow-sm" title="Alta Affidabilità (Fatti ufficiali/Confermati)">
                            🍅 {{.Analysis.ReliabilityScore}}/10
                        </span>
                    {{else if ge .Analysis.ReliabilityScore 5}}
                        <span class="flex items-center gap-1 bg-amber-50 text-amber-700 border border-amber-200 px-2 py-0.5 rounded-md font-bold text-xs shadow-sm" title="Affidabilità Media (Speculazioni basate su dati)">
                            🤔 {{.Analysis.ReliabilityScore}}/10
                        </span>
                    {{else}}
                        <span class="flex items-center gap-1 bg-rose-50 text-rose-700 border border-rose-200 px-2 py-0.5 rounded-md font-bold text-xs shadow-sm" title="Bassa Affidabilità (Rumor/Fonti non verificate)">
                            🤢 {{.Analysis.ReliabilityScore}}/10
                        </span>
                    {{end}}
				</div>
				<p class="text-slate-600 mb-4 text-sm">{{.Analysis.Summary}}</p>
				<div class="bg-slate-50 border-l-4 border-blue-500 p-3 rounded-r-lg mb-4">
					<p class="text-sm text-slate-700"><strong>Impatto:</strong> {{.Analysis.Impact}}</p>
				</div>
				{{if .Analysis.ReferenceLinks}}
				<div class="mt-4 pt-4 border-t border-slate-100">
					<p class="text-xs font-semibold text-slate-400 uppercase tracking-wider mb-2">Fonti Web AI</p>
					<ul class="space-y-1">
						{{range split .Analysis.ReferenceLinks ","}}{{if .}}
						<li><a href="{{.}}" target="_blank" class="text-blue-500 hover:text-blue-700 text-sm truncate block">{{.}}</a></li>
						{{end}}{{end}}
					</ul>
				</div>
				{{end}}
				<div class="mt-4 text-right"><span class="text-xs text-slate-400">Analizzato: {{.Analysis.AnalyzedAt.Time.Format "02/01 15:04"}}</span></div>
			</div>
			{{else}}
				<div class="col-span-full text-center py-12 bg-white border border-slate-200 rounded-xl">
					<p class="text-slate-500 text-lg">Nessuna analisi disponibile. Attendi il demone...</p>
				</div>
			{{end}}
		</div>
	</div>
</body>
</html>`

var tmplFuncs = template.FuncMap{"split": strings.Split}
var homeTmpl = template.Must(template.New("home").Funcs(tmplFuncs).Parse(homeHTML))
