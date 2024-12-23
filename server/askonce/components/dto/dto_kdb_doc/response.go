package dto_kdb_doc

type AddRes struct {
	KdbDataId int64 `json:"kdbDataId"`
}

type RecallRes struct {
	List []RecallItem `json:"list"`
}

type RecallItem struct {
	DataName      string `json:"dataName"`
	DataPath      string `json:"dataPath"`
	SearchContent string `json:"searchContent"`
	DataContent   string `json:"dataContent"`
}

type AddZipRes struct {
	TaskId string `json:"taskId" form:"taskId"`
}

type InfoRes struct {
	KdbDataId int64                  `json:"kdbDataId"`
	DataName  string                 `json:"dataName"`
	DbType    string                 `json:"dbType,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

type ListResp struct {
	List  []ListItem `json:"list"`
	Total int64      `json:"total"`
}

type ListItem struct {
	Id         int64  `json:"id"`
	Type       string `json:"type"`
	DataName   string `json:"dataName"`
	DataPath   string `json:"dataPath,omitempty"`
	DataSuffix string `json:"dataSuffix,omitempty"`

	DbType string `json:"dbType,omitempty"`

	Status     int    `json:"status"` // 0正在构建到知识库 1 成功 2 失败
	CreateTime string `json:"createTime"`
}
