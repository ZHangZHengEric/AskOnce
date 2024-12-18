package data

import (
	"askonce/api/jobd"
	"askonce/helpers"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/xiangtao94/golib/flow"
	"golang.org/x/sync/errgroup"
	"sync"
)

/*
*
文档处理
*/
type DocumentData struct {
	flow.Layer
	jobdApi *jobd.JobdApi
}

func (d *DocumentData) OnCreate() {
	d.jobdApi = flow.Create(d.GetCtx(), new(jobd.JobdApi))
}

// 文本切分
func (d *DocumentData) TextSplit(content string) (segments []jobd.TextChunkItem, err error) {
	documentSplitRes, err := d.jobdApi.DocumentSplit(content)
	if err != nil {
		return nil, err
	}
	segments = documentSplitRes.SentencesList
	return segments, nil
}

// 批量文本转向量
func (d *DocumentData) TextEmbedding(texts []string) (embResAll [][]float32, err error) {
	// 最大批次
	sentsG := slice.Chunk(texts, 30)
	embResAll = make([][]float32, 0)
	lock := sync.Mutex{}
	embResMap := make(map[int][][]float32)
	eg2, _ := errgroup.WithContext(d.GetCtx())
	for i, ss := range sentsG {
		tmp := ss
		index := i
		eg2.Go(func() error {
			embRes, errA := helpers.EmbeddingGpt.CreateEmbedding(d.GetCtx(), tmp)
			if errA != nil {
				return errA
			}
			lock.Lock()
			embResMap[index] = embRes
			lock.Unlock()
			return nil
		})
	}
	if err := eg2.Wait(); err != nil {
		return nil, err
	}
	for i := range sentsG {
		embResAll = append(embResAll, embResMap[i]...)
	}
	return
}
