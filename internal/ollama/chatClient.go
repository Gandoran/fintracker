package ollama

import (
	"net/http"
)

type ChatClient struct {
	baseURL     string
	modelName   string
	temperature float32
	httpClient  *http.Client
}

func NewChatClient(url, model string, temp float32) *ChatClient {
	return &ChatClient{
		baseURL:     url,
		modelName:   model,
		temperature: temp,
		httpClient:  &http.Client{},
	}
}
