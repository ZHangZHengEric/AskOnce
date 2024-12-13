// Package es -----------------------------
// @file      : client.go
// @author    : xiangtao
// @contact   : xiangtao1994@gmail.com
// @time      : 2024/12/13 00:51
// -------------------------------------------
package es

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/knnsearch"
	"github.com/gin-gonic/gin"
)

type ElasticsearchClient struct {
	Client *elasticsearch.Client
}

type ElasticsearchModel interface {
	GetIndexName() string // 获取索引名称
	GetMapping() Mappings // 获取索引 Mapping
}

// Mappings 定义 properties 的配置
type Mappings struct {
	Properties map[string]Field `json:"properties"`
}

// Field 定义单个字段的配置
type Field struct {
	Type       string  `json:"type"`
	Index      *string `json:"index,omitempty"`      // 可选字段
	Analyzer   *string `json:"analyzer,omitempty"`   // 可选字段
	Similarity *string `json:"similarity,omitempty"` // 可选字段
	Dims       *int    `json:"dims,omitempty"`       // 可选字段，仅 dense_vector 使用
}

// CreateIndex 根据提供的 mapping 创建索引
func (ec *ElasticsearchClient) CreateIndex(ctx *gin.Context, indexName string, mapping Mappings) error {
	mappingMap := map[string]interface{}{
		"mappings": mapping,
	}
	// 将结构体转换为 JSON
	body, err := json.Marshal(mappingMap)
	if err != nil {
		return fmt.Errorf("failed to marshal index mapping: %v", err)
	}
	res, err := ec.Client.Indices.Create(indexName,
		ec.Client.Indices.Create.WithContext(ctx),
		ec.Client.Indices.Create.WithBody(bytes.NewReader(body)))
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.IsError() {
		return fmt.Errorf("failed to create index: %s", res.String())
	}

	return nil
}

// 删除索引
func (ec *ElasticsearchClient) DeleteIndex(ctx *gin.Context, indexName string) error {
	res, err := ec.Client.Indices.Delete(
		[]string{indexName},
		ec.Client.Indices.Delete.WithContext(ctx),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.IsError() {
		return fmt.Errorf("failed to delete index: %s", res.String())
	}
	return nil
}

// BulkInsert 批量插入数据，批量限制为 3000 条
func (ec *ElasticsearchClient) BulkInsert(ctx *gin.Context, indexName string, docs []map[string]any) error {
	batchSize := 3000
	for i := 0; i < len(docs); i += batchSize {
		end := i + batchSize
		if end > len(docs) {
			end = len(docs)
		}
		batch := docs[i:end]

		var buf bytes.Buffer
		for _, doc := range batch {
			docJSON, _ := json.Marshal(doc)
			buf.Write(docJSON)
			buf.WriteByte('\n')
		}

		res, err := ec.Client.Bulk(&buf,
			ec.Client.Bulk.WithIndex(indexName),
			ec.Client.Bulk.WithContext(ctx),
		)
		if err != nil {
			return err
		}
		defer res.Body.Close()
		if res.IsError() {
			return fmt.Errorf("bulk insert error: %s", res.String())
		}
	}
	return nil
}

// BulkDelete 批量删除文档
func (ec *ElasticsearchClient) BulkDelete(ctx *gin.Context, indexName string, ids []string) error {
	var buf bytes.Buffer
	for _, id := range ids {
		meta := map[string]interface{}{
			"delete": map[string]interface{}{
				"_id": id,
			},
		}
		metaJSON, _ := json.Marshal(meta)
		buf.Write(metaJSON)
		buf.WriteByte('\n')
	}

	res, err := ec.Client.Bulk(&buf,
		ec.Client.Bulk.WithIndex(indexName),
		ec.Client.Bulk.WithContext(ctx))
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.IsError() {
		return fmt.Errorf("bulk delete error: %s", res.String())
	}
	return nil
}

type ESSearchSource struct {
	Excludes []string `json:"excludes"`
}
type ESSearchBody struct {
	Source ESSearchSource  `json:"_source,omitempty"`
	Knn    *SearchBodyKnn  `json:"knn,omitempty"`
	Query  *SearchBodyBm25 `json:"query,omitempty"`
	Size   *int            `json:"size,omitempty"`
}

type SearchBodyBm25 struct {
	Match map[string]any `json:"match"`
}

type SearchBodyKnn struct {
	Field         string `json:"field"`
	QueryVector   any    `json:"query_vector"`
	K             int    `json:"k"`
	NumCandidates int    `json:"num_candidates"`
}

type ESSearchOutput struct {
	Source SearchOutputSource `json:"source"`
	Score  float32            `json:"score"`
}

type SearchOutputSource struct {
	DocId       int64  `json:"doc_id"`
	DocContent  string `json:"doc_content"`
	DataSplitId int64  `json:"data_split_id"`
	Start       int    `json:"start"`
	End         int    `json:"end"`
}

// Search 混合查询支持 BM25 和向量搜索
func (ec *ElasticsearchClient) Search(ctx *gin.Context, indexName string, query *knnsearch.Request) ([]ESSearchOutput, error) {
	// 将请求体转为 JSON
	body, _ := json.Marshal(query)
	res, err := ec.Client.KnnSearch(
		[]string{indexName},
		ec.Client.KnnSearch.WithBody(bytes.NewReader(body)),
		ec.Client.KnnSearch.WithContext(ctx),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("search error: %s", res.String())
	}
	render := knnsearch.NewResponse()
	if err := json.NewDecoder(res.Body).Decode(&render); err != nil {
		return nil, err
	}
	output := make([]ESSearchOutput, 0, len(render.Hits.Hits))
	for _, hit := range render.Hits.Hits {
		tmp := SearchOutputSource{}
		b, _ := hit.Source_.MarshalJSON()
		err = json.Unmarshal(b, &tmp)
		output = append(output, ESSearchOutput{
			Source: tmp,
			Score:  float32(*hit.Score_),
		})
	}
	return output, nil
}
