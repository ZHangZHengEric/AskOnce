package gpt

import (
	"askonce/components"
	"askonce/components/defines"
	"askonce/components/dto/dto_gpt"
	http2 "github.com/xiangtao94/golib/pkg/http"
	"net/http"
)

func Embedding(gpt IGPT, req *dto_gpt.EmbeddingReq) (resp *dto_gpt.EmbeddingResp, err error) {
	path := gpt.GetPath(defines.MethodTypeEmbedding, "")
	if len(path) == 0 {
		return nil, components.NewOpenAiError("method_no_support", "方法不支持")
	}
	header := gpt.HandleRequestHeader()
	requestOpt := http2.HttpRequestOptions{
		RequestBody: req,
		Encode:      http2.EncodeJson,
		Headers:     header,
	}
	httpResult, err := gpt.GetClient().HttpPost(gpt.GetCtx(), path, requestOpt)
	if err != nil {
		return nil, components.NewOpenAiError("do_request_failed", err.Error())
	}
	if httpResult.HttpCode != http.StatusOK {
		return nil, components.NewOpenAiError("do_request_failed", string(httpResult.Response))
	}
	resp, err = gpt.HandleEmbeddingResponse(httpResult.Response)
	if err != nil {
		return nil, components.NewOpenAiError("handle_response_failed", err.Error())
	}
	return
}
