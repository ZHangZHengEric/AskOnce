package gpt

import (
	"askonce/components"
	"askonce/components/defines"
	"askonce/components/dto/dto_gpt"
	http2 "github.com/xiangtao94/golib/pkg/http"
	"net/http"
)

func ChatCompletion(gpt IGPT, req *dto_gpt.ChatCompletionReq) (resp *dto_gpt.ChatCompletionResp, err error) {
	path := gpt.GetPath(defines.MethodTypeChat, req.Model)
	if len(path) == 0 {
		return nil, components.NewOpenAiError("method_no_support", "方法不支持")
	}
	header := gpt.HandleRequestHeader()
	requestOpt := http2.HttpRequestOptions{
		RequestBody: req,
		Encode:      http2.EncodeJson,
		Headers:     header,
	}
	if req.Stream == true {
		err := gpt.GetClient().HttpPostStream(gpt.GetCtx(), path, requestOpt, func(data string) error {
			if len(data) < 6 { // ignore blank line or wrong format
				return nil
			}
			if data[:6] != "data: " && data[:6] != "[DONE]" {
				return nil
			}
			if data[:6] == "data: " {
				tmpData := data[6:]
				resp, err = gpt.HandleChatResponse([]byte(tmpData))
				if err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			return nil, components.NewOpenAiError("do_request_failed", err.Error())
		}
	} else {
		httpResult, err := gpt.GetClient().HttpPost(gpt.GetCtx(), path, requestOpt)
		if err != nil {
			return nil, components.NewOpenAiError("do_request_failed", err.Error())
		}
		if httpResult.HttpCode != http.StatusOK {
			return nil, components.NewOpenAiError("do_request_failed", string(httpResult.Response))
		}
		resp, err = gpt.HandleChatResponse(httpResult.Response)
		if err != nil {
			return nil, components.NewOpenAiError("handle_response_failed", err.Error())
		}

	}
	return
}
