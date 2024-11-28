package jobd

import "encoding/json"

type AtomEsSearchReq struct {
	Id                string       `json:"id"`
	SearchType        string       `json:"search_type"`
	MapperValueOrPath any          `json:"mapper_value_or_path"`
	SearchBody        []SearchBody `json:"search_body"`
}

type SearchBody struct {
	Knn   *SearchBodyKnn  `json:"knn,omitempty"`
	Query *SearchBodyBm25 `json:"query,omitempty"`
	Size  *int            `json:"size,omitempty"`
}

type SearchBodyBm25 struct {
	Match struct {
		DocContent interface{} `json:"doc_content"`
	} `json:"match"`
}

type SearchBodyKnn struct {
	Field         string `json:"field"`
	QueryVector   any    `json:"query_vector"`
	K             int    `json:"k"`
	NumCandidates int    `json:"num_candidates"`
}

type SearchOutput struct {
	Source SearchOutputSource `json:"source"`
	Score  float32            `json:"score"`
}

type SearchOutputSource struct {
	DocId       int64  `json:"doc_id"`
	DocContent  string `json:"doc_content"`
	DataSplitId int64  `json:"data_split_id"`
}

func (entity *JobdApi) AtomEsSearch(emb any, query string, querySize int, mapValue string) (res []SearchOutput, err error) {
	inputReq := &AtomEsSearchReq{
		SearchType:        "all",
		MapperValueOrPath: json.RawMessage(mapValue),
		SearchBody: []SearchBody{
			{
				Knn: &SearchBodyKnn{
					Field:         "emb",
					QueryVector:   emb,
					K:             querySize,
					NumCandidates: 500,
				},
			},
			{
				Query: &SearchBodyBm25{
					Match: struct {
						DocContent interface{} `json:"doc_content"`
					}(struct{ DocContent interface{} }{
						DocContent: query,
					}),
				},
				Size: &querySize,
			},
		},
	}
	return doTaskProcess[*AtomEsSearchReq, []SearchOutput](entity, "atomes_es8_search", inputReq, 10000)
}
