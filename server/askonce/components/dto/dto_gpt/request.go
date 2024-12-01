package dto_gpt

type ChatCompletionReq struct {
	Model       string                  `json:"model"`
	Messages    []ChatCompletionMessage `json:"messages"`
	Tools       []Tool                  `json:"tools,omitempty"`
	MaxTokens   int                     `json:"max_tokens,omitempty"`
	Temperature float64                 `json:"temperature,omitempty"`
	TopP        float64                 `json:"top_p,omitempty"`
	Stream      bool                    `json:"stream,omitempty"`
}

type EmbeddingReq struct {
	Model string `json:"model"` // 此处model需要替换、获取对应渠道的model
	Input string `json:"input"`
	User  string `json:"user,omitempty"`
}
