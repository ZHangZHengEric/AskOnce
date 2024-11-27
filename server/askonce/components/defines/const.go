package defines

/*
*
上下文相关
*/
const (
	LoginInfo = "loginInfo"
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

/*
*
结构化相关json的固定key
*/
type StructuredKey string

const (
	StructuredTitle      StructuredKey = "title"       // 标题
	StructuredContent    StructuredKey = "content"     // 内容
	StructuredStartIndex StructuredKey = "start_index" // 起始下标
	StructuredEndIndex   StructuredKey = "end_index"   // 结尾下标
)

// 强结构化文件结尾
var StructuredExtension = []string{"json", "csv", "xls", "xlsx"}

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
