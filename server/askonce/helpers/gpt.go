// Package helpers -----------------------------
// @file      : gpt.go
// @author    : xiangtao
// @contact   : xiangtao1994@gmail.com
// @time      : 2024/12/18 10:17
// -------------------------------------------
package helpers

import (
	"askonce/conf"
	"github.com/tmc/langchaingo/llms/openai"
)

var EmbeddingGpt *openai.LLM

func InitGpt() {
	var err error
	gptConf := conf.WebConf.Gpt
	for k, v := range gptConf {
		if k == "embedding" {
			EmbeddingGpt, err = openai.New(
				openai.WithBaseURL(v.Addr),
				openai.WithToken(v.AK),
				openai.WithEmbeddingModel(v.Model),
			)
			if err != nil {
				panic("init gpt err: " + err.Error())
			}
		}
	}
}
