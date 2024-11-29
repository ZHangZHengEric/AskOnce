package dto_history

import "askonce/components/dto"

type AskReq struct {
	dto.PageParam
	QueryType string `json:"queryType"` // simple complex research
	Query     string `json:"query"`
}
