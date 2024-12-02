package data

import (
	"askonce/api/gpt"
	"askonce/components/defines"
	"askonce/components/dto/dto_gpt"
	"askonce/conf"
	"github.com/xiangtao94/golib/flow"
)

type GptData struct {
	flow.Data
}

func (d *GptData) ChatSync(modelType string, id string, req *dto_gpt.ChatCompletionReq) (answer string, use dto_gpt.ChatCompletionUsage, err error) {
	return
}

type ChatAnswer struct {
	Answer string `json:"answer"`
	Status string `json:"status"`
}

func (d *GptData) Chat(modelType string, id string, req *dto_gpt.ChatCompletionReq, f func(offset int, chatAnswer ChatAnswer) error) error {
	return nil
}

func (d *GptData) Embedding(text string) (output []float32, err error) {
	embeddingModel := conf.WebConf.EmbeddingModelConf
	g, err := gpt.CreateGptClient(d.GetCtx(), defines.GPTSource(embeddingModel.Source))
	if err != nil {
		return
	}
	g.Init(embeddingModel.Addr, embeddingModel.AK)
	req := &dto_gpt.EmbeddingReq{
		Model: embeddingModel.Model,
		Input: text,
	}

	resp, err := gpt.Embedding(g, req)
	if err != nil {
		return
	}
	if resp.Error != nil {
		return nil, resp.Error
	}
	return resp.Data[0].Embedding, nil
}
