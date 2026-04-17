package search

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type TavilyClient struct {
	apiKey string
}

func NewTavilyClient(apiKey string) *TavilyClient {
	return &TavilyClient{apiKey: apiKey}
}

type tavilyRequest struct {
	APIKey        string `json:"api_key"`
	Query         string `json:"query"`
	IncludeAnswer bool   `json:"include_answer"`
	MaxResults    int    `json:"max_results"`
}

func (t *TavilyClient) Search(query string) (string, []string) {
	reqBody := tavilyRequest{
		APIKey:        t.apiKey,
		Query:         query,
		IncludeAnswer: false,
		MaxResults:    3,
	}
	jsonData, _ := json.Marshal(reqBody)
	resp, err := http.Post("https://api.tavily.com/search", "application/json", bytes.NewBuffer(jsonData))
	if err != nil || resp.StatusCode != 200 {
		return "Errore connection to Tavily.", nil
	}
	defer resp.Body.Close()
	var result struct {
		Results []struct {
			Title   string `json:"title"`
			Url     string `json:"url"`
			Content string `json:"content"`
		} `json:"results"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "Errore reading tavily.", nil
	}
	var builder strings.Builder
	var links []string
	for i, r := range result.Results {
		builder.WriteString(fmt.Sprintf("[%d] %s: %s\n", i+1, r.Title, r.Content))
		links = append(links, r.Url)
	}
	if builder.Len() == 0 {
		return "No information founded.", nil
	}
	return builder.String(), links
}
