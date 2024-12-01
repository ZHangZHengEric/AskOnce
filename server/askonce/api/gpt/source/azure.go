package source

import (
	"askonce/components/defines"
	"fmt"
)

type AzureGPT struct {
	CommonGPT
}

func (gpt *AzureGPT) GetPath(methodType defines.GPTMethodType, model string) string {
	switch methodType {
	case defines.MethodTypeEmbedding:
		return fmt.Sprintf("/openai/deployments/%s/embeddings?api-version=2023-05-15", model)
	case defines.MethodTypeChat:
		return fmt.Sprintf("/openai/deployments/%s/chat/completions?api-version=2024-02-15-preview", model)
	default:
		return ""
	}
}

func (gpt *AzureGPT) HandleRequestHeader() map[string]string {
	return map[string]string{
		"api-key": fmt.Sprintf("Bearer %s", gpt.accessKey),
	}
}
