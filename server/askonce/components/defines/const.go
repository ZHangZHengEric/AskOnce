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
