package dto_kdb_doc

import (
	"askonce/components/dto"
	"mime/multipart"
)

type AddReq struct {
	KdbId int64                 `json:"kdbId" form:"kdbId" binding:"required"`
	Type  string                `json:"type" form:"type" binding:"required"`
	File  *multipart.FileHeader `json:"file" form:"file"`
	ImportText
	ImportDataBase
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
	Title    string                 `json:"title" form:"title"`
	Text     string                 `json:"text" form:"text"`
	Metadata map[string]interface{} `json:"metadata" form:"metadata"`
}

type ImportDataBase struct {
	DbName    string `json:"dbName" form:"dbName"`
	DbType    string `json:"dbType" form:"dbType"`
	DbHost    string `json:"dbHost" form:"dbHost"`
	DbPort    int    `json:"dbPort" form:"dbPort"`
	DbUser    string `json:"dbUser" form:"dbUser"`
	DbPwd     string `json:"dbPwd" form:"dbPwd"`
	DbComment string `json:"dbComment" form:"dbComment"`
}

type ListReq struct {
	KdbId       int64  `json:"kdbId" binding:"required"`
	QueryName   string `json:"queryName"`
	QueryStatus []int  `json:"queryStatus"`
	dto.PageParam
}

type DeleteReq struct {
	KdbId  int64 `json:"kdbId" binding:"required"`
	DataId int64 `json:"dataId" binding:"required"`
}

type InfoReq struct {
	KdbId  int64 `json:"kdbId" form:"kdbId" binding:"required"`
	DataId int64 `json:"dataId" form:"dataId" binding:"required"`
}

type RedoReq struct {
	KdbId  int64 `json:"kdbId" binding:"required"`
	DataId int64 `json:"dataId" binding:"required"`
}
type RecallReq struct {
	KdbId int64  `json:"kdbId" binding:"required"`
	Query string `query:"query" binding:"required"`
}

type TaskProcessReq struct {
	KdbId  int64  `json:"kdbId" form:"kdbId" binding:"required"`
	TaskId string `json:"taskId" form:"taskId" binding:"required"`
}

type TaskRedoReq struct {
	KdbId  int64  `json:"kdbId" form:"kdbId" binding:"required"`
	TaskId string `json:"taskId" form:"taskId" binding:"required"`
}
