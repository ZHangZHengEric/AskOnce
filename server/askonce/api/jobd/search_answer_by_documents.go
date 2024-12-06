package jobd

import "askonce/components/dto/dto_search"

type AnswerByDocumentsReq struct {
	Id           string                          `json:"id"`
	Question     string                          `json:"question"`
	AnswerStyle  string                          `json:"answer_style"`
	SearchResult []dto_search.CommonSearchOutput `json:"search_result"`
	IsStream     bool                            `json:"is_stream"`
}

type AnswerByDocumentsRes struct {
	Answer string `json:"answer"`
}

func (entity *JobdApi) AnswerByDocuments(question string, answerStyle string, searchResult []dto_search.CommonSearchOutput, f func(data JobdCommonRes) error) (err error) {
	if len(searchResult) == 0 {
		searchResult = make([]dto_search.CommonSearchOutput, 0)
	}
	input := &AnswerByDocumentsReq{
		Id:           "",
		AnswerStyle:  answerStyle,
		Question:     question,
		SearchResult: searchResult,
		IsStream:     true,
	}
	return doTaskProcessStream[*AnswerByDocumentsReq, AnswerByDocumentsRes](entity, "answer_by_documents", input, 1000000, f)
}

func (entity *JobdApi) AnswerByDocumentsSync(question string, answerStyle string, searchResult []dto_search.CommonSearchOutput) (out AnswerByDocumentsRes, err error) {
	if len(searchResult) == 0 {
		searchResult = make([]dto_search.CommonSearchOutput, 0)
	}
	input := &AnswerByDocumentsReq{
		Id:           "",
		AnswerStyle:  answerStyle,
		Question:     question,
		SearchResult: searchResult,
		IsStream:     false,
	}
	return doTaskProcess[*AnswerByDocumentsReq, AnswerByDocumentsRes](entity, "answer_by_documents", input, 1000000)
}
