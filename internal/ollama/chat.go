package ollama

import (
	"context"
	"fmt"
)

func (c *ChatClient) ChatWithArticle(ctx context.Context, articleContent string, userQuestion string) (string, error) {
	userPrompt := fmt.Sprintf("ARTICOLO:\n%s\n\nDOMANDA: %s", articleContent, userQuestion)
	req := ChatRequest{
		Model: c.modelName,
		Messages: []Message{
			{Role: "system", Content: SystemPromptChat},
			{Role: "user", Content: userPrompt},
		},
		Stream:  false,
		Options: map[string]interface{}{"temperature": c.temperature},
	}
	resp, err := doChatRequest(ctx, c.httpClient, c.baseURL, req)
	if err != nil {
		return "", err
	}
	return resp.Message.Content, nil
}
