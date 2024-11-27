package dto

type PageParam struct {
	PageNo   int `json:"pageNo" binding:"required"`
	PageSize int `json:"pageSize" binding:"required"`
}

type EmptyReq struct {
}

type LoginInfo struct {
	UserId   string `json:"UserId" redis:"UserId"`
	UserName string `json:"UserName" redis:"UserName"`
}

// 知识库设置
type KdbSetting struct {
	EmbeddingModel DocEmbeddingModel `json:"embeddingModel"` // Embedding 模型
	RetrievalModel RetrievalSetting  `json:"retrievalModel"` // 召回设置
}

// 文档embedding模式
type DocEmbeddingModel string

const (
	DocEmbeddingModelCommon DocEmbeddingModel = "common"
)

// 检索模型
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

type RetrievalSettingWeights struct {
	KeywordWeight float32 `json:"keywordWeight"` //关键词权重
	VectorWeight  float32 `json:"vectorWeight"`  //向量权重
}
