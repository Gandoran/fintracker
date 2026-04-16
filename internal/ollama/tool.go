package ollama

var WebSearchTool = Tool{
	Type: "function",
	Function: ToolFunction{
		Name:        "web_search",
		Description: "Cerca su internet notizie recenti o contesto. Usalo se mancano informazioni.",
		Parameters: ToolParameters{
			Type: "object",
			Properties: map[string]ToolProperty{
				"query": {Type: "string", Description: "Stringa di ricerca"},
			},
			Required: []string{"query"},
		},
	},
}
