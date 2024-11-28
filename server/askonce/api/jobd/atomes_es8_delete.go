package jobd

type AtomEsDeleteReq struct {
	DocIds            []string `json:"doc_ids"`
	Id                string   `json:"id"`
	MapperValueOrPath any      `json:"mapper_value_or_path"`
	DeleteAll         bool     `json:"delete_all"`
}

type AtomEsDeleteRes struct {
	DeleteResults []string `json:"delete_results"`
}

func (entity *JobdApi) AtomEsDelete(inputReq *AtomEsDeleteReq) (res *AtomEsDeleteRes, err error) {
	return doTaskProcess[*AtomEsDeleteReq, *AtomEsDeleteRes](entity, "atomes_es8_delete", inputReq, 10000)
}
