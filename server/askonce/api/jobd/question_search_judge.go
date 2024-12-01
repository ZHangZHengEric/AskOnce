package jobd

type QuestionSearchJudgeReq struct {
	Id        string   `json:"id"`
	Questions []string `json:"questions"`
	Threshold float32  `json:"threshold"`
}

type QuestionSearchJudgeRes struct {
	NeedSearch []bool `json:"need_search"`
}

func (entity *JobdApi) QuestionSearchJudge(questions []string) (res *QuestionSearchJudgeRes, err error) {
	input := &QuestionSearchJudgeReq{
		Id:        "",
		Questions: questions,
		Threshold: 0.4,
	}
	return doTaskProcess[*QuestionSearchJudgeReq, *QuestionSearchJudgeRes](entity, "question_search_judge_api", input, 100000)
}
