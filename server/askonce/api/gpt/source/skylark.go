package source

import "askonce/components/defines"

type SkylarkGPT struct {
	CommonGPT
}

func (gpt *SkylarkGPT) GetPath(methodType defines.GPTMethodType, model string) string {
	switch methodType {
	case defines.MethodTypeEmbedding:
		return "/api/v3/embeddings"
	case defines.MethodTypeChat:
		return "/api/v3/chat/completions"
	default:
		return ""
	}
}
