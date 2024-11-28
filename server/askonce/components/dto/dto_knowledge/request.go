package dto_knowledge

import (
	"askOnce/components/dto"
	"mime/multipart"
)

type ListReq struct {
	QueryName string `query:"queryName"`
	dto.PageParam
}

type AddReq struct {
	Name     string `json:"name" binding:"required"`     // 知识库名称
	Type     string `json:"type" binding:"required"`     // 知识库类型 doc sql
	Intro    string `json:"intro"`                       // 知识库介绍
	Language string `json:"language" binding:"required"` // 知识库语言
}

type UpdateReq struct {
	KdbId   int64  `json:"kdbId"`
	Name    string `json:"name"`    // 知识库名称
	Intro   string `json:"intro"`   // 知识库介绍
	CoverId int64  `json:"coverId"` // 知识库封面
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
	KdbId    int64   `json:"kdbId" binding:"required"`
	AuthType int     `json:"authType" binding:"required"`
	UserIds  []int64 `json:"userIds"`
}

type UserDeleteReq struct {
	KdbId    int64   `json:"kdbId" binding:"required"`
	AuthType int     `json:"authType" binding:"required"`
	UserIds  []int64 `json:"userIds"`
}

type DetailReq struct {
	KdbId int64 `json:"kdbId"  form:"kdbId" binding:"required"`
}

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

type SearchReq struct {
	KdbId int64  `json:"kdbId" binding:"required"`
	Query string `query:"query" binding:"required"`
}

type DataRedoReq struct {
	KdbId  int64 `json:"kdbId" binding:"required"`
	DataId int64 `json:"dataId" binding:"required"`
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

type SearchAdminReq struct {
	IndexName string `json:"indexName" binding:"required"`
	Query     string `json:"query" binding:"required"`
	Size      int    `json:"size"`
}
