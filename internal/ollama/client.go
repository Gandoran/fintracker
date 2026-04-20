package ollama

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"sync"
)

type Searcher interface {
	Search(query string) (string, []string)
}

type Client struct {
	baseURL     string
	modelName   string
	temperature float32
	httpClient  *http.Client
	searcher    Searcher
	mu          sync.Mutex
}

func NewClient(url, model string, temp float32, searcher Searcher) *Client {
	return &Client{
		baseURL:     url,
		modelName:   model,
		temperature: temp,
		httpClient:  &http.Client{},
		searcher:    searcher,
	}
}

func (c *Client) doChatRequest(ctx context.Context, reqData ChatRequest) (*ChatResponse, error) {
	jsonData, _ := json.Marshal(reqData)
	req, _ := http.NewRequestWithContext(ctx, "POST", c.baseURL, bytes.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var chatResp ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return nil, err
	}
	return &chatResp, nil
}
