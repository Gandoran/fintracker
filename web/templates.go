package web

import (
	"html/template"
	"strings"
)

const templatesHTML = `
{{define "base"}}
<!DOCTYPE html>
<html lang="it">
<head>
    <meta charset="UTF-8">
    <title>Lumina Dashboard</title>
    <script src="https://cdn.tailwindcss.com"></script>
</head>
<body class="bg-slate-50 text-slate-800 font-sans p-8 relative">
    
    <div class="max-w-5xl mx-auto">
        <div class="flex justify-between items-center mb-8">
            <h1 class="text-3xl font-bold text-slate-900 tracking-tight">Lumina <span class="text-blue-600">AI</span></h1>
            <form action="/" method="GET" class="flex items-center space-x-2">
                <input type="text" name="q" value="{{.SearchTerm}}" placeholder="Cerca ticker o parola..." 
                    class="px-4 py-2 border border-slate-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 w-64">
                <button type="submit" class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg transition">Cerca</button>
            </form>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            {{range .Results}}
                {{template "card" .}}
            {{else}}
                <div class="col-span-full text-center py-12 bg-white border border-slate-200 rounded-xl">
                    <p class="text-slate-500 text-lg">Nessuna analisi disponibile.</p>
                </div>
            {{end}}
        </div>
    </div>

    {{template "chat_modal"}}
    {{template "scripts"}}
</body>
</html>
{{end}}

{{define "card"}}
<div class="bg-white border border-slate-200 rounded-xl p-6 shadow-sm hover:shadow-md transition flex flex-col justify-between">
    <div>
        <h3 class="text-xl font-semibold text-slate-900 mb-2">{{.Article.Title}}</h3>
        <div class="flex items-center flex-wrap gap-3 mb-4 text-sm">
            <span class="bg-blue-100 text-blue-800 px-2 py-1 rounded font-mono font-medium">{{.Analysis.Tickers}}</span>
            <span class="font-semibold {{if eq .Analysis.Sentiment "Bullish"}}text-green-600{{else if eq .Analysis.Sentiment "Bearish"}}text-red-600{{else}}text-slate-500{{end}}">
                {{.Analysis.Sentiment}}
            </span>
            <span class="flex items-center gap-1 bg-slate-100 text-slate-700 border border-slate-200 px-2 py-0.5 rounded-md font-bold text-xs shadow-sm">
                🎯 {{.Analysis.ReliabilityScore}}/10
            </span>
        </div>
        <p class="text-slate-600 mb-4 text-sm">{{.Analysis.Summary}}</p>
    </div>
    
    <div class="mt-4 pt-4 border-t border-slate-100 text-right flex justify-between items-center">
        <span class="text-xs text-slate-400">Analizzato: {{.Analysis.AnalyzedAt.Time.Format "02/01 15:04"}}</span>
        <button onclick="openChat({{.Article.ID}}, '{{.Article.Title}}')" class="bg-indigo-50 text-indigo-600 hover:bg-indigo-100 border border-indigo-200 px-3 py-1.5 rounded flex items-center gap-2 text-sm font-semibold transition">
            💬 Chiedi a Lumina
        </button>
    </div>
</div>
{{end}}

{{define "chat_modal"}}
<div id="chatOverlay" class="hidden fixed inset-0 bg-slate-900/60 backdrop-blur-sm z-40 transition-opacity"></div>

<div id="chatModal" class="hidden fixed top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 w-[90%] max-w-2xl bg-white rounded-2xl shadow-2xl z-50 overflow-hidden flex flex-col h-[600px]">
    
    <div class="bg-slate-900 text-white p-4 flex justify-between items-center">
        <div>
            <h2 class="text-lg font-bold flex items-center gap-2">Assistente Lumina</h2>
            <p id="chatArticleTitle" class="text-xs text-slate-300 truncate max-w-md mt-1">Titolo Articolo</p>
        </div>
        <button onclick="closeChat()" class="text-slate-300 hover:text-white text-2xl leading-none">&times;</button>
    </div>

    <div id="chatBox" class="flex-1 p-6 overflow-y-auto bg-slate-50 space-y-4">
        <div class="bg-white border border-slate-200 rounded-lg p-3 inline-block max-w-[85%] text-slate-700 shadow-sm text-sm">
            Ciao! Sono pronto a rispondere alle tue domande su questa notizia. Chiedimi pure!
        </div>
    </div>

    <div class="p-4 bg-white border-t border-slate-200">
        <form id="chatForm" onsubmit="sendMessage(event)" class="flex gap-2">
            <input type="hidden" id="currentArticleId" value="">
            <input type="text" id="chatInput" placeholder="Fai una domanda sull'articolo..." required
                class="flex-1 px-4 py-3 border border-slate-300 rounded-xl focus:outline-none focus:ring-2 focus:ring-indigo-500 bg-slate-50">
            <button type="submit" id="sendBtn" class="bg-indigo-600 hover:bg-indigo-700 text-white px-6 py-3 rounded-xl font-semibold shadow transition disabled:opacity-50">
                Invia
            </button>
        </form>
    </div>
</div>
{{end}}

{{define "scripts"}}
<script>
    const overlay = document.getElementById('chatOverlay');
    const modal = document.getElementById('chatModal');
    const chatBox = document.getElementById('chatBox');
    const input = document.getElementById('chatInput');
    const form = document.getElementById('chatForm');
    const sendBtn = document.getElementById('sendBtn');

    function openChat(articleId, articleTitle) {
        document.getElementById('currentArticleId').value = articleId;
        document.getElementById('chatArticleTitle').innerText = articleTitle;
        chatBox.innerHTML = '<div class="bg-white border border-slate-200 rounded-lg p-3 inline-block max-w-[85%] text-slate-700 shadow-sm text-sm">Ciao! Sono pronto a rispondere alle tue domande su questa notizia. Chiedimi pure!</div>';
        overlay.classList.remove('hidden');
        modal.classList.remove('hidden');
        input.focus();
    }

    function closeChat() {
        overlay.classList.add('hidden');
        modal.classList.add('hidden');
    }

    // Chiudi modale cliccando fuori
    overlay.addEventListener('click', closeChat);

    async function sendMessage(e) {
        e.preventDefault();
        const question = input.value.trim();
        const articleId = document.getElementById('currentArticleId').value;
        if (!question) return;

        // 1. Aggiungi messaggio Utente
        appendMessage(question, 'user');
        input.value = '';
        sendBtn.disabled = true;

        // 2. Aggiungi indicatore caricamento IA
        const loadingId = 'loading-' + Date.now();
        appendMessage('Sto analizzando...', 'ai', loadingId);

        try {
            // 3. Chiamata API al backend Go
            const response = await fetch('/api/chat', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    article_id: parseInt(articleId),
                    question: question
                })
            });

            const data = await response.json();
            
            // Rimuovi loading e metti risposta reale
            document.getElementById(loadingId).remove();
            
            if (!response.ok) throw new Error(data.error || 'Errore server');
            
            appendMessage(data.answer, 'ai');

        } catch (error) {
            document.getElementById(loadingId).remove();
            appendMessage('Scusa, c\'è stato un errore di connessione con l\'IA.', 'error');
        } finally {
            sendBtn.disabled = false;
            input.focus();
        }
    }

    function appendMessage(text, sender, id = null) {
        const div = document.createElement('div');
        div.className = sender === 'user' ? 'text-right' : 'text-left';
        if (id) div.id = id;

        const innerClass = sender === 'user' 
            ? 'bg-indigo-600 text-white rounded-l-xl rounded-tr-xl p-3 inline-block max-w-[85%] text-sm shadow-md'
            : (sender === 'error' ? 'bg-red-100 text-red-700 border border-red-200 rounded-r-xl rounded-tl-xl p-3 inline-block max-w-[85%] text-sm' 
            : 'bg-white border border-slate-200 text-slate-700 rounded-r-xl rounded-tl-xl p-3 inline-block max-w-[85%] text-sm shadow-sm whitespace-pre-wrap');

        div.innerHTML = '<div class="' + innerClass + '">' + text + '</div>';
        chatBox.appendChild(div);
        chatBox.scrollTop = chatBox.scrollHeight;
    }
</script>
{{end}}
`

var tmplFuncs = template.FuncMap{"split": strings.Split}
var tmpl = template.Must(template.New("dashboard").Funcs(tmplFuncs).Parse(templatesHTML))
