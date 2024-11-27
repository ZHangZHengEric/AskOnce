package dto_kdb

// 输入数据源
type InfoList struct {
	DataSourceType DataSourceType `json:"dataSourceType"`
	FileIds        []string       `json:"fileIds"`
}

// 文档处理规则
type ProcessRule struct {
	Mode  string               `json:"mode"` // custom , auto
	Rules DocProcessRuleDetail `json:"rules"`
}

type DocProcessRuleDetail struct {
	PreProcessRules []PreProcessRule `json:"preProcessRules"` // 文本预处理规则
	Segmentation    SegmentationRule `json:"segmentation"`    // 分段规则
}

type PreProcessRule struct {
	Id      string `json:"id"`
	Enabled bool   `json:"enabled"`
}

type SegmentationRule struct {
	SegmentMode  DocSegmentMode `json:"segmentMode"`  // 分段模式 “simple" ，"structured"
	Separator    string         `json:"separator"`    // 分段标识符
	ChunkSize    int            `json:"chunkSize"`    // 分段最大长度
	ChunkOverlap int            `json:"chunkOverlap"` // 分段重叠长度
}

type DataSourceType string

const (
	DataSourceTypeFile = "file"
)

// 文档分段模式
type DocSegmentMode string

const (
	DocSegmentModeSimple     DocSegmentMode = "simple"     // 简单
	DocSegmentModeStructured DocSegmentMode = "structured" // 结构化
)
