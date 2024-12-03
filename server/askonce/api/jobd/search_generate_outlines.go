package jobd

type AnswerOutlineReq struct {
	Id        string `json:"id"`
	Answer    string `json:"answer"`
	ModelName string `json:"model_name"`
}

type AnswerOutlineRes struct {
	AnswerOutline []AnswerOutlineItem `json:"answer_outline"`
}

type AnswerOutlineItem struct {
	Level   string `json:"level"`
	Content string `json:"content"`
}

func (entity *JobdApi) AnswerOutline(answer string) (res *AnswerOutlineRes, err error) {
	input := &AnswerOutlineReq{
		Id:     "",
		Answer: answer,
	}
	return doTaskProcess[*AnswerOutlineReq, *AnswerOutlineRes](entity, "generate_outlines", input, 100000)
}
