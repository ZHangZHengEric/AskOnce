package helpers

import (
	"askonce/conf"
	"fmt"
	"github.com/elastic/elastic-transport-go/v8/elastictransport"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/gin-gonic/gin"
	"net/url"
	"os"
)

var (
	EsClient *ElasticsearchClient
)

func InitElastic() {
	config := conf.WebConf.ElasticSearch
	endpointUrl, err := url.Parse(config.Addr)
	if err != nil {
		panic("init elastic failed, err:" + err.Error())
	}
	cfg := elasticsearch.Config{
		Addresses: []string{endpointUrl.String()},
		Username:  config.Username,
		Password:  config.Password,
		Logger: &elastictransport.JSONLogger{
			Output:             os.Stdout,
			EnableRequestBody:  true,
			EnableResponseBody: true,
		},
	}
	typeClient, err := elasticsearch.NewTypedClient(cfg)
	if err != nil {
		panic("init elastic failed, err:" + err.Error())
	}
	EsClient = &ElasticsearchClient{
		Client: typeClient,
	}
}

type ElasticsearchClient struct {
	Client *elasticsearch.TypedClient
}

// CheckIndex 判断索引是否存在
func (ec *ElasticsearchClient) CheckIndex(ctx *gin.Context, indexName string) (bool, error) {
	res, err := ec.Client.Indices.Exists(indexName).Do(ctx)
	if err != nil {
		return false, err
	}
	return res, nil
}

// CreateIndex 根据提供的 mapping 创建索引
func (ec *ElasticsearchClient) CreateIndex(ctx *gin.Context, indexName string, mapping *types.TypeMapping, setting *types.IndexSettings) error {
	res, err := ec.Client.Indices.Create(indexName).Mappings(mapping).Settings(setting).Do(ctx)
	if err != nil {
		return err
	}
	if !res.Acknowledged {
		return fmt.Errorf("failed to create index: %s", indexName)
	}
	return nil
}

// 删除索引
func (ec *ElasticsearchClient) DeleteIndex(ctx *gin.Context, indexName string) error {
	res, err := ec.Client.Indices.Delete(indexName).IgnoreUnavailable(true).Do(ctx)
	if err != nil {
		return err
	}
	if !res.Acknowledged {
		return fmt.Errorf("failed to delete index: %s", indexName)
	}
	return nil
}

// BulkInsert 批量插入数据，批量限制为 3000 条
func (ec *ElasticsearchClient) DocumentInsert(ctx *gin.Context, indexName string, docs []any) (err error) {

	bulk := ec.Client.Bulk().Index(indexName)
	for _, doc := range docs {
		err = bulk.CreateOp(types.CreateOperation{Index_: &indexName}, doc)
		if err != nil {
			return err
		}
	}
	resp, err := bulk.Do(ctx)
	if err != nil {
		return err
	}
	if resp.Errors {
		return fmt.Errorf("elastic search error: %v", resp.Errors)
	}
	return nil
}

// BulkDelete 批量删除文档
func (ec *ElasticsearchClient) DocumentDelete(ctx *gin.Context, indexName string, query *types.Query) error {
	_, err := ec.Client.DeleteByQuery(indexName).Query(query).Do(ctx)
	if err != nil {
		return err
	}
	return nil
}

// Search 混合查询
func (ec *ElasticsearchClient) Search(ctx *gin.Context, indexName string, query *search.Request) (*search.Response, error) {
	res, err := ec.Client.Search().Index(indexName).Request(query).Do(ctx)
	if err != nil {
		return nil, err
	}
	if res.TimedOut {
		return nil, fmt.Errorf("knn search time out")
	}
	return res, nil
}
