package jobd

type NetRagAssessmentReq struct {
	Id        string  `json:"id"`
	Question  string  `json:"question"`
	Threshold float32 `json:"threshold"`
}

type NetRagAssessmentRes struct {
	Result bool `json:"result"`
}

func (entity *JobdApi) NetRagAssessment(question string) (useRag bool, err error) {
	input := &NetRagAssessmentReq{
		Id:        "",
		Question:  question,
		Threshold: 0.4,
	}
	res, err := doTaskProcess[*NetRagAssessmentReq, NetRagAssessmentRes](entity, "net_rag_assessment", input, 100000)
	if err != nil {
		entity.LogErrorf("net_rag_assessment error: %s", err.Error())
		return false, nil
	}
	return res.Result, nil
}
