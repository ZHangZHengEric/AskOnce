package components

import "github.com/xiangtao94/golib/pkg/errors"

var ErrorNotLogin = errors.Error{
	Code:    3,
	Message: "用户未登录，请登录！",
}

var ErrorApiAuthError = errors.Error{
	Code:    3,
	Message: "鉴权失败！",
}

var ErrorLoginError = errors.Error{
	Code:    100000,
	Message: "登录失败，账户名或密码不正确!",
}

// 数据库问题
var ErrorMysqlError = errors.Error{
	Code:    100001,
	Message: "数据库操作失败！",
}

var ErrorJobdError = errors.Error{
	Code:    100002,
	Message: "任务执行异常",
}

var ErrorKdbDataRedoError = errors.Error{
	Code:    100003,
	Message: "数据正在处理或处理完成，不允许操作！",
}

var ErrorQueryError = errors.Error{
	Code:    100004,
	Message: "搜索异常，请联系管理员！",
}

var ErrorKdbEmptyError = errors.Error{
	Code:    100005,
	Message: "知识库没有数据，请添加数据！",
}

var ErrorChatError = errors.Error{
	Code:    100006,
	Message: "对话异常，请联系管理员！",
}

var ErrorFormatError = errors.Error{
	Code:    100007,
	Message: "文件格式不支持！",
}

var ErrorFileClientError = errors.Error{
	Code:    100008,
	Message: "文件服务器初始化失败！",
}

var ErrorKdbExist = errors.Error{
	Code:    100009,
	Message: "知识库已存在！请修改名称",
}

var ErrorKdbNoOperate = errors.Error{
	Code:    100010,
	Message: "知识库不能操作！",
}

var ErrorFileUploadError = errors.Error{
	Code:    100011,
	Message: "文件上传失败！",
}

var ErrorFileNoExist = errors.Error{
	Code:    100012,
	Message: "文件不存在！",
}

var ErrorFileNoContent = errors.Error{
	Code:    100013,
	Message: "文件内容未识别，请重新上传！",
}

var ErrorTextSplitError = errors.Error{
	Code:    100014,
	Message: "文本切分失败！",
}

var ErrorDocNoExitError = errors.Error{
	Code:    100015,
	Message: "文档不存在！",
}

var ErrorShareEmpty = errors.Error{
	Code:    100016,
	Message: "分享不正确或已过期，请确认！",
}

// *************************************//

var ErrorUserExistError = errors.Error{
	Code:    200001,
	Message: "用户已存在！",
}
var ErrorUserNotExistError = errors.Error{
	Code:    200002,
	Message: "用户不存在！",
}

var ErrorUserAddError = errors.Error{
	Code:    200003,
	Message: "用户添加失败！",
}

var ErrorLoginPwdError = errors.Error{
	Code:    200004,
	Message: "用户不存在或密码错误！",
}

var ErrorAskSessionError = errors.Error{
	Code:    100011,
	Message: "会话失败，请重新提问！",
}

var ErrorTextCheckError = errors.Error{
	Code:    100008,
	Message: "文本检测不通过，请更换问题！",
}

var ErrorAskSessionNoExist = errors.Error{
	Code:    100009,
	Message: "session不存在！",
}

var ErrorQueryEmpty = errors.Error{
	Code:    100010,
	Message: "未搜索到内容，请换个问题提问！",
}
