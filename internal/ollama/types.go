package ollama

type Message struct {
	Role      string     `json:"role"`
	Content   string     `json:"content"`
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
}

type ChatRequest struct {
	Model     string                 `json:"model"`
	Messages  []Message              `json:"messages"`
	Stream    bool                   `json:"stream"`
	Format    string                 `json:"format"`
	KeepAlive int                    `json:"keep_alive"`
	Options   map[string]interface{} `json:"options"`
	Tools     []Tool                 `json:"tools,omitempty"`
}

type ChatResponse struct {
	Message Message `json:"message"`
}

type Tool struct {
	Type     string       `json:"type"`
	Function ToolFunction `json:"function"`
}

type ToolFunction struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Parameters  ToolParameters `json:"parameters"`
}

type ToolParameters struct {
	Type       string                  `json:"type"`
	Properties map[string]ToolProperty `json:"properties"`
	Required   []string                `json:"required"`
}

type ToolProperty struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

type ToolCall struct {
	Function ToolCallFunction `json:"function"`
}

type ToolCallFunction struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}
