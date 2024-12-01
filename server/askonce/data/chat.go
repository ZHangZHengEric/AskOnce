package data

import (
	"askonce/components/dto/dto_gpt"
	"github.com/xiangtao94/golib/flow"
)

type ChatData struct {
	flow.Data
}

func (d *ChatData) ChatSync(modelType string, id string, req *dto_gpt.ChatCompletionReq) (answer string, use dto_gpt.ChatCompletionUsage, err error) {
	return
}

func (d *ChatData) Chat(modelType string, id string, req *dto_gpt.ChatCompletionReq, f func(offset int, chatCompletionResp dto_gpt.ChatCompletionResp) error) error {
	return nil
}
