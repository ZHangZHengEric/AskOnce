package dto_search

import "askonce/components/dto"

type CaseReq struct {
	KdbId int64 `json:"kdbId"` // 为kdb时有值
}

type SettingReq struct {
	SearchEngine string `json:"searchEngine" binding:"required"`
}

type HisReq struct {
	SessionId string `json:"sessionId"`
}

type HisDetailReq struct {
	SessionId string `form:"sessionId" binding:"required"`
}

type AskReq struct {
	SessionId string `json:"sessionId" binding:"required"`
	Question  string `json:"question" binding:"required"` // 问题
	Type      string `json:"type" binding:"required"`     // simple complex research
	KdbId     int64  `json:"kdbId"`                       // 为kdb时有值
}

type ReferReq struct {
	SessionId string `json:"sessionId"`
}

type OutlineReq struct {
	SessionId string `json:"sessionId"`
}

type RelationReq struct {
	SessionId string `json:"sessionId"`
}

type UnlikeReq struct {
	SessionId string   `json:"sessionId"`
	Reasons   []string `json:"reasons"`
}

type ProcessReq struct {
	SessionId string `json:"sessionId"`
}

type KdbListReq struct {
	Query     string `json:"query" `
	OrderType int    `json:"orderType"` // 0,1 创建时间倒序 2 创建时间正序 3 最近使用
	dto.PageParam
}

type ChatAskReq struct {
	SessionId string `json:"sessionId"`
	Question  string `json:"question" binding:"required"` // 问题
	Type      string `json:"type" `                       // simple complex research
	KdbId     int64  `json:"kdbId"`                       // 为kdb时有值
}
