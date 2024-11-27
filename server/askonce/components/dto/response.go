package dto

type FileUploadRes struct {
	FileId    string        `json:"fileId"`    // 文件id
	Name      string        `json:"name"`      // 文件名
	Extension string        `json:"extension"` // 文件后缀
	Size      FileSizeLabel `json:"size"`
}

type ModelProviderRes struct {
	Data []ModelProviderItem `json:"data"`
}

type ModelProviderItem struct {
	Provider  string               `json:"provider"`
	Label     I18n                 `json:"label"`
	IconSmall I18n                 `json:"iconSmall"`
	IconLarge I18n                 `json:"iconLarge"`
	Status    string               `json:"status"`
	Models    []ModelProviderModel `json:"models"`
}

type ModelProviderModel struct {
	Model           string      `json:"model"`
	Label           I18n        `json:"label"`
	ModelType       string      `json:"modelType"`
	Features        interface{} `json:"features"`
	FetchFrom       string      `json:"fetchFrom"`
	ModelProperties struct {
		ContextSize int `json:"contextSize"`
		MaxChunks   int `json:"maxChunks"`
	} `json:"modelProperties"`
	Deprecated           bool   `json:"deprecated"`
	Status               string `json:"status"`
	LoadBalancingEnabled bool   `json:"loadBalancingEnabled"`
}

type I18n struct {
	ZhHans string `json:"zh_Hans"`
	EnUS   string `json:"en_US"`
}
