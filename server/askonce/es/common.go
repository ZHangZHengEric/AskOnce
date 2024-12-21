// Package es -----------------------------
// @file      : common_doc.go
// @author    : xiangtao
// @contact   : xiangtao1994@gmail.com
// @time      : 2024/12/13 19:05
// -------------------------------------------
package es

import (
	"askonce/helpers"
	"encoding/json"
	"fmt"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/gin-gonic/gin"
	"github.com/xiangtao94/golib/pkg/env"
	"github.com/xiangtao94/golib/pkg/zlog"
	"golang.org/x/sync/errgroup"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type CommonDocument struct {
	DocId        int64     `json:"doc_id"`
	DocContent   string    `json:"doc_content"`
	DocSegmentId int64     `json:"doc_segment_id"`
	Start        int       `json:"start"`
	End          int       `json:"end"`
	Emb          []float32 `json:"emb"`
	Score        float64   `json:"score,omitempty"`
	Source       string    `json:"source,omitempty"`
}

func CommonIndexCreate(ctx *gin.Context, indexName string) (err error) {
	envPath := filepath.Join(env.GetConfDirPath(), "mount/es")
	mappingPath := filepath.Join(envPath, "common_mapping.json")
	mappingFile, err := os.ReadFile(mappingPath)
	if err != nil {
		return fmt.Errorf("read common mapping file failed: %v", err)
	}
	mapping := &types.TypeMapping{}
	err = mapping.UnmarshalJSON(mappingFile)
	if err != nil {
		return
	}
	settingsPath := filepath.Join(envPath, "common_setting.json")
	setting := &types.IndexSettings{}
	settingsFile, err := os.ReadFile(settingsPath)
	if err != nil {
		return fmt.Errorf("read common mapping file failed: %v", err)
	}
	err = setting.UnmarshalJSON(settingsFile)
	if err != nil {
		return
	}
	err = helpers.EsClient.CreateIndex(ctx, indexName, mapping, setting)
	if err != nil {
		return
	}
	return nil
}

func CommonIndexDelete(ctx *gin.Context, indexName string) (err error) {
	err = helpers.EsClient.DeleteIndex(ctx, indexName)
	if err != nil {
		return
	}
	return nil
}

func CommonDocumentInsert(ctx *gin.Context, indexName string, documents []*CommonDocument) (err error) {
	exist, err := helpers.EsClient.CheckIndex(ctx, indexName)
	if err != nil {
		return err
	}
	if !exist {
		err = CommonIndexCreate(ctx, indexName)
		if err != nil {
			return err
		}
	}
	batchSize := 100
	docsChunk := slice.Chunk(documents, batchSize)
	wg, _ := errgroup.WithContext(ctx)
	for _, docT := range docsChunk {
		wg.Go(func() error {
			// 转换为 []any
			var anySlice []any
			for _, v := range docT {
				anySlice = append(anySlice, v)
			}
			err := helpers.EsClient.DocumentInsert(ctx, indexName, anySlice)
			if err != nil {
				return err
			}
			return nil
		})
	}
	if err := wg.Wait(); err != nil {
		docIds := make([]int64, 0)
		for _, doc := range documents {
			docIds = append(docIds, doc.DocId)
		}
		_ = CommonDocumentDelete(ctx, indexName, docIds)
		return err
	}
	return nil
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

func CommonDocumentSearch(ctx *gin.Context, indexName string, docContent string, vec []float32, size int) (docs []*CommonDocument, err error) {
	docs = make([]*CommonDocument, 0)
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
			tmpD := &CommonDocument{}
			err = json.Unmarshal(hit.Source_, &tmpD)
			if err != nil {
				return err
			}
			tmpD.Score = float64(*hit.Score_)
			tmpD.Source = "vec"
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
			tmpD := &CommonDocument{}
			err = json.Unmarshal(hit.Source_, &tmpD)
			if err != nil {
				return err
			}
			tmpD.Source = "bm25"
			tmpD.Score = float64(*hit.Score_)
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
