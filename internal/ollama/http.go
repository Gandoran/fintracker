package ollama

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strings"
)

func doChatRequest(ctx context.Context, httpClient *http.Client, baseURL string, reqData ChatRequest) (*ChatResponse, error) {
	jsonData, _ := json.Marshal(reqData)
	cleanBase := strings.TrimRight(baseURL, "/")
	endpoint := cleanBase
	if !strings.HasSuffix(cleanBase, "/api/chat") {
		endpoint = cleanBase + "/api/chat"
	}
	req, _ := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")
	resp, err := httpClient.Do(req)
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
