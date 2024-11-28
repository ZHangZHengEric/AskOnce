package dto_kdb_doc

type ListResp struct {
	List  []ListItem `json:"list"`
	Total int64      `json:"total"`
}

type ListItem struct {
	Id         int64  `json:"id"`
	Type       string `json:"type"`
	DataSuffix string `json:"dataSuffix"`
	DataName   string `json:"dataName"`
	DataPath   string `json:"dataPath"`
	Status     int    `json:"status"` // 0正在构建到知识库 1 成功 2 失败
	CreateTime string `json:"createTime"`
}

type SearchRes struct {
	List []SearchItem `json:"list"`
}

type SearchItem struct {
	DataName      string `json:"dataName"`
	DataPath      string `json:"dataPath"`
	SearchContent string `json:"searchContent"`
	DataContent   string `json:"dataContent"`
}
