package jobd

import "askonce/components/dto/dto_search"

type SimpleQAConstructReq struct {
	Id           string                          `json:"id"`
	Question     string                          `json:"question"`
	AnswerStyle  string                          `json:"answer_style"`
	SearchResult []dto_search.CommonSearchOutput `json:"search_result"`
	ModelName    string                          `json:"model_name"`
}

type SimpleQAConstructRes struct {
	Prompt        string            `json:"prompt"`
	GenerateParam GenerateParamItem `json:"generate_param"`
}

type GenerateParamItem struct {
	Temperature     float64 `json:"temperature"`
	PresencePenalty float64 `json:"presence_penalty"`
	MaxNewTokens    int     `json:"max_new_tokens"`
}

func (entity *JobdApi) SimpleQAConstruct(question, answerStyle string, searchResult []dto_search.CommonSearchOutput) (res *SimpleQAConstructRes, err error) {
	if len(searchResult) == 0 {
		searchResult = make([]dto_search.CommonSearchOutput, 0)
	}
	input := &SimpleQAConstructReq{
		Id:           "",
		Question:     question,
		AnswerStyle:  answerStyle,
		SearchResult: searchResult,
	}
	return doTaskProcess[*SimpleQAConstructReq, *SimpleQAConstructRes](entity, "atom_app_askonce_simple_qa_construct", input, 100000)
}
