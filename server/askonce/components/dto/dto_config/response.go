package dto_config

import "askOnce/components/defines"

type ConfigResp struct {
	Language  string                      `json:"language"`
	ModelType defines.ChatCompletionModel `json:"modelType"`
}

type Dict struct {
	Name   string `json:"name"`
	EnName string `json:"enName"`
	Value  string `json:"value"`
}
