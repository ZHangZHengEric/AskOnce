// Package es -----------------------------
// @file      : common.go
// @author    : xiangtao
// @contact   : xiangtao1994@gmail.com
// @time      : 2024/12/23 15:49
// -------------------------------------------
package es

import (
	"askonce/helpers"
	"encoding/json"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/gin-gonic/gin"
	"github.com/xiangtao94/golib/pkg/zlog"
	"golang.org/x/sync/errgroup"
	"sync"
	"time"
)

type Document interface {
	ConvertAny() any
	GetDocId() int64
	SetScore(score float64)
	SetSource(source string)
}

type CommonDocument struct {
	DocId      int64     `json:"doc_id"`
	DocContent string    `json:"doc_content"`
	Emb        []float32 `json:"emb"`
	Score      float64   `json:"score,omitempty"`
	Source     string    `json:"source,omitempty"`
}

func (cd *CommonDocument) ConvertAny() any {
	return cd
}

func (cd *CommonDocument) GetDocId() int64 {
	return cd.DocId
}

func (cd *CommonDocument) SetScore(score float64) {
	cd.Score = score
}

func (cd *CommonDocument) SetSource(source string) {
	cd.Source = source
}

func CommonDocumentSearch[T Document](ctx *gin.Context, indexName string, docContent string, vec []float32, size int) (docs []T, err error) {
	docs = make([]T, 0)
	startTime := time.Now()
	wg, _ := errgroup.WithContext(ctx)
	lock := sync.RWMutex{}
	wg.Go(func() error {
		explain := true
		candidates := 1000
		knnSearch := types.KnnSearch{
			Field:         "emb",
			K:             &size,
			NumCandidates: &candidates,
			QueryVector:   vec,
		}
		req := &search.Request{
			Source_: types.SourceFilter{
				Excludes: []string{"emb"},
			},
			Knn:     []types.KnnSearch{knnSearch},
			Size:    &size,
			Explain: &explain,
		}
		searchRes, err := helpers.EsClient.Search(ctx, indexName, req)
		if err != nil {
			return err
		}
		for _, hit := range searchRes.Hits.Hits {
			var tmpD T
			err = json.Unmarshal(hit.Source_, &tmpD)
			if err != nil {
				return err
			}
			tmpD.SetScore(float64(*hit.Score_))
			tmpD.SetSource("vec")
			docs = append(docs, tmpD)
		}
		lock.Lock()
		defer lock.Unlock()
		return nil
	})
	wg.Go(func() error {
		req := &search.Request{
			Source_: types.SourceFilter{
				Excludes: []string{"emb"},
			},
			Query: &types.Query{
				Match: map[string]types.MatchQuery{
					"doc_content": {
						Query: docContent,
					},
				},
			},
			Size: &size,
		}
		searchRes, err := helpers.EsClient.Search(ctx, indexName, req)
		if err != nil {
			return err
		}
		lock.Lock()
		defer lock.Unlock()
		for _, hit := range searchRes.Hits.Hits {
			var tmpD T
			err = json.Unmarshal(hit.Source_, &tmpD)
			if err != nil {
				return err
			}
			tmpD.SetScore(float64(*hit.Score_))
			tmpD.SetSource("bm25")
			docs = append(docs, tmpD)
		}
		return nil
	})
	if err := wg.Wait(); err != nil {
		return nil, err
	}
	zlog.Infof(ctx, "es搜索耗时 %v", time.Since(startTime).Milliseconds())
	return
}

func CommonBatchInsert[T Document](ctx *gin.Context, indexName string, data []T) (err error) {
	if len(data) == 0 {
		return
	}
	batchSize := 100
	docsChunk := slice.Chunk(data, batchSize)
	wg, _ := errgroup.WithContext(ctx)
	for _, ant := range docsChunk {
		wg.Go(func() error {
			anySlice := []any{}
			for _, d := range ant {
				anySlice = append(anySlice, d.ConvertAny())
			}
			err := helpers.EsClient.DocumentInsert(ctx, indexName, anySlice)
			if err != nil {
				return err
			}
			return nil
		})
	}
	if err = wg.Wait(); err != nil {
		docIds := make([]int64, 0)
		for _, d := range data {
			docIds = append(docIds, d.GetDocId())
		}
		_ = CommonDocumentDelete(ctx, indexName, docIds)
	}
	return
}

func CommonDocumentDelete(ctx *gin.Context, indexName string, docIds []int64) (err error) {
	exist, err := helpers.EsClient.CheckIndex(ctx, indexName)
	if err != nil {
		return err
	}
	if !exist {
		return
	}
	query := &types.Query{
		Terms: &types.TermsQuery{
			TermsQuery: map[string]types.TermsQueryField{
				"doc_id": docIds, // 批量值
			},
		},
	}
	err = helpers.EsClient.DocumentDelete(ctx, indexName, query)
	if err != nil {
		return err
	}
	return nil
}
