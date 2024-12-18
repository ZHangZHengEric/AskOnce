package jobd

import "askonce/components/dto/dto_search"

type AnswerByDocumentsReq struct {
	Id             string                          `json:"id"`
	Question       string                          `json:"question"`
	AnswerStyle    string                          `json:"answer_style"`
	AnswerOutlines []Outline                       `json:"answer_outlines"`
	SearchResult   []dto_search.CommonSearchOutput `json:"search_result"`
	IsStream       bool                            `json:"is_stream"`
	SearchCode     string                          `json:"search_code"`
}

type AnswerByDocumentsRes struct {
	Answer string `json:"answer"`
}

func (entity *JobdApi) AnswerByDocuments(sessionId string, question string, answerStyle string, searchResult []dto_search.CommonSearchOutput, f func(data JobdCommonRes) error) (err error) {
	if len(searchResult) == 0 {
		searchResult = make([]dto_search.CommonSearchOutput, 0)
	}
	input := &AnswerByDocumentsReq{
		Id:           "",
		AnswerStyle:  answerStyle,
		Question:     question,
		SearchCode:   sessionId,
		SearchResult: searchResult,
		IsStream:     true,
	}
	return doTaskProcessStream[*AnswerByDocumentsReq, AnswerByDocumentsRes](entity, "answer_by_documents", input, 1000000, f)
}

func (entity *JobdApi) AnswerByDocumentsSync(sessionId string, question string, answerStyle string, searchResult []dto_search.CommonSearchOutput, outline []Outline) (out AnswerByDocumentsRes, err error) {
	if len(searchResult) == 0 {
		searchResult = make([]dto_search.CommonSearchOutput, 0)
	}
	input := &AnswerByDocumentsReq{
		Id:             "",
		AnswerOutlines: outline,
		AnswerStyle:    answerStyle,
		Question:       question,
		SearchResult:   searchResult,
		SearchCode:     sessionId,
		IsStream:       false,
	}
	return doTaskProcess[*AnswerByDocumentsReq, AnswerByDocumentsRes](entity, "answer_by_documents", input, 1000000)
}
