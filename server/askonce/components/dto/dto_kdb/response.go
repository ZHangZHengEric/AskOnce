package dto_kdb

import (
	"askonce/components/dto"
)

type AddRes struct {
	KdbId int64 `json:"kdbId"`
}

type ListResp struct {
	List  []ListItem `json:"list"`
	Total int64      `json:"total"`
}

type ListItem struct {
	Id         int64  `json:"id"`
	Name       string `json:"name"`
	CreateTime string `json:"createTime"`
	DataSource string `json:"dataSource"`
	Type       string `json:"type"`   // 类型  private 私有知识库 public 共有知识库
	DocNum     int    `json:"docNum"` // 文档数量
}

type DataListResp struct {
	List  []DataListItem `json:"list"`
	Total int64          `json:"total"`
}

type DataListItem struct {
	KdbDataId  int64  `json:"kdbDataId"`
	Type       string `json:"type"`
	DataSuffix string `json:"dataSuffix"`
	DataName   string `json:"dataName"`
	DataPath   string `json:"dataPath"`
	Status     int    `json:"status"` // 0正在构建到知识库 1 成功 2 失败
	CreateTime string `json:"createTime"`
}

type InfoRes struct {
	Id             int64          `json:"id"`
	Name           string         `json:"name"`
	Description    string         `json:"description"`
	CreatedAt      int64          `json:"createdAt"`
	CreatedBy      string         `json:"createdBy"`
	UpdatedAt      int64          `json:"updatedAt"`
	DataSourceType DataSourceType `json:"dataSourceType"`
	WordCount      int64          `json:"wordCount"`
	DocumentCount  int64          `json:"documentCount"`
	dto.KdbSetting
}

type DocAddRes struct {
	KdbDocId int64 `json:"kdbDocId"`
}

type DocInfoRes struct {
	Id             int64  `json:"id"`
	Name           string `json:"name"`
	Enable         bool   `json:"enable"`         // 是否启用
	IndexingStatus string `json:"indexingStatus"` // 索引建立状态
	WordCount      int64  `json:"wordCount"`      // 字符数
	HitCount       int64  `json:"hitCount"`       // 召回数
	CreatedAt      int64  `json:"createdAt"`
	CreatedBy      string `json:"createdBy"`
	UpdatedAt      int64  `json:"updatedAt"`
}

type DocListRes struct {
	List  []DocListItem `json:"list"`
	Total int64         `json:"total"`
}

type DocListItem struct {
	Id             int64  `json:"id"`
	Name           string `json:"name"`
	Enable         bool   `json:"enable"`         // 是否启用
	IndexingStatus string `json:"indexingStatus"` // 索引建立状态
	WordCount      int64  `json:"wordCount"`      // 字符数
	HitCount       int64  `json:"hitCount"`       // 召回数
	CreatedAt      int64  `json:"createdAt"`
	CreatedBy      string `json:"createdBy"`
	UpdatedAt      int64  `json:"updatedAt"`
}
