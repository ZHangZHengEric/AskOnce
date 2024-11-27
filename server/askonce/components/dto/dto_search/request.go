package dto_search

import "rag-backend/components/dto"

type AskReq struct {
	SessionId string `json:"sessionId" binding:"required"`
	Question  string `json:"question" binding:"required"` // 问题
	KdbId     int64  `json:"kdbId" binding:"required"`    //
}

type SettingReq struct {
	SearchEngine string `json:"searchEngine" binding:"required"`
}

type HisReq struct {
	dto.PageParam
}

type HisDetailReq struct {
	SessionId string `form:"sessionId" binding:"required"`
}
