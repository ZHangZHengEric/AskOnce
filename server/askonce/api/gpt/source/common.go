package source

import (
	"askonce/components/defines"
	"askonce/components/dto/dto_gpt"
	"encoding/json"
	"fmt"
	"github.com/xiangtao94/golib/flow"
	"github.com/xiangtao94/golib/pkg/http"
)

/**
通用开放平台接口调用 路径和鉴权统一
*/

type CommonGPT struct {
	flow.Api
	accessKey string
	version   string
}

func (gpt *CommonGPT) Init(endpoint string, ak string, proxyUrl string) {
	gpt.accessKey = ak
	gpt.Client = &http.HttpClientConf{
		Domain: endpoint,
		Retry:  3,
		Proxy:  proxyUrl,
	}
}

func (gpt *CommonGPT) GetClient() *http.HttpClientConf {
	return gpt.Client
}

func (gpt *CommonGPT) GetPath(methodType defines.GPTMethodType, model string) string {
	switch methodType {
	case defines.MethodTypeEmbedding:
		return "/v1/embeddings"
	case defines.MethodTypeChat:
		return "/v1/chat/completions"
	default:
		return ""
	}
}

func (gpt *CommonGPT) HandleRequestHeader() map[string]string {
	return map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", gpt.accessKey),
	}
}

func (gpt *CommonGPT) HandleChatResponse(responseByte []byte) (resp *dto_gpt.ChatCompletionResp, err error) {
	err = json.Unmarshal(responseByte, &resp)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}
	return resp, nil
}

func (gpt *CommonGPT) HandleEmbeddingResponse(responseByte []byte) (resp *dto_gpt.EmbeddingResp, err error) {
	err = json.Unmarshal(responseByte, &resp)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}
	return resp, nil
}
