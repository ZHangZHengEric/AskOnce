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
	Name     string `json:"name"`      // 知识库名称
	Type     string `json:"type"`      // 知识库类型 doc sql
	Intro    string `json:"intro"`     // 知识库介绍
	Language string `json:"language" ` // 知识库语言
}

type UpdateReq struct {
	KdbId   int64  `json:"kdbId"`
	Name    string `json:"name"`    // 知识库名称
	Intro   string `json:"intro"`   // 知识库介绍
	CoverId int64  `json:"coverId"` // 知识库封面
	*dto.KdbSetting
}

type DeleteReq struct {
	KdbId int64 `json:"kdbId" binding:"required"`
}

type DeleteSelfReq struct {
	KdbId int64 `json:"kdbId" binding:"required"`
}

type AuthReq struct {
	KdbId int64 `json:"kdbId" binding:"required"`
}

type InfoReq struct {
	KdbId int64 `form:"kdbId" binding:"required"`
}

type UserListReq struct {
	KdbId     int64  `json:"kdbId" binding:"required"`
	AuthType  int    `json:"authType" binding:"required"`
	QueryName string `json:"queryName"`
	dto.PageParam
}

type UserQueryReq struct {
	QueryName string `json:"queryName"`
}

type UserAddReq struct {
	KdbId    int64    `json:"kdbId" binding:"required"`
	AuthType int      `json:"authType" binding:"required"`
	UserIds  []string `json:"userIds"`
}

type UserDeleteReq struct {
	KdbId    int64    `json:"kdbId" binding:"required"`
	AuthType int      `json:"authType"`
	UserIds  []string `json:"userIds"`
}

type GenShareCodeReq struct {
	KdbId    int64 `json:"kdbId" binding:"required"`
	AuthType int   `json:"authType" binding:"required"`
}

type VerifyShareCodeReq struct {
	ShareCode string `form:"shareCode" binding:"required"`
}

type InfoShareCodeReq struct {
	ShareCode string `form:"shareCode" binding:"required"`
}

type RecallReq struct {
	KdbId int64  `json:"kdbId" binding:"required"`
	Query string `query:"query" binding:"required"`
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
