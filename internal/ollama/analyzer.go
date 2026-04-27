package ollama

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"fintracker/internal/models"
)

func (c *AnalyzerClient) AnalyzeArticle(ctx context.Context, art models.Article) (*models.Analysis, error) {
	messages := c.buildInitialMessages(art)
	return c.processChatLoop(ctx, messages, art)
}

func (c *AnalyzerClient) buildInitialMessages(art models.Article) []Message {
	return []Message{
		{Role: "system", Content: SystemPromptFinancial},
		{Role: "user", Content: fmt.Sprintf("Titolo: %s\nContenuto: %s", art.Title, art.Content)},
	}
}

func (c *AnalyzerClient) processChatLoop(ctx context.Context, msgs []Message, art models.Article) (*models.Analysis, error) {
	var allFoundLinks []string
	for {
		req := c.buildChatRequest(msgs)
		resp, err := doChatRequest(ctx, c.httpClient, c.baseURL, req)
		if err != nil {
			return nil, err
		}
		if len(resp.Message.ToolCalls) > 0 {
			//tool called by ollama model
			var newLinks []string
			msgs, newLinks = c.handleToolCalls(msgs, resp.Message)
			allFoundLinks = append(allFoundLinks, newLinks...)
			continue
		}
		return c.parseFinalResponse(resp.Message.Content, art, allFoundLinks)
	}
}

func (c *AnalyzerClient) buildChatRequest(messages []Message) ChatRequest {
	req := ChatRequest{
		Model:     c.modelName,
		Messages:  messages,
		Stream:    false,
		Format:    "json",
		KeepAlive: 0,
		Options:   map[string]interface{}{"temperature": c.temperature},
	}
	if c.searcher != nil {
		req.Tools = []Tool{WebSearchTool}
	}
	return req
}

func (c *AnalyzerClient) handleToolCalls(messages []Message, aiMessage Message) ([]Message, []string) {
	query := aiMessage.ToolCalls[0].Function.Arguments["query"].(string)
	fmt.Printf("Searching: '%s'\n", query)
	var searchResults string
	var links []string
	if c.searcher != nil {
		searchResults, links = c.searcher.Search(query)
	} else {
		searchResults = "Web Search disabled or Run out of Token"
	}
	messages = append(messages, aiMessage)
	messages = append(messages, Message{Role: "tool", Content: searchResults})
	return messages, links
}

func (c *AnalyzerClient) parseFinalResponse(content string, art models.Article, links []string) (*models.Analysis, error) {
	var analysis models.Analysis
	if err := json.Unmarshal([]byte(content), &analysis); err != nil {
		return nil, fmt.Errorf("JSON non valido: %v", err)
	}
	analysis.AnalysisAt = time.Now()
	analysis.Original = art
	analysis.ReferenceLinks = links
	return &analysis, nil
}
