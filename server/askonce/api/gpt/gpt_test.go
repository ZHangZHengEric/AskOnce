package gpt

import (
	"askonce/components/defines"
	"askonce/components/dto/dto_gpt"
	"askonce/test"
	"encoding/json"
	"fmt"

	"testing"
)

func TestCommonGPT_ChatCompletion(t *testing.T) {
	test.Init()
	g, errB := CreateGptClient(test.Ctx, defines.GPTSourceSkylark)
	if errB != nil {
		return
	}
	g.Init("", "")
	req := &dto_gpt.ChatCompletionReq{
		Model: "gpt-4-vision-preview",
		Messages: []dto_gpt.ChatCompletionMessage{
			{
				Role:    dto_gpt.ChatCompletionUser,
				Content: "",
			},
		},
	}

	resp, errA := ChatCompletion(g, req)
	if errA != nil {
		fmt.Println(errA)
		return
	}
	respStr, _ := json.Marshal(resp)

	fmt.Println(string(respStr))
}

func TestEmbedding(t *testing.T) {
	test.Init()
	g, errB := CreateGptClient(test.Ctx, defines.GPTSourceSkylark)
	if errB != nil {
		return
	}
	g.Init("", "")
	req := &dto_gpt.EmbeddingReq{
		Model: "ep-20241101093918-nwnsj",
		Input: "今天天气真好",
	}

	resp, errA := Embedding(g, req)
	if errA != nil {
		fmt.Println(errA)
		return
	}
	respStr, _ := json.Marshal(resp)

	fmt.Println(string(respStr))
}
