package dto_knowledge

import (
	"askOnce/components/dto"
	"mime/multipart"
)

type DataAddReq struct {
	KdbId int64                 `json:"kdbId" form:"kdbId" binding:"required"`
	Type  string                `json:"type" form:"type" binding:"required"`
	File  *multipart.FileHeader `json:"file" form:"file"`
	Text  string                `json:"text" form:"text"`
	Title string                `json:"title" form:"title"`
}

type DataBatchAddReq struct {
	KdbId int64                 `json:"kdbId" form:"kdbId" binding:"required"`
	File  *multipart.FileHeader `json:"file" form:"file" binding:"required"`
}

type DataListReq struct {
	KdbId     int64  `json:"kdbId" binding:"required"`
	QueryName string `query:"queryName"`
	dto.PageParam
}

type DataDeleteReq struct {
	KdbId  int64 `json:"kdbId" binding:"required"`
	DataId int64 `json:"dataId" binding:"required"`
}

type DataRedoReq struct {
	KdbId  int64 `json:"kdbId" binding:"required"`
	DataId int64 `json:"dataId" binding:"required"`
}
