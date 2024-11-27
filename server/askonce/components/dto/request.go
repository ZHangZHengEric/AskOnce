package dto

import "mime/multipart"

type FileUploadReq struct {
	Source string                `form:"source"` // 来源
	File   *multipart.FileHeader `form:"file" binding:"required"`
}

type ModelProviderReq struct {
	ModelType string `form:"modelType" binding:"required"`
}
