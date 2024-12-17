package data

import (
	"askonce/api/jobd"
	"askonce/components/dto"
	"askonce/es"
	"askonce/helpers"
	"askonce/models"
	"github.com/xiangtao94/golib/flow"
	"github.com/xiangtao94/golib/pkg/orm"
	"time"
)

type KdbDocData struct {
	flow.Redis
	kdbDocDao        *models.KdbDocDao
	fileData         *FileData
	kdbDocContentDao *models.KdbDocContentDao
	kdbDocSegmentDao *models.KdbDocSegmentDao
	jobdApi          *jobd.JobdApi
}

func (k *KdbDocData) OnCreate() {
	k.kdbDocDao = flow.Create(k.GetCtx(), new(models.KdbDocDao))
	k.kdbDocContentDao = flow.Create(k.GetCtx(), new(models.KdbDocContentDao))
	k.kdbDocSegmentDao = flow.Create(k.GetCtx(), new(models.KdbDocSegmentDao))
	k.jobdApi = flow.Create(k.GetCtx(), new(jobd.JobdApi))
	k.fileData = flow.Create(k.GetCtx(), new(FileData))
}

func (k *KdbDocData) DeleteDocs(kdb *models.Kdb, docIds []int64, deleteAll bool) (err error) {
	var docs []*models.KdbDoc
	if deleteAll {
		docs, err = k.kdbDocDao.GetByKdbId(kdb.Id)
		if err != nil {
			return err
		}
	} else {
		docs, err = k.kdbDocDao.GetByIds(docIds)
		if err != nil {
			return err
		}
	}
	if len(docs) == 0 {
		return
	}
	db := helpers.MysqlClient.WithContext(k.GetCtx())
	tx := db.Begin()
	k.kdbDocDao.SetDB(tx)
	k.kdbDocContentDao.SetDB(tx)
	k.kdbDocSegmentDao.SetDB(tx)
	docSuccessIds := make([]int64, 0)
	docFileIds := make([]string, 0)
	for _, doc := range docs {
		if doc.Status == models.KdbDocSuccess {
			docSuccessIds = append(docSuccessIds, doc.Id)
		}
		docFileIds = append(docFileIds, doc.SourceId)
	}
	err = k.fileData.DeleteByFileIds(docFileIds)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = k.kdbDocDao.DeleteByIds(docIds)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = k.kdbDocContentDao.DeleteByDocIds(docIds)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = k.kdbDocSegmentDao.DeleteByDocIds(kdb.Id, docIds)
	if err != nil {
		tx.Rollback()
		return err
	}
	if deleteAll {
		err = es.CommonIndexDelete(k.GetCtx(), kdb.GetIndexName())
		if err != nil {
			tx.Rollback()
			return
		}
	} else {
		err = es.CommonDocumentDelete(k.GetCtx(), kdb.GetIndexName(), docSuccessIds)
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
	for _, split := range splitList {
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
	}
	db := helpers.MysqlClient.WithContext(k.GetCtx())
	tx := db.Begin()

	k.kdbDocDao.SetDB(tx)
	k.kdbDocContentDao.SetDB(tx)
	k.kdbDocSegmentDao.SetDB(tx)
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
	esInsertCorpus := make([]*es.CommonDocument, 0)
	for i, s := range segments {
		esInsertCorpus = append(esInsertCorpus, &es.CommonDocument{
			DocId:        s.DocId,
			DocContent:   s.Text,
			DocSegmentId: s.Id,
			Start:        s.StartIndex,
			End:          s.EndIndex,
			Emb:          embeddingAll[i],
		})
	}
	err = es.CommonDocumentInsert(k.GetCtx(), kdb.GetIndexName(), esInsertCorpus)
	if err != nil {
		tx.Rollback()
		k.LogErrorf("saveEs error，docId %v, error %v", doc.Id, err.Error())
		return err
	}
	tx.Commit()
	return
}

func (k *KdbDocData) GetList(kdbId int64, queryName string, queryStatus []int, param dto.PageParam) (list []*models.KdbDoc, cnt int64, err error) {
	return k.kdbDocDao.GetList(kdbId, queryName, queryStatus, param)
}

func (k *KdbDocData) AddDoc(doc *models.KdbDoc) (err error) {
	err = k.kdbDocDao.Insert(doc)
	if err != nil {
		return err
	}
	return
}
