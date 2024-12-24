package dto_search

type CommonSearchOutput struct {
	// 标题
	Title string `json:"title"`
	// 文档路径或页面地址
	Url string `json:"url"`
	// 检索的内容
	Content string `json:"content"`
	// 元数据
	Metadata map[string]interface{} `json:"metadata,omitempty"`

	// 来源
	Form string `json:"form,omitempty"`
	// 文档id
	DocId int64 `json:"docId,omitempty"`
	// 文档切片id
	DocSegmentId int64 `json:"docSegmentId,omitempty"`
	// 文档全文
	FullContent string `json:"fullContent,omitempty"`
	// 分数
	Score float64 `json:"score,omitempty"`
	// 数据库相关信息
	DatabaseInfo *CommonSearchDatabaseInfo `json:"databaseInfo,omitempty"`
}

type CommonSearchDatabaseInfo struct {
	DatabaseName    string `json:"databaseName"`
	DatabaseComment string `json:"databaseComment"`
	TableName       string `json:"tableName"`
	TableComment    string `json:"tableComment"`
	ColumnName      string `json:"columnName"`
	ColumnType      string `json:"columnType"`
	ColumnComment   string `json:"columnComment"`
	Sql             string `json:"sql"`
}

type DoReferItem struct {
	Start       int                `json:"start"` // 答案开始下标
	End         int                `json:"end"`   // 答案结束下标
	NumberIndex int                `json:"numberIndex"`
	Refers      []DoReferReferItem `json:"refers"` // 参考文档信息
}

type DoReferReferItem struct {
	Index      int `json:"index"`      // 参考文档下标
	ReferStart int `json:"referStart"` // 参考文档content文字开始下标
	ReferEnd   int `json:"referEnd"`   // 参考文档content文字结束下标
}
