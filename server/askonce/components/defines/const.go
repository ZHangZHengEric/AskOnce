package defines

import "github.com/xiangtao94/golib/flow"

/*
*
上下文相关
*/
const (
	LoginInfo = "loginInfo"
)

const (
	COOKIE_KEY         = "askOnceSession"
	COOKIE_PATH        = "/"
	COOKIE_DEFAULT_AGE = flow.EXPIRE_TIME_1_WEEK
)

/*
*
minio存储桶定义
*/
const (
	BucketOrigin string = "origin" // 原始数据存储桶

	BucketTmp string = "tmp" // 临时文件

	BucketText string = "text" // 解析完成后文本数据存储桶
)

type GPTSource string

const (
	GPTSourceOpenAI   GPTSource = "openai"
	GPTSourceKimi     GPTSource = "kimi"
	GPTSourceBaiChuan GPTSource = "baichuan"
	GPTSourceQwen     GPTSource = "qwen"
	GPTSourceGlm      GPTSource = "glm"
	GPTSourceMiniMax  GPTSource = "minimax"
	GPTSourceSkylark  GPTSource = "skylark" // 字节豆包大模型 v3接口版本
	GPTSourceAzure    GPTSource = "azure"
)

type GPTMethodType string

const (
	MethodTypeChat      GPTMethodType = "chat"
	MethodTypeEmbedding GPTMethodType = "embedding"
)
