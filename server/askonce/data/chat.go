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

type ChatAnswer struct {
	Answer string `json:"answer"`
	Status string `json:"status"`
}

func (d *ChatData) Chat(modelType string, id string, req *dto_gpt.ChatCompletionReq, f func(offset int, chatAnswer ChatAnswer) error) error {
	return nil
}
