package source

import (
	"askonce/components/defines"
	"askonce/components/dto/dto_gpt"
	"encoding/json"
	"github.com/xiangtao94/golib/pkg/errors"
)

type MiniMaxGPT struct {
	CommonGPT
}

func (gpt *MiniMaxGPT) GetPath(methodType defines.GPTMethodType, model string) string {
	switch methodType {
	case defines.MethodTypeEmbedding:
		return "/v1/embeddings"
	case defines.MethodTypeChat:
		return "/v1/text/chatcompletion_v2"
	default:
		return ""
	}
}

func (gpt *MiniMaxGPT) HandleChatResponse(responseByte []byte) (resp *dto_gpt.ChatCompletionResp, err error) {
	err = json.Unmarshal(responseByte, &resp)
	if err != nil {
		return nil, err
	}
	if resp.BaseResp != nil && resp.BaseResp.StatusCode != 0 {
		return nil, errors.NewError(resp.BaseResp.StatusCode, resp.BaseResp.StatusMsg)
	}
	return resp, nil
}

func (gpt *MiniMaxGPT) HandleEmbeddingResponse(responseByte []byte) (resp *dto_gpt.EmbeddingResp, err error) {
	err = json.Unmarshal(responseByte, &resp)
	if err != nil {
		return nil, err
	}
	if resp.BaseResp != nil && resp.BaseResp.StatusCode != 0 {
		return nil, errors.NewError(resp.BaseResp.StatusCode, resp.BaseResp.StatusMsg)
	}
	return resp, nil
}
