package data

import (
	"askonce/api/jobd"
	"askonce/conf"
	"askonce/helpers"
	"askonce/models"
	"encoding/json"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/xiangtao94/golib/flow"
	"github.com/xiangtao94/golib/pkg/orm"
	"golang.org/x/sync/errgroup"
	"strconv"
	"strings"
	"time"
)

type KdbDocData struct {
	flow.Redis
	kdbDocDao        *models.KdbDocDao
	kdbDocContentDao *models.KdbDocContentDao
	kdbDocSegmentDao *models.KdbDocSegmentDao
	jobdApi          *jobd.JobdApi
}

func (k *KdbDocData) OnCreate() {
	k.kdbDocDao = flow.Create(k.GetCtx(), new(models.KdbDocDao))
	k.kdbDocContentDao = flow.Create(k.GetCtx(), new(models.KdbDocContentDao))
	k.kdbDocSegmentDao = flow.Create(k.GetCtx(), new(models.KdbDocSegmentDao))
	k.jobdApi = flow.Create(k.GetCtx(), new(jobd.JobdApi))
}

func (k *KdbDocData) DeleteDoc(kdb *models.Kdb, docId int64) (err error) {
	doc, err := k.kdbDocDao.GetById(docId)
	if err != nil {
		return err
	}
	if doc == nil {
		return
	}
	db := helpers.MysqlClient.WithContext(k.GetCtx())
	k.kdbDocDao.SetDB(db)
	k.kdbDocContentDao.SetDB(db)
	k.kdbDocSegmentDao.SetDB(db)
	tx := db.Begin()
	err = k.kdbDocDao.DeleteById(docId)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = k.kdbDocContentDao.DeleteByDocIds([]int64{docId})
	if err != nil {
		tx.Rollback()
		return err
	}
	err = k.kdbDocSegmentDao.DeleteByDocIds(kdb.Id, []int64{docId})
	if err != nil {
		tx.Rollback()
		return err
	}
	if doc.Status == models.KdbDocSuccess {
		esDbConfigStr := strings.Replace(conf.WebConf.EsDbConfig, "@indexName@", kdb.GetIndexName(), 1)
		_, err = k.jobdApi.EsDelete(&jobd.ESDeleteReq{DocIds: []string{
			strconv.FormatInt(docId, 10)},
			MapperValueOrPath: json.RawMessage(esDbConfigStr),
		})
		if err != nil {
			tx.Rollback()
			return
		}
	}
	tx.Commit()
	return
}

func (k *KdbDocData) SaveDocBuild(kdb *models.Kdb, doc *models.KdbDoc, content string, splitList []jobd.TextChunkItem, embeddingAll [][]float32) (err error) {
	segments := make([]*models.KdbDocSegment, 0, len(splitList))
	esInsertCorpus := make([]map[string]any, 0)
	for i, split := range splitList {
		if len(split.PassageContent) == 0 {
			continue
		}
		segments = append(segments, &models.KdbDocSegment{
			DocId:      doc.Id,
			KdbId:      doc.KdbId,
			StartIndex: split.Start,
			EndIndex:   split.End,
			Text:       split.PassageContent,
			CrudModel: orm.CrudModel{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		})
		esInsertCorpus = append(esInsertCorpus, map[string]any{
			"doc_id":      doc.Id,
			"doc_content": split.PassageContent,
			"start":       split.Start,
			"end":         split.End,
			"emb":         embeddingAll[i],
		})
	}
	db := helpers.MysqlClient.WithContext(k.GetCtx())
	k.kdbDocDao.SetDB(db)
	k.kdbDocContentDao.SetDB(db)
	k.kdbDocSegmentDao.SetDB(db)
	tx := db.Begin()
	err = k.kdbDocContentDao.Insert(&models.KdbDocContent{
		DocId:   doc.Id,
		KdbId:   doc.KdbId,
		Content: content,
	})
	if err != nil {
		tx.Rollback()
		k.LogErrorf("kdbDocContentDao insert error，docId %v, error %v", doc.Id, err.Error())
		return err
	}
	err = k.kdbDocSegmentDao.BatchInsert(doc.Id, segments)
	if err != nil {
		tx.Rollback()
		k.LogErrorf("kdbDocSegmentDao insert error，docId %v, error %v", doc.Id, err.Error())
		return err
	}
	err = k.saveEs(kdb, doc, esInsertCorpus, len(embeddingAll[0]))
	if err != nil {
		tx.Rollback()
		k.LogErrorf("saveEs error，docId %v, error %v", doc.Id, err.Error())
		return err
	}
	tx.Commit()
	return
}

func (k *KdbDocData) saveEs(kdb *models.Kdb, doc *models.KdbDoc, corpus []map[string]any, embLength int) (err error) {
	corpusG := slice.Chunk(corpus, 500)
	eg, _ := errgroup.WithContext(k.GetCtx())
	esDbConfigStr := strings.ReplaceAll(conf.WebConf.EsDbConfig, "@indexName@", kdb.GetIndexName())
	esDbConfigStr = strings.ReplaceAll(esDbConfigStr, "@dimsLength@", strconv.Itoa(embLength))
	for _, ccc := range corpusG {
		tmp := ccc
		eg.Go(func() error {
			_, err := k.jobdApi.EsInsert(jobd.ESInsertReq{
				Corpus:            tmp,
				MapperValueOrPath: json.RawMessage(esDbConfigStr),
			})
			if err != nil {
				return err
			}
			return nil
		})

	}
	if err := eg.Wait(); err != nil {
		_, _ = k.jobdApi.EsDelete(&jobd.ESDeleteReq{
			DocIds:            []string{strconv.FormatInt(doc.Id, 10)},
			MapperValueOrPath: json.RawMessage(esDbConfigStr),
		})
		return err
	}
	return
}
