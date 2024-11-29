package data

import (
	"askonce/api/jobd"
	"askonce/components/dto"
	"askonce/conf"
	"askonce/helpers"
	"askonce/models"
	"encoding/json"
	"github.com/xiangtao94/golib/flow"
	"github.com/xiangtao94/golib/pkg/orm"
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
	documentData     *DocumentData
	fileData         *FileData
}

func (k *KdbDocData) OnCreate() {
	k.kdbDocDao = flow.Create(k.GetCtx(), new(models.KdbDocDao))
	k.kdbDocContentDao = flow.Create(k.GetCtx(), new(models.KdbDocContentDao))
	k.kdbDocSegmentDao = flow.Create(k.GetCtx(), new(models.KdbDocSegmentDao))
	k.jobdApi = flow.Create(k.GetCtx(), new(jobd.JobdApi))
	k.documentData = flow.Create(k.GetCtx(), new(DocumentData))
	k.fileData = flow.Create(k.GetCtx(), new(FileData))
}

func (k *KdbDocData) GetDocList(kdbId int64, queryName string, pageParam dto.PageParam) (list []*models.KdbDoc, cnt int64, err error) {
	list = make([]*models.KdbDoc, 0)
	list, cnt, err = k.kdbDocDao.GetList(kdbId, queryName, pageParam)
	return
}

func (k *KdbDocData) AddDocFormFile(kdbId int64, userId string, file *models.File, split bool) (add *models.KdbDoc, err error) {
	add = &models.KdbDoc{
		KdbId:      kdbId,
		DocName:    file.OriginName,
		DataSource: "file",
		SourceId:   file.Id,
		NeedSplit:  split,
		Status:     models.KdbDocWaiting,
		UserId:     userId,
		CrudModel: orm.CrudModel{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	err = k.kdbDocDao.Insert(add)

	return
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
		esDbConfigStr := strings.Replace(conf.WebConf.EsDbConfig, "${indexName}", kdb.GetIndexName(), 1)
		_, err = k.jobdApi.AtomEsDelete(jobd.ESDeleteReq{DocIds: []string{
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

// 文档构建到内存数据库
func (k *KdbDocData) DocBuild(docIds []int64) (err error) {
	//1. 文件导入加锁

	//2. 文件解析文本

	//3. 文本切分

	//4. 文本转向量

	//5. 存向量数据库

	// 6.更新db
	return
}
