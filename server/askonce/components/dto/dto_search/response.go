package dto_search

import "rag-backend/components/dto"

type AskRes struct {
	Type    dto.ChatType `json:"type"`
	Content string       `json:"content"`
	Id      string       `json:"id,omitempty"`
	Title   string       `json:"title,omitempty"`
}

type HistoryRes struct {
	List  []HistoryItem `json:"list"`
	Total int64         `json:"total"`
}

type HistoryItem struct {
	SessionId  string `json:"sessionId"`
	Name       string `json:"name"`
	CreateTime string `json:"createTime"`
}

type HisDetailRes struct {
	List  []HistoryDetailItem `json:"list"`
	Total int64               `json:"total"`
}

type HistoryDetailItem struct {
	Id          int64   `json:"id"`
	Role        string  `json:"role"`
	Content     string  `json:"content"`
	ContentTime string  `json:"contentTime"`
	UseData     UseData `json:"useData"` //使用到的数据
}

type UseData struct {
	KdbSearch []UseDataKdbSearch `json:"kdbSearch,omitempty"` // 知识库搜索结果
}

type UseDataKdbSearch struct {
	Content string `json:"content"`
	Id      int64  `json:"id"`
	Title   string `json:"title"`
}
