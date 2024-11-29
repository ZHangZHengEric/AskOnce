package dto_history

type AskRes struct {
	List  []AskItem `json:"list"`
	Total int64     `json:"total"`
}

type AskItem struct {
	SessionId  string `json:"sessionId"`
	CreateTime string `json:"createTime"`
	KdbName    string `json:"kdbName"`
	KdbId      int64  `json:"kdbId"`
	Question   string `json:"question"`
	AskType    string `json:"askType"`
}
