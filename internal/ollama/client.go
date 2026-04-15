package ollama

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"fintracker/internal/models"
)

type Client struct {
	baseURL     string
	modelName   string
	temperature float32
	httpClient  *http.Client
}

func NewClient(url, model string, temp float32) *Client {
	return &Client{
		baseURL:     url,
		modelName:   model,
		temperature: temp,
		httpClient:  &http.Client{},
	}
}

func (c *Client) AnalyzeArticle(ctx context.Context, art models.Article) (*models.Analysis, error) {
	prompt := fmt.Sprintf(SystemPromptFinancial, art.Title, art.Content)
	rawJSONResponse, err := c.doOllamaRequest(ctx, prompt)
	if err != nil {
		return nil, err
	}
	var analysis models.Analysis
	if err := json.Unmarshal([]byte(rawJSONResponse), &analysis); err != nil {
		return nil, fmt.Errorf("JSON non valido da Gemma: %v", err)
	}
	analysis.AnalysisAt = time.Now()
	analysis.Original = art
	return &analysis, nil
}

func (c *Client) doOllamaRequest(ctx context.Context, prompt string) (string, error) {
	reqData := GenerateRequest{
		Model:   c.modelName,
		Prompt:  prompt,
		Stream:  false,
		Format:  "json",
		Options: map[string]interface{}{"temperature": c.temperature},
	}
	jsonData, _ := json.Marshal(reqData)
	req, _ := http.NewRequestWithContext(ctx, "POST", c.baseURL, bytes.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("errore connessione Ollama: %v", err)
	}
	defer resp.Body.Close()
	bodyBytes, _ := io.ReadAll(resp.Body)

	var ollamaResp GenerateResponse
	if err := json.Unmarshal(bodyBytes, &ollamaResp); err != nil {
		return "", fmt.Errorf("impossibile leggere la risposta di rete: %v", err)
	}
	return ollamaResp.Response, nil
}
