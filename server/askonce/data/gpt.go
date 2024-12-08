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
	gptClient gpt.IGPT
}

func (d *GptData) OnCreate() {
	embeddingModel := conf.WebConf.EmbeddingModelConf
	gptClient, err := gpt.CreateGptClient(d.GetCtx(), defines.GPTSource(embeddingModel.Source))
	if err != nil {
		return
	}
	gptClient.Init(embeddingModel.Addr, embeddingModel.AK)
	d.gptClient = gptClient
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

func (d *GptData) Embedding(texts []string) (output [][]float32, err error) {
	embeddingModel := conf.WebConf.EmbeddingModelConf
	req := &dto_gpt.EmbeddingReq{
		Model: embeddingModel.Model,
		Input: texts,
	}

	resp, err := gpt.Embedding(d.gptClient, req)
	if err != nil {
		return
	}
	if resp.Error != nil {
		return nil, resp.Error
	}
	output = make([][]float32, 0)
	for _, bb := range resp.Data {
		output = append(output, bb.Embedding)
	}
	return output, nil
}
