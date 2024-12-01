package dto_config

type DetailReq struct {
}

type SaveReq struct {
	Language  string `json:"language"`
	ModelType string `json:"modelType"`
}
