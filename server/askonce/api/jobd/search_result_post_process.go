// Package jobd -----------------------------
// @file      : search_result_post_process.go
// @author    : xiangtao
// @contact   : xiangtao1994@gmail.com
// @time      : 2024/12/14 20:14
// -------------------------------------------
package jobd

import "askonce/es"

type SearchResultPostProcessReq struct {
	Id           string               `json:"id"`
	Question     string               `json:"question"`
	SearchResult []*es.CommonDocument `json:"search_result"`
}

type SearchResultPostProcessRes struct {
	SearchResult []*es.CommonDocument `json:"search_result"`
}

func (entity *JobdApi) SearchResultPostProcess(question string, input []*es.CommonDocument) (res *SearchResultPostProcessRes, err error) {
	req := &SearchResultPostProcessReq{
		Id:           "",
		Question:     question,
		SearchResult: input,
	}
	return doTaskProcess[*SearchResultPostProcessReq, *SearchResultPostProcessRes](entity, "search_result_post_process", req, 100000)
}
