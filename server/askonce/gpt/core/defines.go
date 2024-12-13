// Package gpt -----------------------------
// @file      : defines.go
// @author    : xiangtao
// @contact   : xiangtao1994@gmail.com
// @time      : 2024/12/13 17:02
// -------------------------------------------
package core

type ChatCompletionReq struct {
	Model       string                  `json:"model"`
	Messages    []ChatCompletionMessage `json:"messages"`
	Tools       []Tool                  `json:"tools,omitempty"`
	MaxTokens   int                     `json:"max_tokens,omitempty"`
	Temperature float64                 `json:"temperature,omitempty"`
	TopP        float64                 `json:"top_p,omitempty"`
	Stream      bool                    `json:"stream,omitempty"`
}

type ChatCompletionResp struct {
	Id                string                 `json:"id,omitempty"`
	Object            string                 `json:"object,omitempty"`
	Created           int                    `json:"created,omitempty"`
	Model             string                 `json:"model,omitempty"`
	Choices           []ChatCompletionChoice `json:"choices,omitempty"` // 聊天完成选项的列表
	Usage             *ChatCompletionUsage   `json:"usage,omitempty"`
	SystemFingerprint string                 `json:"system_fingerprint,omitempty"`
	Error             *OpenAIError           `json:"error,omitempty"`
	BaseResp          *struct {
		StatusCode int    `json:"status_code"`
		StatusMsg  string `json:"status_msg"`
	} `json:"base_resp,omitempty"` // minimax 错误
}

type ChatCompletionChoice struct {
	Index        int                    `json:"index"`
	FinishReason string                 `json:"finish_reason"`
	Message      *ChatCompletionMessage `json:"message,omitempty"`
	Delta        *ChatCompletionMessage `json:"delta,omitempty"`
	Usage        *ChatCompletionUsage   `json:"usage,omitempty"`
}

type ChatCompletionRole string

const (
	ChatCompletionUser      ChatCompletionRole = "user"
	ChatCompletionSystem    ChatCompletionRole = "system"
	ChatCompletionAssistant ChatCompletionRole = "assistant"
)

type ChatCompletionUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type ChatCompletionMessage struct {
	Role      ChatCompletionRole `json:"role"`
	Content   string             `json:"content"`
	ToolCalls *[]ToolCall        `json:"tool_calls,omitempty"`
}

type ToolCall struct {
	Index    int          `json:"index"`
	Id       string       `json:"id"`
	Type     string       `json:"type"` //function
	Function FunctionCall `json:"function"`
}

type FunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type ToolType string

const (
	ToolTypeFunction ToolType = "function"
)

type Tool struct {
	Type     ToolType            `json:"type"`
	Function *FunctionDefinition `json:"function,omitempty"`
}

type FunctionDefinition struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Strict      bool   `json:"strict,omitempty"`
	// Parameters is an object describing the function.
	// You can pass json.RawMessage to describe the schema,
	// or you can pass in a struct which serializes to the proper JSON schema.
	// The jsonschema package is provided for convenience, but you should
	// consider another specialized library if you require more complex schemas.
	Parameters any `json:"parameters"`
}

type EmbeddingReq struct {
	Model string   `json:"model"`
	Input []string `json:"input"`
	User  string   `json:"user,omitempty"`
}

type EmbeddingResp struct {
	Object   string          `json:"object"`
	Data     []EmbeddingData `json:"data"`
	Model    string          `json:"model"`
	Usage    EmbeddingUsage  `json:"usage"`
	Error    *OpenAIError    `json:"error,omitempty"`
	BaseResp *struct {
		StatusCode int    `json:"status_code"`
		StatusMsg  string `json:"status_msg"`
	} `json:"base_resp,omitempty"` // minimax 错误
}

type EmbeddingData struct {
	Object    string    `json:"object"`
	Embedding []float32 `json:"embedding"`
	Index     int       `json:"index"`
}

type EmbeddingUsage struct {
	PromptTokens int `json:"prompt_tokens"`
	TotalTokens  int `json:"total_tokens"`
}
