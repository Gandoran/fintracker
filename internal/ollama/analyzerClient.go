package ollama

import (
	"net/http"
	"sync"
)

type Searcher interface {
	Search(query string) (string, []string)
}

type AnalyzerClient struct {
	baseURL     string
	modelName   string
	temperature float32
	httpClient  *http.Client
	searcher    Searcher
	mu          sync.Mutex
}

func NewAnalyzerClient(url, model string, temp float32, searcher Searcher) *AnalyzerClient {
	return &AnalyzerClient{
		baseURL:     url,
		modelName:   model,
		temperature: temp,
		httpClient:  &http.Client{},
		searcher:    searcher,
	}
}
