package dto

type PageParam struct {
	PageNo   int `json:"pageNo" binding:"required"`
	PageSize int `json:"pageSize" binding:"required"`
}

type EmptyReq struct {
}

type LoginInfoSession struct {
	UserId     string `json:"UserId" redis:"UserId"`
	Account    string `json:"Account" redis:"Account"`
	Session    string `json:"Session" redis:"Session"`
	LoginTime  int64  `json:"LoginTime" redis:"LoginTime"`
	ExpireTime int64  `json:"ExpireTime" redis:"ExpireTime"`
}

// 知识库设置
type KdbSetting struct {
	RetrievalModel     RetrievalSetting `json:"retrievalModel"`     // 召回设置
	ReferenceThreshold float32          `json:"referenceThreshold"` // 引用参考阈值
	KdbAttach          KdbAttach        `json:"kdbAttach"`          // 附加属性
}

// 检索设置
type RetrievalSetting struct {
	SearchMethod          DocSearchMethod         `json:"searchMethod"`          // 搜索方法 “keyword", "vector", "all"
	TopK                  int                     `json:"topK"`                  // 召回多少条
	ScoreThresholdEnabled bool                    `json:"scoreThresholdEnabled"` // 搜索分数阈值开启
	ScoreThreshold        float32                 `json:"scoreThreshold"`        // 搜索分数阈值
	Weights               RetrievalSettingWeights `json:"weights"`               // 混合搜索 搜索权重
}

// 文档搜索方法
type DocSearchMethod string

const (
	DocSearchMethodKeyWord DocSearchMethod = "keyword"
	DocSearchMethodVector  DocSearchMethod = "vector"
	DocSearchMethodAll     DocSearchMethod = "all"
)

type KdbAttach struct {
	Language   string   `json:"language"`
	CoverId    int64    `json:"coverId"`
	CoverColor bool     `json:"coverColor"`
	Cases      []string `json:"cases"`
}

type RetrievalSettingWeights struct {
	KeywordWeight float32 `json:"keywordWeight"` //关键词权重
	VectorWeight  float32 `json:"vectorWeight"`  //向量权重
}

// 文档处理规则
type DocProcessSetting struct {
	Mode         string           `json:"mode"`         // custom , auto
	Segmentation SegmentationRule `json:"segmentation"` // 分段规则
}

type SegmentationRule struct {
	Separator    string `json:"separator"`    // 分段标识符
	ChunkSize    int    `json:"chunkSize"`    // 分段最大长度
	ChunkOverlap int    `json:"chunkOverlap"` // 分段重叠长度
}
