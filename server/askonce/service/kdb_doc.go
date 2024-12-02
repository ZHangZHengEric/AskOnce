package service

import (
	"askonce/components"
	"askonce/components/dto/dto_kdb"
	"askonce/components/dto/dto_kdb_doc"
	"askonce/data"
	"askonce/helpers"
	"askonce/models"
	"askonce/utils"
	"fmt"
	"github.com/xiangtao94/golib/flow"
	"github.com/xiangtao94/golib/pkg/orm"
	"time"
)

type KdbDocService struct {
	flow.Service
	kdbData      *data.KdbData
	fileData     *data.FileData
	documentData *data.DocumentData
	kdbDocData   *data.KdbDocData
	kdbDocDao    *models.KdbDocDao
}

func (k *KdbDocService) OnCreate() {
	k.kdbData = flow.Create(k.GetCtx(), new(data.KdbData))
	k.fileData = flow.Create(k.GetCtx(), new(data.FileData))
	k.documentData = flow.Create(k.GetCtx(), new(data.DocumentData))
	k.kdbDocData = flow.Create(k.GetCtx(), new(data.KdbDocData))
	k.kdbDocDao = flow.Create(k.GetCtx(), new(models.KdbDocDao))
}

func (k *KdbDocService) DocList(req *dto_kdb_doc.ListReq) (res *dto_kdb.DataListResp, err error) {
	userInfo, _ := utils.LoginInfo(k.GetCtx())
	kdb, err := k.kdbData.CheckKdbAuth(req.KdbId, userInfo.UserId, models.AuthTypeRead)
	if err != nil {
		return
	}
	res = &dto_kdb.DataListResp{
		List:  make([]dto_kdb.DataListItem, 0),
		Total: 0,
	}
	docs, cnt, err := k.kdbDocDao.GetList(kdb.Id, req.QueryName, req.QueryStatus, req.PageParam)
	if err != nil {
		return nil, err
	}
	res.Total = cnt
	fileIds := make([]string, 0)
	for _, doc := range docs {
		fileIds = append(fileIds, doc.SourceId)
	}
	fileMap, err := k.fileData.GetFileByFileIds(fileIds)
	if err != nil {
		return nil, err
	}
	for _, doc := range docs {
		t := dto_kdb.DataListItem{
			KdbDataId: doc.Id,
			Type:      doc.DataSource,
			DataName:  doc.DocName,

			Status:     doc.Status,
			CreateTime: doc.CreatedAt.Format(time.DateTime),
		}
		file := fileMap[doc.SourceId]
		if file != nil {
			t.DataPath = file.Path
			t.DataSuffix = file.Extension
		}
		res.List = append(res.List, t)
	}
	return
}

func (k *KdbDocService) DocAdd(req *dto_kdb_doc.AddReq) (res interface{}, err error) {
	userInfo, _ := utils.LoginInfo(k.GetCtx())
	kdb, err := k.kdbData.CheckKdbAuth(req.KdbId, userInfo.UserId, models.AuthTypeWrite)
	if err != nil {
		return
	}
	needSplit := true
	var file *models.File
	if req.Type == "text" {
		fileName := ""
		if len(req.Title) > 0 {
			fileName = fmt.Sprintf("%s.txt", req.Title)
		} else {
			fileName = fmt.Sprintf("%v.txt", helpers.GenID())
		}
		if len([]rune(req.Text)) < 1024 {
			needSplit = false
		}
		file, err = k.fileData.UploadContent(userInfo.UserId, fileName, req.Text, "knowledge")
		if err != nil {
			return nil, err
		}
	} else {
		file, err = k.fileData.Upload(userInfo.UserId, req.File, "knowledge")
		if err != nil {
			return nil, err
		}
	}

	doc := &models.KdbDoc{
		KdbId:      kdb.Id,
		DocName:    file.OriginName,
		DataSource: "file",
		SourceId:   file.Id,
		NeedSplit:  needSplit,
		Status:     models.KdbDocRunning,
		UserId:     userInfo.UserId,
		CrudModel: orm.CrudModel{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	err = k.kdbDocDao.Insert(doc)
	if err != nil {
		return nil, err
	}
	go func(k *KdbDocService) {
		defer func() {
			if r := recover(); r != nil {
				k.LogErrorf("文档【%v】构建内存数据库失败 %s", doc.Id, r)
				_ = k.kdbDocDao.UpdateStatus(doc.Id, models.KdbDocFail)
			}
		}()
		_ = k.kdbDocDao.UpdateStatus(doc.Id, models.KdbDocRunning)
		err = k.DocBuild(kdb, doc)
		if err != nil {
			k.LogErrorf("文档【%v】构建内存数据库失败 %s", doc.Id, err.Error())
			_ = k.kdbDocDao.UpdateStatus(doc.Id, models.KdbDocFail)
		} else {
			k.LogInfof("文档【%v】构建内存数据库成功 %s", doc.Id)
			_ = k.kdbDocDao.UpdateStatus(doc.Id, models.KdbDocSuccess)
		}
	}(k.CopyWithCtx(k.GetCtx()).(*KdbDocService))
	return
}

func (k *KdbDocService) DocDelete(req *dto_kdb_doc.DeleteReq) (res interface{}, err error) {
	userInfo, _ := utils.LoginInfo(k.GetCtx())
	kdb, err := k.kdbData.CheckKdbAuth(req.KdbId, userInfo.UserId, models.AuthTypeWrite)
	if err != nil {
		return
	}
	err = k.kdbDocData.DeleteDoc(kdb, req.DocId)
	if err != nil {
		return nil, err
	}
	return
}

func (k *KdbDocService) DataRedo(req *dto_kdb_doc.RedoReq) (res any, err error) {
	userInfo, _ := utils.LoginInfo(k.GetCtx())
	kdb, err := k.kdbData.CheckKdbAuth(req.KdbId, userInfo.UserId, models.AuthTypeWrite)
	if err != nil {
		return
	}
	go func(k *KdbDocService) {
		defer func() {
			if r := recover(); r != nil {
				k.LogErrorf("文档【%v】构建内存数据库失败 %s", req.DocId, r)
				_ = k.kdbDocDao.UpdateStatus(req.DocId, models.KdbDocFail)
			}
		}()
		_ = k.kdbDocDao.UpdateStatus(req.DocId, models.KdbDocRunning)
		doc, err := k.kdbDocDao.GetById(req.DocId)
		if err != nil {
			return
		}
		_ = k.kdbDocDao.UpdateStatus(req.DocId, models.KdbDocSuccess)
		err = k.DocBuild(kdb, doc)
		if err != nil {
			k.LogErrorf("文档【%v】构建内存数据库失败 %s", req.DocId, err.Error())
			_ = k.kdbDocDao.UpdateStatus(doc.Id, models.KdbDocFail)
		} else {
			k.LogInfof("文档【%v】构建内存数据库成功", req.DocId)
			_ = k.kdbDocDao.UpdateStatus(doc.Id, models.KdbDocSuccess)
		}
	}(k.CopyWithCtx(k.GetCtx()).(*KdbDocService))
	return
}

// 文档构建到内存数据库
func (k *KdbDocService) DocBuild(kdb *models.Kdb, doc *models.KdbDoc) (err error) {
	//2. 文件解析文本段
	k.LogInfof("开始文件解析文本，docId %v", doc.Id)
	_, content, err := k.fileData.ConvertFileToText(doc.SourceId)
	if err != nil {
		k.LogErrorf("文件解析文本，docId %v, error %v", doc.Id, err.Error())
		return err
	}
	//3. 文本切分
	k.LogInfof("开始文本切分，docId %v", doc.Id)
	splitList, err := k.documentData.TextSplit(content)
	if err != nil {
		k.LogErrorf("文本切分error，docId %v,error %v", doc.Id, err.Error())
		return components.ErrorTextSplitError
	}
	if len(splitList) == 0 {
		return components.ErrorTextSplitError
	}
	contents := make([]string, 0, len(splitList))
	for _, split := range splitList {
		contents = append(contents, split.PassageContent)
	}
	//4. 文本转向量
	k.LogInfof("开始文本转向量，docId %v", doc.Id)
	embeddingAll, err := k.documentData.TextEmbedding(contents)
	if err != nil || len(embeddingAll) != len(contents) {
		k.LogErrorf("文本转向量error，docId %v,error %v", doc.Id, err.Error())
		return err
	}
	//5. 存向量数据库和mysql
	k.LogInfof("开始存数据库，docId %v", doc.Id)
	err = k.kdbDocData.SaveDocBuild(kdb, doc, content, splitList, embeddingAll)
	if err != nil {
		k.LogErrorf("存mysql error，docId %v,error %v", doc.Id, err.Error())
		return err
	}
	return
}
