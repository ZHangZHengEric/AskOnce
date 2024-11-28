package jobd

type AtomEsText2EmbInput struct {
	Sentences []string `json:"sentences"`
	Id        string   `json:"id"`
}

func (entity *JobdApi) Embedding(texts []string) (res [][]float32, err error) {
	input := &AtomEsText2EmbInput{
		Sentences: texts,
	}
	entity.Client.MaxRespBodyLen = -1
	entity.Client.MaxReqBodyLen = -1
	res, err = doTaskProcess[*AtomEsText2EmbInput, [][]float32](entity, "text2emb", input, 100000)
	if res == nil {
		res = make([][]float32, 0)
	}
	return res, err
}

func (entity *JobdApi) EmbeddingForQuery(sentences []string) (res [][]float32, err error) {
	input := &AtomEsText2EmbInput{
		Sentences: sentences,
	}
	entity.Client.MaxRespBodyLen = -1
	res, err = doTaskProcess[*AtomEsText2EmbInput, [][]float32](entity, "query2emb", input, 100000)
	if res == nil {
		res = make([][]float32, 0)
	}
	return res, err
}
