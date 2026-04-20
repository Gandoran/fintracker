<p align="center">
  <img src="https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat-square&logo=go&logoColor=white" />
  <img src="https://img.shields.io/badge/SQLite-3-003B57?style=flat-square&logo=sqlite&logoColor=white" />
  <img src="https://img.shields.io/badge/TailwindCSS-3-38B2AC?style=flat-square&logo=tailwindcss&logoColor=white" />
  <img src="https://img.shields.io/badge/Ollama-Local%20LLM-000000?style=flat-square" />
  <img src="https://img.shields.io/badge/Gemma-Google%20AI-4285F4?style=flat-square" />
  <img src="https://img.shields.io/badge/Telegram%20Bot-API-26A5E4?style=flat-square&logo=telegram&logoColor=white" />
</p>

# FinTracker: Automated Financial Intelligence Pipeline

## 📖 Descrizione Dettagliata
FinTracker è una piattaforma di intelligenza finanziaria automatizzata, progettata per risolvere il problema critico dell'eccesso di informazione (information overload) che caratterizza l'ecosistema finanziario globale. Il sistema funge da pipeline end-to-end, ingesting contenuti finanziari da fonti multiple e disordinate (RSS feed, articoli web) e trasformandoli in analisi strutturate, digeribili e immediatamente fruibili.

L'architettura è basata su un flusso di lavoro continuo (Continuous Integration/Continuous Analysis). Il processo inizia con il **Web Scraping** e l'aggiornamento di *feeds* finanziari di alto valore (come quelli di Yahoo Finance o CNBC), salvando i contenuti grezzi come entità `Article`. Successivamente, ogni articolo viene inoltrato al **Modulo di Analisi LLM** (Large Language Model), che utilizza un servizio Ollama per processare il testo. Questo modulo avanzato non si limita a riassumere; esso estrae metriche critiche — quali il *sentimento*, l'*impatto* di mercato, i *ticker* azionari rilevanti e i *link di riferimento* — per popolare l'entità `Analysis`.

L'intero processo è persistito e reso interrogabile tramite un database SQLite, consentendo agli utenti di eseguire query sofisticate non solo su articoli recenti, ma su tendenze analitiche storiche. Infine, il sistema completa il ciclo di valore inviando notifiche in tempo reale via Telegram e esponendo i risultati attraverso un'interfaccia web dedicata. È uno strumento critico per data scientist e analisti che necessitano di trasformare il "rumore" finanziario in segnali di mercato quantificabili.

## 🚀 Funzionalità Principali
**1. Ingestion e Scraping di Feed Finanziari (Scraper):**
Il sistema è in grado di connettersi e scandire periodicamente (configurabile in minuti) specifici RSS feed finanziari. Il `Fetcher` estrae il contenuto, titolando l'articolo e salvandolo come record sorgente. È supportato anche da un modulo di ricerca web potenziato (Tavily) per arricchire o cercare informazioni correlate su query specifiche.

**2. Analisi Semantica e Estrattiva LLM (Ollama Module):**
Questa è la funzione intellettuale centrale. Il `OllamaClient` si interfaccia con un modello LLM locale (specificato da `config.yaml`) per eseguire un'analisi profonda. Questa analisi va oltre il semplice riassunto, strutturando l'output in campi definiti:
*   **Sentiment:** Classificazione del tono (positivo, negativo, neutro).
*   **Impact:** Valutazione dell'impatto potenziale sull'indice o sul settore.
*   **Tickers:** Estrazione automatica di simboli azionari (es. AAPL, MSFT).
*   **Summary/Rationale:** Sintesi concisa del contenuto.
*   Inoltre, il modulo supporta le **Tool Calls**, permettendo al modello di eseguire ricerche web aggiuntive (tramite `WebSearchTool`) per confermare o contestualizzare le informazioni estratte dall'articolo originale.

**3. Gestione dello Stato e Persistenza (Database Store):**
Tutti i dati vengono modellati e persistiti in un database SQLite. Il sistema gestisce automaticamente le migrazioni dello schema (`000001`, `000002`), garantendo l'integrità dei dati tra gli articoli (`Article`) e le loro relative analisi (`Analysis`). Permette la ricerca granulare sia per titolo che per ticker.

**4. Orchestrazione Pipeline (Worker):**
Il `Worker` è il cuore pulsante del sistema. È responsabile del scheduling, del ciclo di vita dell'analisi: preleva gli articoli da processare, li passa all'LLM, attende l'analisi strutturata, salva il risultato nel DB e, infine, avvia il meccanismo di notifica.

**5. Sistema di Notifica (Notifier):**
Implementa un modulo di allerta via Telegram. Quando viene completata un'analisi ritenuta particolarmente rilevante (es. alto impatto o cambiamento di sentiment significativo), l'utente riceve immediatamente un riepilogo conciso del segnale di mercato.

**6. Interfaccia Utente (Web):**
Fornisce un'interfaccia web che espone le funzionalità di ricerca e visualizzazione dei risultati archiviati, permettendo agli utenti di navigare attraverso le analisi più recenti e filtrabili.

## 🛠️ Architettura e Tecnologie
FinTracker adotta un'architettura a strati (Layered Architecture) con un'orchestrazione centralizzata in Go. L'intero progetto è monolitico nell'esecuzione (un unico processo Go), ma logicamente separato in moduli distinti che aderiscono al principio della Separazione delle Preoccupazioni (Separation of Concerns).

**Stack Tecnologico Principale:**
*   **Backend Language:** Go (Golang). Questo fornisce un ambiente performante e sincrono, ideale per gestire pipeline di I/O intensive (scrapping, rete, DB).
*   **Persistenza Dati:** SQLite (tramite la libreria `modernc.org/sqlite`). Utilizzato per la gestione transazionale e l'archiviazione strutturata dei risultati.
*   **Interazione AI:** Ollama (via API REST, simulata nell'interfaccia `Searcher` e `Client`). Il modello `gemma4:e4b` è specificato per l'analisi.
*   **Data Ingestion:** Librerie dedicate per il parsing di RSS (`gofeed`) e web scraping (`goquery`, implicito).
*   **Comunicazione Inter-Componenti:** La comunicazione è prevalentemente sincrona, passando oggetti Go (structs) tra i servizi (es. `Worker` chiama `Fetcher` che chiama `Client`). Il flusso di dati dal "Vuoto" (Feed) al "Stato" (DB) è gestito sequenzialmente.

**Flusso Architetturale dei Dati (Data Flow):**
1. **Configurazione:** `config.yaml` e `.env` inizializzano `Config`.
2. **Scraping:** `Fetcher` pulla URL $\rightarrow$ crea istanze `Article` (Payload Grezzo).
3. **Analisi:** `Worker` $\rightarrow$ passa `Article` $\rightarrow$ `OllamaClient` $\rightarrow$ chiama LLM $\rightarrow$ riceve `Analysis` (Payload Strutturato).
4. **Persistenza:** `Worker` $\rightarrow$ passa `Article` + `Analysis` $\rightarrow$ `Store` (via Transazione DBTX) $\rightarrow$ Database.
5. **Output:** `Worker` $\rightarrow$ `TelegramBot` (Notifica) e `AppServer` (UI).

## 🧩 Moduli e Componenti Core

### `internal/config`
*   **`Config` (Struct):** Modello che incapsula tutte le configurazioni del sistema. È cruciale per l'inizializzazione.
    *   `DotEnvLoad(cfg *Config)`: Gestisce il caricamento delle variabili d'ambiente da file `.env`.
    *   `Load(filepath string)`: Gestisce il caricamento delle configurazioni statiche da file YAML, superando o complementando le variabili d'ambiente.

### `internal/db`
Questo modulo è il gestore di stato persistente e rappresenta la fonte di verità.
*   **`DBTX` (Interface):** L'interfaccia di transazione del database. Garantisce che operazioni multiple (es. inserimento articolo + analisi) avvengano atomicamente.
*   **`Queries` (Struct):** Wrapper su `DBTX`, fornisce metodi tipizzati per interagire con la logica di business del database (es. `CreateArticle`, `SearchAnalyses`).
    *   `WithTx(tx *sql.Tx)`: Implementazione per eseguire operazioni complesse all'interno di una singola transazione.
*   **`Article` (Struct):** Rappresenta un articolo finanziario grezzo (Titolo, Link, Contenuto, Fonte, Data di Pubblicazione). È l'input di base.
*   **`Analysis` (Struct):** Rappresenta l'output dell'analisi LLM. È la chiave di valore del sistema. Contiene metadati strutturati come `sentiment`, `impact`, `tickers`, e un `reliability_score` (aggiunto tramite migrazione).
*   **`Store` (Struct):** L'oggetto singleton che gestisce la connettività fisica al DB.
    *   `NewStore(dbPath string)`: Esegue l'inizializzazione del database e, criticamente, il `runMigrations`.
*   **`migrations`:** Contiene lo schema evolutivo del DB (SQLite).
    *   `000001_init.up.sql`: Crea le tabelle `articles` e `analyses` con il foreign key vincolato.
    *   `000002_init.up.sql`: Aggiunge la colonna `reliability_score` al contesto di analisi, migliorando la validazione dei dati.

### `internal/ollama`
Questo gruppo gestisce l'intelligenza artificiale e le capacità esterne del modello.
*   **`Searcher` (Interface):** Definizione del contratto per qualsiasi componente di ricerca web esterno.
*   **`Client` (Struct):** Il client principale che orchestra la comunicazione con il servizio LLM.
    *   `NewClient(...)`: Configura l'accesso con specifici parametri (URL, Modello, Temperatura).
    *   `doChatRequest()`: Gestisce il payload JSON per la richiesta al modello.
    *   `AnalyzeArticle()`: Il metodo pubblico che implementa il ciclo di analisi, ricevendo l'articolo e restituendo un oggetto `Analysis` strutturato.
*   **`Tool` (Struct):** Definisce un meccanismo di chiamata a strumenti esterni (Tool Calling), essenziale per dotare l'LLM di capacità dinamiche (es. ricerca web).
*   **`WebSearchTool` (Global Variable):** Implementazione specifica di uno strumento di ricerca che viene offerto all'LLM, permettendogli di "verificare" fatti esterni prima di concludere l'analisi.

### `internal/scraper`
*   **`Fetcher` (Struct):** Responsabile della raccolta dati grezzi.
    *   `Fetch(ctx, urls)`: Esegue il fetching HTTP per un batch di URL.
    *   `parseFeed(feed)`: Implementa la logica di parsing specifica per standard RSS/XML.
    *   `sanitizeText(htmlStr)`: Utility per pulire il contenuto HTML prima di salvarlo o analizzarlo.

### `internal/search`
*   **`TavilyClient` (Struct):** Wrapper API specifico per l'integrazione con il motore di ricerca Tavily.
    *   `Search(query)`: Esegue la ricerca web, restituendo non solo testo ma anche una lista di link rilevanti.

### `internal/pipeline`
*   **`Worker` (Struct):** L'Orchestratore di Business.
    *   `NewWorker(...)`: Dipende da tutte le componenti core (Scraper, Ollama, Store, Notifier).
    *   `Start()`: Il ciclo di vita principale, che avvia il processo di alimentazione dei dati.
    *   `processFeeds()`: Coordina il scraping e l'invio degli articoli al processo di analisi.
    *   `processSingleArticle()`: Contiene il flusso critico: Scrape $\rightarrow$ LLM Analyze $\rightarrow$ DB Store.
    *   `SendTelegramNotify()`: Gestisce la decisione di inviare un allerta di sistema.

### `notifier`
*   **`TelegramBot` (Struct):** Componente di output asincrono.
    *   `SendAlert(message string)`: Abilita la comunicazione di stato critico all'utente finale.

### `web`
*   **`AppServer` (Struct):** Gestore del contesto web.
    *   `NewAppServer()`: Inizializza la dipendenza dallo `Store` del DB.
    *   `RegisterRoutes()`: Mappare le richieste HTTP ai handler.
*   **`handlers.go`:** Contiene la logica di accesso ai dati, recuperando i risultati archiviati dallo `Store` e preparandoli per la visualizzazione web.

## 💻 Installazione e Avvio

### Prerequisiti
1. **Go Environment:** Assicurarsi di avere installato Go 1.26.1 o superiore.
2. **Dependency Management:** Le dipendenze principali sono gestite tramite `go.mod`.
3. **Servizi Esterni:**
    *   **Ollama:** Il server LLM deve essere in esecuzione localmente (specificato in `config.yaml`) e il modello `gemma4:e4b` deve essere pre-scaricato.
    *   **API Keys:** È necessario configurare le chiavi API per Tavily e Telegram attraverso variabili d'ambiente (`.env.example`).

### Passo 1: Configurazione dell'Ambiente e Dipendenze
Eseguire il setup delle dipendenze:
```bash
go mod tidy
```

### Esempio App
<img width="1043" height="1026" alt="Screenshot 2026-04-17 235358" src="https://github.com/user-attachments/assets/6d05e446-7c51-404b-8ffa-c6d48842d5a3" />

## Project Structure:
```text
fintracker/
├── cmd
│   └── server
│       └── main.go
├── internal
│   ├── config
│   │   └── config.go
│   ├── db
│   │   ├── migrations
│   │   │   ├── 000001_init.up.sql
│   │   │   └── 000002_init.up.sql
│   │   ├── db.go
│   │   ├── models.go
│   │   ├── queries.sql
│   │   ├── queries.sql.go
│   │   └── store.go
│   ├── models
│   │   ├── analysis.go
│   │   └── article.go
│   ├── ollama
│   │   ├── analyzer.go
│   │   ├── client.go
│   │   ├── systemPrompt.go
│   │   ├── tool.go
│   │   └── types.go
│   ├── pipeline
│   │   └── pipeline.go
│   ├── scraper
│   │   └── fetcher.go
│   └── search
│       └── tavily.go
├── notifier
│   └── telegram.go
├── web
│   ├── handlers.go
│   ├── server.go
│   └── templates.go
├── .env.example
├── config.yaml
├── fintracker.db
├── go.mod
└── sqlc.yaml
<<<<<<< HEAD
```
=======
```
>>>>>>> 4f18e1b (feat: aggiunto modalità di chat per interagire con l'IA sull'articolo selezionato)
