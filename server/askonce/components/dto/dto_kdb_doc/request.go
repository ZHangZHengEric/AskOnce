package dto_kdb_doc

import (
	"askonce/components/dto"
	"mime/multipart"
)

type AddReq struct {
	KdbId int64                 `json:"kdbId" form:"kdbId" binding:"required"`
	Type  string                `json:"type" form:"type" binding:"required"`
	File  *multipart.FileHeader `json:"file" form:"file"`
	Text  string                `json:"text" form:"text"`
	Title string                `json:"title" form:"title"`
}

type AddZipReq struct {
	KdbId  int64  `json:"kdbId" form:"kdbId" binding:"required"`
	ZipUrl string `json:"zipUrl" form:"zipUrl" binding:"required"`
}

type AddByBatchTextReq struct {
	KdbName    string       `json:"kdbName"`
	AutoCreate bool         `json:"autoCreate"`
	Docs       []ImportText `json:"docs"`
}

type ImportText struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

type ListReq struct {
	KdbId       int64  `json:"kdbId" binding:"required"`
	QueryName   string `json:"queryName"`
	QueryStatus []int  `json:"queryStatus"`
	dto.PageParam
}

type DeleteReq struct {
	KdbId int64 `json:"kdbId" binding:"required"`
	DocId int64 `json:"dataId" binding:"required"`
}

type RedoReq struct {
	KdbId int64 `json:"kdbId" binding:"required"`
	DocId int64 `json:"dataId" binding:"required"`
}
type RecallReq struct {
	KdbId int64  `json:"kdbId" binding:"required"`
	Query string `query:"query" binding:"required"`
}
