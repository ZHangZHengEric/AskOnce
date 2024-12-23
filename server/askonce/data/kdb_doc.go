package data

import (
	"askonce/api/jobd"
	"askonce/components/dto"
	"askonce/components/dto/dto_kdb_doc"
	"askonce/es"
	"askonce/helpers"
	"askonce/models"
	"askonce/utils"
	"fmt"
	"github.com/xiangtao94/golib/flow"
	"github.com/xiangtao94/golib/pkg/errors"
	"github.com/xiangtao94/golib/pkg/orm"
	"time"
)

type KdbDocData struct {
	flow.Redis
	kdbDocDao        *models.KdbDocDao
	fileData         *FileData
	datasourceData   *DatasourceData
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
	k.datasourceData = flow.Create(k.GetCtx(), new(DatasourceData))
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
	datasourceIds := make([]string, 0)
	for _, doc := range docs {
		if doc.Status == models.KdbDocSuccess {
			docSuccessIds = append(docSuccessIds, doc.Id)
		}
		if doc.DataSource == models.DataSourceFile {
			docFileIds = append(docFileIds, doc.SourceId)
		} else {
			datasourceIds = append(datasourceIds, doc.SourceId)
		}

	}
	err = k.datasourceData.DeleteByIds(datasourceIds)
	if err != nil {
		tx.Rollback()
		return err
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

func (k *KdbDocData) GetList(kdbId int64, queryName string, queryStatus []int, param dto.PageParam) (list []dto_kdb_doc.InfoItem, cnt int64, err error) {
	docs, cnt, err := k.kdbDocDao.GetList(kdbId, queryName, queryStatus, param)
	if err != nil {
		return
	}
	fileIds := make([]string, 0)
	datasourceIds := make([]string, 0)
	for _, doc := range docs {
		if doc.DataSource == models.DataSourceDatabase {
			datasourceIds = append(datasourceIds, doc.SourceId)
			continue
		}
		fileIds = append(fileIds, doc.SourceId)
	}
	fileMap, err := k.fileData.GetFileByFileIds(fileIds)
	if err != nil {
		return
	}
	datasourceMap, err := k.datasourceData.GetByIds(datasourceIds)
	if err != nil {
		return
	}
	for _, doc := range docs {
		t := dto_kdb_doc.InfoItem{
			Id:         doc.Id,
			Type:       doc.DataSource,
			DataName:   doc.DocName,
			Status:     doc.Status,
			CreateTime: doc.CreatedAt.Format(time.DateTime),
		}

		if file, ok := fileMap[doc.SourceId]; ok {
			t.DataPath = file.Path
			t.DataSuffix = file.Extension
		}
		if datasource, ok := datasourceMap[doc.SourceId]; ok {
			t.DbType = datasource.Type
		}
		list = append(list, t)
	}
	return
}

func (k *KdbDocData) AddDoc(req *dto_kdb_doc.AddReq, kdb *models.Kdb, userId string) (doc *models.KdbDoc, err error) {
	var sourceId, sourceName, sourceType string
	switch req.Type {
	case "text":
		sourceType = models.DataSourceFile
		if len(req.Text) == 0 {
			return nil, errors.NewError(10034, "文本内容为空！")
		}
		file, err := k.fileData.UploadByText(userId, req.Title, req.Text, "knowledge")
		if err != nil {
			return nil, err
		}
		sourceId = file.Id
		sourceName = file.OriginName
	case "file":
		sourceType = models.DataSourceFile
		file, err := k.fileData.UploadByFile(userId, req.File, "knowledge")
		if err != nil {
			return nil, err
		}
		sourceId = file.Id
		sourceName = file.OriginName
	case "database":
		sourceType = models.DataSourceDatabase
		datasource, err := k.datasourceData.Add(userId, req.ImportDataBase)
		if err != nil {
			return nil, err
		}
		sourceId = datasource.Id
		sourceName = datasource.DatabaseName
	}
	doc = &models.KdbDoc{
		KdbId:      kdb.Id,
		DocName:    sourceName,
		DataSource: sourceType,
		SourceId:   sourceId,
		Status:     models.KdbDocRunning,
		UserId:     userId,
		Metadata:   req.Metadata,
		CrudModel: orm.CrudModel{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	err = k.kdbDocDao.Insert(doc)
	if err != nil {
		return
	}
	return
}

func (k *KdbDocData) GetDoc(kdbId int64, kdbDataId int64) (res dto_kdb_doc.InfoItem, err error) {
	doc, err := k.kdbDocDao.GetById(kdbDataId)
	if err != nil {
		return
	}
	res = dto_kdb_doc.InfoItem{
		Id:         doc.Id,
		Type:       doc.DataSource,
		DataName:   doc.DocName,
		Status:     doc.Status,
		CreateTime: doc.CreatedAt.Format(time.DateTime),
	}
	if doc.DataSource == models.DataSourceDatabase {
		datasourceMap, err := k.datasourceData.GetByIds([]string{doc.SourceId})
		if err != nil {
			return res, err
		}

		if datasource, ok := datasourceMap[doc.SourceId]; ok {
			res.DbType = datasource.Type
			res.DbSchema = datasource.Schema
		}
	} else {
		fileMap, err := k.fileData.GetFileByFileIds([]string{doc.SourceId})
		if err != nil {
			return res, err
		}
		if file, ok := fileMap[doc.SourceId]; ok {
			res.DataPath = file.Path
			res.DataSuffix = file.Extension
		}
	}
	return
}

func (k *KdbDocData) AddDocByFiles(kdb *models.Kdb, files []*models.File, userId string) (taskId string, err error) {
	taskId = utils.Get16MD5Encode(fmt.Sprintf("%s%v", userId, time.Now().UnixNano()))
	docs := make([]*models.KdbDoc, 0)
	for _, file := range files {
		doc := &models.KdbDoc{
			KdbId:      kdb.Id,
			TaskId:     taskId,
			DocName:    file.OriginName,
			DataSource: models.DataSourceFile,
			SourceId:   file.Id,
			Status:     models.KdbDocWaiting,
			UserId:     userId,
			Metadata:   file.Metadata,
			CrudModel: orm.CrudModel{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}
		docs = append(docs, doc)
	}
	err = k.kdbDocDao.BatchInsert(docs)
	if err != nil {
		return
	}
	return
}

func (k *KdbDocData) AddDocByDatasource(kdb *models.Kdb, datasource *models.Datasource, userId string) (err error) {
	doc := &models.KdbDoc{
		KdbId:      kdb.Id,
		DocName:    datasource.DatabaseName,
		DataSource: models.DataSourceDatabase,
		SourceId:   datasource.Id,
		Status:     models.KdbDocWaiting,
		UserId:     userId,
		CrudModel: orm.CrudModel{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	err = k.kdbDocDao.Insert(doc)
	return
}
