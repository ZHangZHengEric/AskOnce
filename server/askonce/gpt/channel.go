// Package gpt -----------------------------
// @file      : channel.go
// @author    : xiangtao
// @contact   : xiangtao1994@gmail.com
// @time      : 2024/12/13 16:31
// -------------------------------------------
package gpt

import (
	"askonce/gpt/client"
	"askonce/gpt/core"
	"github.com/gin-gonic/gin"
	http2 "github.com/xiangtao94/golib/pkg/http"
	"net/http"
)

type ChannelConf struct {
	Source string `yaml:"source"` // 渠道
	Addr   string `yaml:"addr"`
	AK     string `yaml:"ak"`
	Model  string `yaml:"model"`
}

type Channel struct {
	ctx       *gin.Context
	gptClient client.IClient
}

func CreatChannel(ctx *gin.Context, conf ChannelConf) (*Channel, error) {
	gptClient, err := client.CreateGptClient(conf.Source, conf.Addr, conf.AK)
	if err != nil {
		return nil, err
	}
	return &Channel{ctx: ctx, gptClient: gptClient}, nil
}

func (c *Channel) Embedding(req *core.EmbeddingReq) (resp *core.EmbeddingResp, err error) {
	path := c.gptClient.GetPath(client.MethodTypeEmbedding, "")
	if len(path) == 0 {
		return nil, core.NewOpenAiError("method_no_support", "方法不支持")
	}
	header := c.gptClient.HandleRequestHeader()
	requestOpt := http2.HttpRequestOptions{
		RequestBody: req,
		Encode:      http2.EncodeJson,
		Headers:     header,
	}
	httpResult, err := c.gptClient.GetClient().HttpPost(c.ctx, path, requestOpt)
	if err != nil {
		return nil, core.NewOpenAiError("do_request_failed", err.Error())
	}
	if httpResult.HttpCode != http.StatusOK {
		return nil, core.NewOpenAiError("do_request_failed", string(httpResult.Response))
	}
	resp, err = c.gptClient.HandleEmbeddingResponse(httpResult.Response)
	if err != nil {
		return nil, core.NewOpenAiError("handle_response_failed", err.Error())
	}
	if resp.Error != nil {
		return nil, resp.Error
	}
	return
}

func (c *Channel) ChatCompletion(req *core.ChatCompletionReq) (resp *core.ChatCompletionResp, err error) {
	path := c.gptClient.GetPath(client.MethodTypeChat, req.Model)
	if len(path) == 0 {
		return nil, core.NewOpenAiError("method_no_support", "方法不支持")
	}
	header := c.gptClient.HandleRequestHeader()
	requestOpt := http2.HttpRequestOptions{
		RequestBody: req,
		Encode:      http2.EncodeJson,
		Headers:     header,
	}
	if req.Stream == true {
		err := c.gptClient.GetClient().HttpPostStream(c.ctx, path, requestOpt, func(data string) error {
			if len(data) < 6 { // ignore blank line or wrong format
				return nil
			}
			if data[:6] != "data: " && data[:6] != "[DONE]" {
				return nil
			}
			if data[:6] == "data: " {
				tmpData := data[6:]
				resp, err = c.gptClient.HandleChatResponse([]byte(tmpData))
				if err != nil {
					return err
				}
				if resp.Error != nil {
					return resp.Error
				}
			}
			return nil
		})
		if err != nil {
			return nil, core.NewOpenAiError("do_request_failed", err.Error())
		}
		return resp, nil
	}
	httpResult, err := c.gptClient.GetClient().HttpPost(c.ctx, path, requestOpt)
	if err != nil {
		return nil, core.NewOpenAiError("do_request_failed", err.Error())
	}
	if httpResult.HttpCode != http.StatusOK {
		return nil, core.NewOpenAiError("do_request_failed", string(httpResult.Response))
	}
	resp, err = c.gptClient.HandleChatResponse(httpResult.Response)
	if err != nil {
		return nil, core.NewOpenAiError("handle_response_failed", err.Error())
	}
	if resp.Error != nil {
		return nil, resp.Error
	}
	return
}
