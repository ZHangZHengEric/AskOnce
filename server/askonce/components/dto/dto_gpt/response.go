package dto_gpt

import (
	"askonce/components"
)

type ChatCompletionResp struct {
	Id                string                  `json:"id,omitempty"`
	Object            string                  `json:"object,omitempty"`
	Created           int                     `json:"created,omitempty"`
	Model             string                  `json:"model,omitempty"`
	Choices           []ChatCompletionChoice  `json:"choices,omitempty"` // 聊天完成选项的列表
	Usage             *ChatCompletionUsage    `json:"usage,omitempty"`
	SystemFingerprint string                  `json:"system_fingerprint,omitempty"`
	Error             *components.OpenAIError `json:"error,omitempty"`
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

type EmbeddingResp struct {
	Object   string                  `json:"object"`
	Data     []EmbeddingData         `json:"data"`
	Model    string                  `json:"model"`
	Usage    EmbeddingUsage          `json:"usage"`
	Error    *components.OpenAIError `json:"error,omitempty"`
	BaseResp *struct {
		StatusCode int    `json:"status_code"`
		StatusMsg  string `json:"status_msg"`
	} `json:"base_resp,omitempty"` // minimax 错误
}

type EmbeddingData struct {
	Object    string    `json:"object"`
	Embedding []float64 `json:"embedding"`
	Index     int       `json:"index"`
}
