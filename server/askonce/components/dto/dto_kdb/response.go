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
	Id           int64  `json:"id"`
	Name         string `json:"name"`
	CreateTime   string `json:"createTime"`
	DataSource   string `json:"dataSource"`
	DocNum       int    `json:"docNum"` // 文档数量
	Cover        string `json:"cover"`
	DefaultColor bool   `json:"defaultColor"`
	Creator      string `json:"creator"`
	Type         int    `json:"type"`  // 1 公共数据
	Intro        string `json:"intro"` // 介绍
}

type InfoRes struct {
	KdbId          int64          `json:"kdbId"`
	Name           string         `json:"name"`  // 知识库名称
	Intro          string         `json:"intro"` // 知识库介绍
	Cover          string         `json:"cover"` // 知识库封面
	CreatedAt      int64          `json:"createdAt"`
	CreatedBy      string         `json:"createdBy"`
	UpdatedAt      int64          `json:"updatedAt"`
	DataSourceType DataSourceType `json:"dataSourceType"`
	WordCount      int64          `json:"wordCount"`
	DocumentCount  int64          `json:"documentCount"`
	dto.KdbSetting
}

type CoversRes struct {
	List []CoverItem `json:"list"`
}

type CoverItem struct {
	Id   int64  `json:"id"`
	Type string `json:"type"`
	Url  string `json:"url"`
}

type AuthRes struct {
	AuthType int `json:"authType"`
}

type UserListRes struct {
	List  []UserListItem `json:"list"`
	Total int64          `json:"total"`
}

type UserListItem struct {
	UserId   string `json:"userId"`
	UserName string `json:"userName"`
	JoinTime string `json:"joinTime"`
}

type UserQueryRes struct {
	List  []UserQueryItem `json:"list"`
	Total int64           `json:"total"`
}

type UserQueryItem struct {
	UserId   string `json:"userId"`
	UserName string `json:"userName"`
}

type GenShareCodeRes struct {
	ShareCode string `json:"shareCode" `
}

type ShareCodeInfoRes struct {
	Creator  string `json:"creator"`
	KdbName  string `json:"kdbName"`
	AuthType int    `json:"authType"`
}

type LoadProcessRes struct {
	Success    int64 `json:"success"`
	Fail       int64 `json:"fail"`
	InProgress int64 `json:"inProgress"`
	Waiting    int64 `json:"waiting"`

	Total       int64 `json:"total"`
	TaskProcess int64 `json:"taskProcess"`
}
