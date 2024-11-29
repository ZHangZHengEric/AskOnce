package jobd

type AtomEsRankReq struct {
	Query         string   `json:"query"`
	RecallResults []string `json:"recall_results"`
	Id            string   `json:"id"`
}

func (entity *JobdApi) AtomEsRank(query string, recallResults []string) (res []float32, err error) {
	input := &AtomEsRankReq{
		Query:         query,
		RecallResults: recallResults,
	}
	return doTaskProcess[*AtomEsRankReq, []float32](entity, "bgem3_rerank", input, 10000)
}
