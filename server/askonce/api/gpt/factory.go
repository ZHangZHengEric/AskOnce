package gpt

import (
	"askonce/api/gpt/source"
	"askonce/components/defines"
	"askonce/components/dto/dto_gpt"
	"github.com/gin-gonic/gin"
	"github.com/xiangtao94/golib/flow"
	"github.com/xiangtao94/golib/pkg/errors"
	"github.com/xiangtao94/golib/pkg/http"
)

type IGPT interface {
	flow.IApi
	Init(endpoint string, ak string)
	GetClient() *http.HttpClientConf
	GetPath(methodType defines.GPTMethodType, model string) string
	HandleRequestHeader() map[string]string
	HandleChatResponse(responseByte []byte) (*dto_gpt.ChatCompletionResp, error)
	HandleEmbeddingResponse(responseByte []byte) (*dto_gpt.EmbeddingResp, error)
}

func CreateGptClient(ctx *gin.Context, GPTSource defines.GPTSource) (IGPT, error) {
	switch GPTSource {
	case defines.GPTSourceOpenAI, defines.GPTSourceKimi, defines.GPTSourceBaiChuan, defines.GPTSourceQwen:
		return flow.Create(ctx, new(source.CommonGPT)), nil
	case defines.GPTSourceAzure:
		return flow.Create(ctx, new(source.AzureGPT)), nil
	case defines.GPTSourceGlm:
		return flow.Create(ctx, new(source.GlmGPT)), nil
	case defines.GPTSourceMiniMax:
		return flow.Create(ctx, new(source.MiniMaxGPT)), nil
	case defines.GPTSourceSkylark:
		return flow.Create(ctx, new(source.SkylarkGPT)), nil
	default:
		return nil, errors.NewError(errors.SYSTEM_ERROR, "GPT来源不支持")
	}
}
