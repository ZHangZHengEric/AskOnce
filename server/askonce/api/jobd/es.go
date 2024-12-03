package jobd

import "encoding/json"

type ESSearchReq struct {
	Id                string         `json:"id"`
	SearchType        string         `json:"search_type"`
	MapperValueOrPath any            `json:"mapper_config"`
	SearchBody        []ESSearchBody `json:"search_body"`
}

type ESSearchBody struct {
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

type ESSearchOutput struct {
	Source SearchOutputSource `json:"source"`
	Score  float32            `json:"score"`
}

type SearchOutputSource struct {
	DocId       string `json:"doc_id"`
	DocContent  string `json:"doc_content"`
	DataSplitId int64  `json:"data_split_id"`
	Start       int    `json:"start"`
	End         int    `json:"end"`
}

func (entity *JobdApi) EsSearch(emb any, query string, querySize int, mapValue string) (res []ESSearchOutput, err error) {
	inputReq := &ESSearchReq{
		SearchType:        "all",
		MapperValueOrPath: json.RawMessage(mapValue),
		SearchBody: []ESSearchBody{
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
	return doTaskProcess[*ESSearchReq, []ESSearchOutput](entity, "search_engine_search", inputReq, 10000)
}

type ESInsertReq struct {
	Corpus            []map[string]any `json:"corpus"`
	Id                string           `json:"id"`
	MapperValueOrPath any              `json:"mapper_config"`
}

func (entity *JobdApi) EsInsert(inputReq ESInsertReq) (res any, err error) {
	entity.Client.MaxReqBodyLen = -1
	return doTaskProcessString[ESInsertReq](entity, "search_engine_insert", inputReq, 50000)
}

type ESDeleteReq struct {
	DocIds            []string `json:"doc_ids"`
	Id                string   `json:"id"`
	MapperValueOrPath any      `json:"mapper_config"`
	DeleteAll         bool     `json:"delete_all"`
}

type ESDeleteRes struct {
	DeleteResults []string `json:"delete_results"`
}

func (entity *JobdApi) EsDelete(inputReq *ESDeleteReq) (res ESDeleteRes, err error) {
	return doTaskProcess[*ESDeleteReq, ESDeleteRes](entity, "search_engine_delete", inputReq, 10000)
}
