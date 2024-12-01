package jobd

type MergeAnswerReq struct {
	Id        string   `json:"id"`
	Question  string   `json:"question"`
	Answers   []string `json:"answers"`
	Titles    []string `json:"titles"`
	ModelName string   `json:"model_name"`
}

type MergeAnswerRes struct {
	Answer string `json:"answer"`
}

func (entity *JobdApi) MergeAnswers(question string, titles, answers []string) (res *SimpleQAConstructRes, err error) {
	input := &MergeAnswerReq{
		Question: question,
		Answers:  answers,
		Titles:   titles,
	}
	return doTaskProcess[*MergeAnswerReq, *SimpleQAConstructRes](entity, "atom_app_askonce_merge_answers", input, 100000)
}
