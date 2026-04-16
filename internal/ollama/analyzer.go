package ollama

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"fintracker/internal/models"
)

func (c *Client) AnalyzeArticle(ctx context.Context, art models.Article) (*models.Analysis, error) {
	messages := c.buildInitialMessages(art)
	return c.processChatLoop(ctx, messages, art)
}

func (c *Client) buildInitialMessages(art models.Article) []Message {
	return []Message{
		{Role: "system", Content: SystemPromptFinancial},
		{Role: "user", Content: fmt.Sprintf("Titolo: %s\nContenuto: %s", art.Title, art.Content)},
	}
}

func (c *Client) processChatLoop(ctx context.Context, msgs []Message, art models.Article) (*models.Analysis, error) {
	for {
		req := c.buildChatRequest(msgs)
		resp, err := c.doChatRequest(ctx, req)
		if err != nil {
			return nil, err
		}
		if len(resp.Message.ToolCalls) > 0 {
			msgs = c.handleToolCalls(msgs, resp.Message)
			continue
		}
		return c.parseFinalResponse(resp.Message.Content, art)
	}
}

func (c *Client) buildChatRequest(messages []Message) ChatRequest {
	return ChatRequest{
		Model:     c.modelName,
		Messages:  messages,
		Stream:    false,
		Format:    "json",
		KeepAlive: 0,
		Options:   map[string]interface{}{"temperature": c.temperature},
		Tools:     []Tool{WebSearchTool},
	}
}

func (c *Client) handleToolCalls(messages []Message, aiMessage Message) []Message {
	query := aiMessage.ToolCalls[0].Function.Arguments["query"].(string)
	fmt.Printf("Searching: '%s'\n", query)

	//PLACEHOLDER
	mockResult := "Risultati finti per: " + query

	messages = append(messages, aiMessage)
	messages = append(messages, Message{Role: "tool", Content: mockResult})
	return messages
}

func (c *Client) parseFinalResponse(content string, art models.Article) (*models.Analysis, error) {
	var analysis models.Analysis
	if err := json.Unmarshal([]byte(content), &analysis); err != nil {
		return nil, fmt.Errorf("JSON non valido: %v", err)
	}
	analysis.AnalysisAt = time.Now()
	analysis.Original = art
	return &analysis, nil
}
