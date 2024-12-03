package helpers

import (
	"encoding/json"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	green20220302 "github.com/alibabacloud-go/green-20220302/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/gin-gonic/gin"
	"github.com/minio/pkg/env"
	"github.com/xiangtao94/golib/pkg/zlog"
	"net/http"
)

var GreenAliClient *green20220302.Client

func InitAliGreenClient() {
	config := &openapi.Config{
		// 必填，您的 AccessKey ID
		AccessKeyId: tea.String(env.Get("BACKEND_GREEN_ALI_AK", "")),
		// 必填，您的 AccessKey Secret
		AccessKeySecret: tea.String(env.Get("BACKEND_GREEN_ALI_SK", "")),
		RegionId:        tea.String("cn-hangzhou"),
		Endpoint:        tea.String("green-cip.cn-hangzhou.aliyuncs.com"),
		ConnectTimeout:  tea.Int(3000),
		ReadTimeout:     tea.Int(6000),
	}
	// 注意，此处实例化的client请尽可能重复使用，避免重复建立连接，提升检测性能。
	var _err error
	GreenAliClient, _err = green20220302.NewClient(config)
	if _err != nil {
		return
	}
}

func TextCheck(ctx *gin.Context, text string) (pass bool, advise string) {
	pass = true
	if len(text) == 0 {
		return
	}
	// 检测数据。
	serviceParameters, _ := json.Marshal(
		map[string]interface{}{
			"content": text,
		},
	)
	textModerationRequest := &green20220302.TextModerationPlusRequest{
		/*
		   文本检测service：内容安全控制台文本增强版规则配置的serviceCode，示例：chat_detection
		*/
		Service:           tea.String("llm_query_moderation"),
		ServiceParameters: tea.String(string(serviceParameters)),
	}
	// 创建RuntimeObject实例并设置运行参数。
	runtime := &util.RuntimeOptions{}
	runtime.ReadTimeout = tea.Int(10000)
	runtime.ConnectTimeout = tea.Int(10000)
	// 复制代码运行请自行打印API的返回值。
	result, _err := GreenAliClient.TextModerationPlusWithOptions(textModerationRequest, runtime)
	if _err != nil {
		zlog.Errorf(ctx, "文本检测报错，error %s", _err.Error())
		return
	}
	if *result.StatusCode != http.StatusOK {
		zlog.Errorf(ctx, "文本检测返回非200")
		return
	}
	body := result.Body
	if *body.Code != http.StatusOK {
		zlog.Errorf(ctx, "文本检测不成功. code:%d\n", *body.Code)
		return
	}
	data := body.Data
	if len(data.Result) == 0 || *data.Result[0].Label == "nonLabel" {
		return
	}
	zlog.Infof(ctx, "文本检测不通过 requestId:%s, msg:%s\n", *body.RequestId, *data.Result[0].Label)
	pass = false
	if len(data.Advice) > 0 {
		advise = *data.Advice[0].Answer
	}
	return
}
