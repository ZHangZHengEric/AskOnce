package dto_config

import "askOnce/components/defines"

type DetailReq struct {
}

type SaveReq struct {
	Language  string                      `json:"language"`
	ModelType defines.ChatCompletionModel `json:"modelType"`
}
