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
}

type DoReferItem struct {
	Start       int                `json:"start"`
	End         int                `json:"end"`
	NumberIndex int                `json:"numberIndex"`
	Refers      []DoReferReferItem `json:"refers"`
}

type DoReferReferItem struct {
	Index      int `json:"index"`
	ReferStart int `json:"referStart"`
	ReferEnd   int `json:"referEnd"`
}
