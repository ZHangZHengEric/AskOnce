package dto_kdb

// 输入数据源
type InfoList struct {
	DataSourceType DataSourceType `json:"dataSourceType"`
	FileIds        []string       `json:"fileIds"`
}

type DataSourceType string

const (
	DataSourceTypeFile = "file"
)
