package helpers

import (
	"askonce/conf"
	"context"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/gin-gonic/gin"
	"github.com/xiangtao94/golib/pkg/zlog"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	_defaultPrintRequestLen  = 512
	_defaultPrintResponseLen = 10240
)

const (
	ES_LOG_MAX_REQ_LEN  = "ES_LOG_MAX_REQ_LEN"
	ES_LOG_MAX_RESP_LEN = "ES_LOG_MAX_RESP_LEN"
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
		Logger:    newLogger(),
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
	Client        *elasticsearch.TypedClient
	MaxReqBodyLen int
	// response body 最大长度展示，0表示采用默认的10240，-1表示不打印。指定长度的时候需注意，返回的json可能被截断
	MaxRespBodyLen int
}

// CheckIndex 判断索引是否存在
func (ec *ElasticsearchClient) CheckIndex(ctx *gin.Context, indexName string) (bool, error) {
	ec.appendContext(ctx)
	res, err := ec.Client.Indices.Exists(indexName).Do(ctx)
	if err != nil {
		return false, err
	}
	return res, nil
}

// CreateIndex 根据提供的 mapping 创建索引
func (ec *ElasticsearchClient) CreateIndex(ctx *gin.Context, indexName string, mapping *types.TypeMapping, setting *types.IndexSettings) error {
	ec.appendContext(ctx)
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

	ec.appendContext(ctx)
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
	ec.appendContext(ctx)
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
	ec.appendContext(ctx)
	_, err := ec.Client.DeleteByQuery(indexName).Query(query).Do(ctx)
	if err != nil {
		return err
	}
	return nil
}

// Search 混合查询
func (ec *ElasticsearchClient) Search(ctx *gin.Context, indexName string, query *search.Request) (*search.Response, error) {
	ec.appendContext(ctx)
	res, err := ec.Client.Search().Index(indexName).Request(query).Do(ctx)
	if err != nil {
		return nil, err
	}
	if res.TimedOut {
		return nil, fmt.Errorf("knn search time out")
	}
	return res, nil
}

type elasticLogger struct {
	logger *zlog.Logger
}

func (e *elasticLogger) LogRoundTrip(request *http.Request, response *http.Response, err error, start time.Time, duration time.Duration) error {
	request.Context()
	fields := e.AppendCustomField(request.Context())
	fields = append(fields,
		zlog.String("path", request.URL.Path),
		zlog.String("method", request.Method),
		zlog.String("query", request.URL.RawQuery),
	)
	reqBody, _ := io.ReadAll(request.Body)
	defer request.Body.Close()
	respBody, _ := io.ReadAll(response.Body)
	defer response.Body.Close()
	requestData, respData := formatLogMsg(request.Context(), reqBody, respBody)
	fields = append(fields, zlog.String("requestParam", string(requestData)))
	fields = append(fields, zlog.Int("responseStatus", response.StatusCode))
	fields = append(fields, zlog.String("response", string(respData)))
	fields = append(fields, zlog.AppendCostTime(start, time.Now())...)
	msg := "success"
	if err != nil {
		msg = err.Error()
		zlog.ErrorLogger(covertGinContext(request.Context()), msg, fields...)
		return nil
	}
	zlog.InfoLogger(covertGinContext(request.Context()), msg, fields...)
	return nil
}

func (e *elasticLogger) RequestBodyEnabled() bool {
	return true
}

func (e *elasticLogger) ResponseBodyEnabled() bool {
	return true
}

func newLogger() *elasticLogger {
	return &elasticLogger{
		logger: zlog.GetZapLogger(),
	}
}

func covertGinContext(ctx context.Context) *gin.Context {
	if c, ok := ctx.(*gin.Context); ok && c != nil {
		return c
	}
	return nil
}

func (e *elasticLogger) AppendCustomField(ctx context.Context) []zlog.Field {
	var requestID string
	ctx1 := covertGinContext(ctx)
	if ctx1 != nil {
		requestID, _ = ctx1.Value(zlog.ContextKeyRequestID).(string)
	}
	fields := []zlog.Field{
		zlog.String("requestId", requestID),
	}
	return fields
}

func formatLogMsg(context context.Context, requestParam, responseData []byte) (req, resp []byte) {
	ctx := covertGinContext(context)
	if ctx == nil {
		return requestParam, responseData
	}
	maxReqBodyLen := ctx.GetInt(ES_LOG_MAX_REQ_LEN)
	if maxReqBodyLen == 0 {
		maxReqBodyLen = _defaultPrintRequestLen
	}

	maxRespBodyLen := ctx.GetInt(ES_LOG_MAX_RESP_LEN)
	if maxRespBodyLen == 0 {
		maxRespBodyLen = _defaultPrintResponseLen
	}

	if maxReqBodyLen != -1 {
		req = requestParam
		if len(requestParam) > maxReqBodyLen {
			req = req[:maxReqBodyLen]
		}
	}

	if maxRespBodyLen != -1 {
		resp = responseData
		if len(responseData) > maxRespBodyLen {
			resp = resp[:maxRespBodyLen]
		}
	}
	return req, resp
}

func (ec *ElasticsearchClient) appendContext(ctx *gin.Context) {
	ctx.Set(ES_LOG_MAX_REQ_LEN, ec.MaxReqBodyLen)
	ctx.Set(ES_LOG_MAX_RESP_LEN, ec.MaxRespBodyLen)
}
