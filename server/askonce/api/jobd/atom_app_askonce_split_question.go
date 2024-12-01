package jobd

type SplitQuestionReq struct {
	Question  string `json:"question"`
	Id        string `json:"id"`
	ModelName string `json:"model_name"`
}

type SplitQuestionRes struct {
	SubTitles      []string `json:"sub_titles"`
	SearchContents []string `json:"search_contents"`
}

func (entity *JobdApi) SplitQuestion(question string) (res *SplitQuestionRes, err error) {
	input := &SplitQuestionReq{
		Question: question,
		Id:       "",
	}
	return doTaskProcess[*SplitQuestionReq, *SplitQuestionRes](entity, "atom_app_askonce_split_question", input, 100000)
}
