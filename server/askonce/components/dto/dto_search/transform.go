package dto_search

type CommonSearchOutput struct {
	Title   string `json:"title"`
	Url     string `json:"url"`
	Content string `json:"content"`
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
