package jobd

type AtomEsInsertReq struct {
	Corpus            []map[string]any `json:"corpus"`
	Id                string           `json:"id"`
	MapperValueOrPath any              `json:"mapper_value_or_path"`
}

func (entity *JobdApi) AtomEsInsert(inputReq AtomEsInsertReq) (res any, err error) {
	entity.Client.MaxRespBodyLen = -1
	entity.Client.MaxReqBodyLen = -1
	return doTaskProcessString[AtomEsInsertReq](entity, "atomes_es8_insert", inputReq, 50000)
}
