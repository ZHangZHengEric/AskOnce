package jobd

import "askonce/components/dto/dto_search"

type ResearchQAConstructReq struct {
	Id               string                          `json:"id"`
	Question         string                          `json:"question"`
	AnswerStyle      string                          `json:"answer_style"`
	SearchResult     []dto_search.CommonSearchOutput `json:"search_result"`
	ModelName        string                          `json:"model_name"`
	Chapter          string                          `json:"chapter"`
	PaperOutlineList []AnswerKeyPointsItem           `json:"paper_outline_list"`
}

func (entity *JobdApi) ResearchQAConstruct(question, answerStyle string, searchResult []dto_search.CommonSearchOutput, chapter string, paperOutlineList []AnswerKeyPointsItem) (res *SimpleQAConstructRes, err error) {
	if len(searchResult) == 0 {
		searchResult = make([]dto_search.CommonSearchOutput, 0)
	}
	input := &ResearchQAConstructReq{
		Id:               "",
		Question:         question,
		AnswerStyle:      answerStyle,
		SearchResult:     searchResult,
		Chapter:          chapter,
		PaperOutlineList: paperOutlineList,
	}
	return doTaskProcess[*ResearchQAConstructReq, *SimpleQAConstructRes](entity, "atom_app_askonce_research_qa_construct", input, 100000)
}
