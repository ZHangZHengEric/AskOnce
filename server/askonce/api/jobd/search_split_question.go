package jobd

type SplitQuestionReq struct {
	Question  string `json:"question"`
	Id        string `json:"id"`
	ModelName string `json:"model_name"`
}

type SplitQuestionRes struct {
	Questions []string `json:"questions"`
}

func (entity *JobdApi) SplitQuestion(question string) (res *SplitQuestionRes, err error) {
	input := &SplitQuestionReq{
		Question: question,
		Id:       "",
	}
	return doTaskProcess[*SplitQuestionReq, *SplitQuestionRes](entity, "split_question", input, 100000)
}
