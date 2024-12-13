// Package gpt -----------------------------
// @file      : client.go
// @author    : xiangtao
// @contact   : xiangtao1994@gmail.com
// @time      : 2024/12/13 16:31
// -------------------------------------------

package client

import (
	"askonce/gpt/core"
	"encoding/json"
	"fmt"
	"github.com/xiangtao94/golib/pkg/errors"
	"github.com/xiangtao94/golib/pkg/http"
	"time"
)

type GPTSource string

const (
	GPTSourceOpenAI    GPTSource = "openai"
	GPTSourceKimi      GPTSource = "kimi"
	GPTSourceBaiChuan  GPTSource = "baichuan"
	GPTSourceQwen      GPTSource = "qwen"
	GPTSourceGlm       GPTSource = "glm"
	GPTSourceMiniMax   GPTSource = "minimax"
	GPTSourceSkylark   GPTSource = "skylark" // 字节豆包大模型 v3接口版本
	GPTSourceAzure     GPTSource = "azure"
	GPTSourceDeepInfra GPTSource = "deepinfra"
)

type GPTMethodType string

const (
	MethodTypeChat      GPTMethodType = "chat"
	MethodTypeEmbedding GPTMethodType = "embedding"
)

var SourceModelDefault = map[GPTSource]map[GPTMethodType]string{}

type IClient interface {
	init(endpoint string, ak string)
	GetClient() *http.HttpClientConf
	GetPath(methodType GPTMethodType, model string) string
	HandleRequestHeader() map[string]string
	HandleChatResponse(responseByte []byte) (*core.ChatCompletionResp, error)
	HandleEmbeddingResponse(responseByte []byte) (*core.EmbeddingResp, error)
}

func CreateGptClient(source string, addr, ak string) (IClient, error) {
	var gptClient IClient
	switch GPTSource(source) {
	case GPTSourceOpenAI, GPTSourceKimi, GPTSourceBaiChuan, GPTSourceQwen, GPTSourceDeepInfra:
		gptClient = new(CommonGPT)
	case GPTSourceAzure:
		gptClient = new(AzureGPT)
	case GPTSourceGlm:
		gptClient = new(GlmGPT)
	case GPTSourceMiniMax:
		gptClient = new(MiniMaxGPT)
	case GPTSourceSkylark:
		gptClient = new(SkylarkGPT)
	default:
		return nil, errors.NewError(errors.SYSTEM_ERROR, "GPT来源不支持")
	}
	gptClient.init(addr, ak)
	return gptClient, nil
}

/**
通用开放平台接口调用 路径和鉴权统一
*/

type CommonGPT struct {
	Client    *http.HttpClientConf
	accessKey string
	version   string
}

func (g *CommonGPT) init(endpoint string, ak string) {
	g.accessKey = ak
	g.Client = &http.HttpClientConf{
		Domain:  endpoint,
		Retry:   3,
		Timeout: 60 * time.Second,
	}
}

func (g *CommonGPT) GetClient() *http.HttpClientConf {
	return g.Client
}

func (g *CommonGPT) GetPath(methodType GPTMethodType, model string) string {
	switch methodType {
	case MethodTypeEmbedding:
		return "/v1/embeddings"
	case MethodTypeChat:
		return "/v1/chat/completions"
	default:
		return ""
	}
}

func (g *CommonGPT) HandleRequestHeader() map[string]string {
	return map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", g.accessKey),
	}
}

func (g *CommonGPT) HandleChatResponse(responseByte []byte) (resp *core.ChatCompletionResp, err error) {
	err = json.Unmarshal(responseByte, &resp)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}
	return resp, nil
}

func (g *CommonGPT) HandleEmbeddingResponse(responseByte []byte) (resp *core.EmbeddingResp, err error) {
	err = json.Unmarshal(responseByte, &resp)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}
	return resp, nil
}

type AzureGPT struct {
	CommonGPT
}

func (g *AzureGPT) GetPath(methodType GPTMethodType, model string) string {
	switch methodType {
	case MethodTypeEmbedding:
		return fmt.Sprintf("/openai/deployments/%s/embeddings?api-version=2023-05-15", model)
	case MethodTypeChat:
		return fmt.Sprintf("/openai/deployments/%s/chat/completions?api-version=2024-02-15-preview", model)
	default:
		return ""
	}
}

func (g *AzureGPT) HandleRequestHeader() map[string]string {
	return map[string]string{
		"api-key": fmt.Sprintf("Bearer %s", g.accessKey),
	}
}

type GlmGPT struct {
	CommonGPT
}

func (gpt *GlmGPT) GetPath(methodType GPTMethodType, model string) string {
	switch methodType {
	case MethodTypeEmbedding:
		return "/api/paas/v4/embeddings"
	case MethodTypeChat:
		return "/api/paas/v4/chat/completions"
	default:
		return ""
	}
}

type MiniMaxGPT struct {
	CommonGPT
}

func (gpt *MiniMaxGPT) GetPath(methodType GPTMethodType, model string) string {
	switch methodType {
	case MethodTypeEmbedding:
		return "/v1/embeddings"
	case MethodTypeChat:
		return "/v1/text/chatcompletion_v2"
	default:
		return ""
	}
}

func (gpt *MiniMaxGPT) HandleChatResponse(responseByte []byte) (resp *core.ChatCompletionResp, err error) {
	err = json.Unmarshal(responseByte, &resp)
	if err != nil {
		return nil, err
	}
	if resp.BaseResp != nil && resp.BaseResp.StatusCode != 0 {
		return nil, errors.NewError(resp.BaseResp.StatusCode, resp.BaseResp.StatusMsg)
	}
	return resp, nil
}

func (gpt *MiniMaxGPT) HandleEmbeddingResponse(responseByte []byte) (resp *core.EmbeddingResp, err error) {
	err = json.Unmarshal(responseByte, &resp)
	if err != nil {
		return nil, err
	}
	if resp.BaseResp != nil && resp.BaseResp.StatusCode != 0 {
		return nil, errors.NewError(resp.BaseResp.StatusCode, resp.BaseResp.StatusMsg)
	}
	return resp, nil
}

type SkylarkGPT struct {
	CommonGPT
}

func (gpt *SkylarkGPT) GetPath(methodType GPTMethodType, model string) string {
	switch methodType {
	case MethodTypeEmbedding:
		return "/api/v3/embeddings"
	case MethodTypeChat:
		return "/api/v3/chat/completions"
	default:
		return ""
	}
}
