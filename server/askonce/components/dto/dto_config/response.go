package dto_config

type ConfigResp struct {
	Language  string `json:"language"`
	ModelType string `json:"modelType"`
}

type Dict struct {
	Name   string `json:"name"`
	EnName string `json:"enName"`
	Value  string `json:"value"`
}
