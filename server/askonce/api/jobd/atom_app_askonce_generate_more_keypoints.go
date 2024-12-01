package jobd

type GenerateMoreKeyPointsRes struct {
	MoreKeyPoints []string `json:"more_keypoints"`
}

type GenerateMoreKeyPointsReq struct {
	Id              string                `json:"id"`
	Question        string                `json:"question"`
	AnswerKeyPoints []AnswerKeyPointsItem `json:"answer_keypoints"`
	ModelName       string                `json:"model_name"`
}

func (entity *JobdApi) GenerateMoreKeyPoints(q string, ak []AnswerKeyPointsItem) (res *GenerateMoreKeyPointsRes, err error) {
	input := &GenerateMoreKeyPointsReq{
		Question:        q,
		AnswerKeyPoints: ak,
	}
	return doTaskProcess[*GenerateMoreKeyPointsReq, *GenerateMoreKeyPointsRes](entity, "atom_app_askonce_generate_more_keypoints", input, 100000)
}
