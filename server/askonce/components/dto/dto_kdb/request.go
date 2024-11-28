package dto_kdb

import (
	"askonce/components/dto"
	"mime/multipart"
)

type ListReq struct {
	QueryName string `json:"queryName"`
	dto.PageParam
}

type DataAddReq struct {
	KdbId int64                 `json:"kdbId" form:"kdbId" binding:"required"`
	File  *multipart.FileHeader `json:"file" form:"file"`
}

type DataListReq struct {
	KdbId     int64  `json:"kdbId" form:"kdbId" binding:"required"`
	QueryName string `query:"queryName"`
	dto.PageParam
}

type DataDeleteReq struct {
	KdbId     int64 `json:"kdbId" form:"kdbId" binding:"required"`
	KdbDataId int64 `json:"kdbDataId" binding:"required"`
}

type DataRedoReq struct {
	KdbId     int64 `json:"kdbId" form:"kdbId" binding:"required"`
	KdbDataId int64 `json:"kdbDataId" binding:"required"`
}

type AddReq struct {
	Name  string `json:"name"`  // 知识库名称
	Type  string `json:"type"`  // 知识库类型 doc sql
	Intro string `json:"intro"` // 知识库介绍
}

type UpdateReq struct {
	KdbId int64  `json:"kdbId"`
	Name  string `json:"name"`  // 知识库名称
	Intro string `json:"intro"` // 知识库介绍
	*dto.KdbSetting
}

type InfoReq struct {
	KdbId int64 `form:"kdbId" binding:"required"`
}

type DocInfoReq struct {
	KdbId    int64 `json:"kdbId" form:"kdbId" binding:"required"`
	KdbDocId int64 `json:"kdbDocId" form:"kdbDocId" binding:"required"`
}

type DocDeleteReq struct {
	KdbId    int64 `json:"kdbId" form:"kdbId" binding:"required"`
	KdbDocId int64 `json:"kdbDocId" form:"kdbDocId" binding:"required"`
}

type DocListReq struct {
	KdbId     int64  `json:"kdbId" form:"kdbId" binding:"required"`
	QueryName string `json:"queryName"`
	dto.PageParam
}

type DocRenameReq struct {
	KdbId    int64  `json:"kdbId" form:"kdbId" binding:"required"`
	KdbDocId int64  `json:"kdbDocId" form:"kdbDocId" binding:"required"`
	Name     string `json:"name"`
}

type DocStatusSettingReq struct {
	KdbId    int64 `json:"kdbId" form:"kdbId" binding:"required"`
	KdbDocId int64 `json:"kdbDocId" form:"kdbDocId" binding:"required"`
}

type DocSegmentReq struct {
	KdbId    int64 `json:"kdbId" form:"kdbId" binding:"required"`
	KdbDocId int64 `json:"kdbDocId" form:"kdbDocId" binding:"required"`
	dto.PageParam
}
