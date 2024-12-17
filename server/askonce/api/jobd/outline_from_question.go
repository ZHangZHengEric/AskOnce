package jobd

import "askonce/components/dto/dto_search"

type QuestionOutlineReq struct {
	Id           string                          `json:"id"`
	Question     string                          `json:"question"`
	SearchResult []dto_search.CommonSearchOutput `json:"search_result"`
}

type QuestionOutlineRes struct {
	Outline []Outline `json:"outline"`
}

func (entity *JobdApi) QuestionOutline(question string, searchResult []dto_search.CommonSearchOutput) (res *QuestionOutlineRes, err error) {
	input := &QuestionOutlineReq{
		Id:           "",
		Question:     question,
		SearchResult: searchResult,
	}
	return doTaskProcess[*QuestionOutlineReq, *QuestionOutlineRes](entity, "generate_outlines_from_question", input, 100000)
}
