package service

import (
	"askonce/components"
	"askonce/components/dto/dto_kdb"
	"askonce/components/dto/dto_kdb_doc"
	"askonce/data"
	"askonce/models"
	"askonce/utils"
	"fmt"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/xiangtao94/golib/flow"
	"github.com/xiangtao94/golib/pkg/errors"
	"github.com/xiangtao94/golib/pkg/orm"
	"sync"
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
	docs, cnt, err := k.kdbDocData.GetList(kdb.Id, req.QueryName, req.QueryStatus, req.PageParam)
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
			Id:       doc.Id,
			Type:     doc.DataSource,
			DataName: doc.DocName,

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

func (k *KdbDocService) DocAdd(req *dto_kdb_doc.AddReq) (res *dto_kdb_doc.AddRes, err error) {
	userInfo, _ := utils.LoginInfo(k.GetCtx())
	kdb, err := k.kdbData.CheckKdbAuth(req.KdbId, userInfo.UserId, models.AuthTypeWrite)
	if err != nil {
		return
	}
	var file *models.File
	if req.Type == "text" {
		if len(req.Text) == 0 {
			return nil, errors.NewError(10034, "文本内容为空！")
		}
		file, err = k.fileData.UploadByText(userInfo.UserId, req.Title, req.Text, "knowledge")
		if err != nil {
			return nil, err
		}
	} else {
		file, err = k.fileData.UploadByFile(userInfo.UserId, req.File, "knowledge")
		if err != nil {
			return nil, err
		}
	}
	doc := &models.KdbDoc{
		KdbId:      kdb.Id,
		DocName:    file.OriginName,
		DataSource: models.DataTypeCommon,
		SourceId:   file.Id,
		NeedSplit:  true,
		Status:     models.KdbDocRunning,
		UserId:     userInfo.UserId,
		CrudModel: orm.CrudModel{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	k.kdbDocData.AddDoc(doc)
	go func(k *KdbDocService) {
		_ = k.DocBuild(kdb, doc)
	}(k.CopyWithCtx(k.GetCtx()).(*KdbDocService))
	res = &dto_kdb_doc.AddRes{
		KdbDataId: doc.Id,
	}
	return
}

func (k *KdbDocService) DocDelete(req *dto_kdb_doc.DeleteReq) (res interface{}, err error) {
	userInfo, _ := utils.LoginInfo(k.GetCtx())
	kdb, err := k.kdbData.CheckKdbAuth(req.KdbId, userInfo.UserId, models.AuthTypeWrite)
	if err != nil {
		return
	}
	err = k.kdbDocData.DeleteDocs(kdb, []int64{req.DocId}, false)
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
	doc, err := k.kdbDocDao.GetById(req.DocId)
	if err != nil {
		return
	}
	go func(k *KdbDocService) {
		_ = k.DocBuild(kdb, doc)
	}(k.CopyWithCtx(k.GetCtx()).(*KdbDocService))
	return
}

func (k *KdbDocService) DocAddZip(req *dto_kdb_doc.AddZipReq) (res *dto_kdb_doc.AddZipRes, err error) {
	taskId := utils.Get16MD5Encode(fmt.Sprintf("%s%v", req.ZipUrl, time.Now().UnixNano()))
	userInfo, _ := utils.LoginInfo(k.GetCtx())
	kdb, err := k.kdbData.CheckKdbAuth(req.KdbId, userInfo.UserId, models.AuthTypeWrite)
	if err != nil {
		return
	}
	files, err := k.fileData.UploadByZip(userInfo.UserId, req.ZipUrl, "knowledge")
	if err != nil {
		return
	}
	docs := make([]*models.KdbDoc, 0)
	for _, file := range files {
		doc := &models.KdbDoc{
			KdbId:      kdb.Id,
			TaskId:     taskId,
			DocName:    file.OriginName,
			DataSource: "file",
			SourceId:   file.Id,
			NeedSplit:  true,
			Status:     models.KdbDocWaiting,
			UserId:     userInfo.UserId,
			CrudModel: orm.CrudModel{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}
		docs = append(docs, doc)
	}
	err = k.kdbDocDao.BatchInsert(docs)
	if err != nil {
		return nil, err
	}
	res = &dto_kdb_doc.AddZipRes{
		TaskId: taskId,
	}
	return
}

func (k *KdbDocService) DocAddByBatchText(req *dto_kdb_doc.AddByBatchTextReq) (res interface{}, err error) {
	userInfo, err := utils.LoginInfo(k.GetCtx())
	if err != nil {
		return nil, err
	}
	kdb, err := k.kdbData.GetKdbByName(req.KdbName, userInfo, req.AutoCreate)
	if err != nil {
		return nil, err
	}
	docs := make([]*models.KdbDoc, 0)
	for _, doc := range req.Docs {
		if len(doc.Text) == 0 {
			return nil, errors.NewError(10034, "文本内容为空！")
		}
		file, err := k.fileData.UploadByText(userInfo.UserId, doc.Title, doc.Text, "knowledge")
		if err != nil {
			return nil, err
		}
		doc := &models.KdbDoc{
			KdbId:      kdb.Id,
			DocName:    file.OriginName,
			DataSource: "file",
			SourceId:   file.Id,
			NeedSplit:  true,
			Status:     models.KdbDocWaiting,
			UserId:     userInfo.UserId,
			CrudModel: orm.CrudModel{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}
		docs = append(docs, doc)
	}
	err = k.kdbDocDao.BatchInsert(docs)
	if err != nil {
		return nil, err
	}
	return
}

func (k *KdbDocService) BuildWaitingDoc() (err error) {
	docs, err := k.kdbDocDao.GetListByStatus(models.KdbDocWaiting)
	if err != nil {
		return
	}
	kdbIds := make([]int64, 0)
	for _, doc := range docs {
		kdbIds = append(kdbIds, doc.KdbId)
	}
	kdbs, err := k.kdbData.GetKdbByIds(kdbIds)
	kdbMap := make(map[int64]*models.Kdb)
	for _, kdb := range kdbs {
		kdbMap[kdb.Id] = kdb
	}
	wg := sync.WaitGroup{}
	for _, doc := range docs {
		wg.Add(1)
		kdb := kdbMap[doc.KdbId]
		go func(k *KdbDocService) {
			defer wg.Done()
			_ = k.DocBuild(kdb, doc)
		}(k.CopyWithCtx(k.GetCtx()).(*KdbDocService))
	}
	wg.Wait()
	return
}

func (k *KdbDocService) DocBuild(kdb *models.Kdb, doc *models.KdbDoc) (err error) {
	_ = k.kdbDocDao.UpdateStatus(doc.Id, models.KdbDocRunning)
	err = k.docBuildDo(kdb, doc)
	if err != nil {
		k.LogErrorf("文档【%v】构建内存数据库失败 %s", doc.Id, err.Error())
		_ = k.kdbDocDao.UpdateStatus(doc.Id, models.KdbDocFail)
		return
	}
	k.LogInfof("文档【%v】构建内存数据库成功", doc.Id)
	_ = k.kdbDocDao.UpdateStatus(doc.Id, models.KdbDocSuccess)
	return
}

// 文档构建到内存数据库
func (k *KdbDocService) docBuildDo(kdb *models.Kdb, doc *models.KdbDoc) (err error) {
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

func (k *KdbDocService) LoadProcess(req *dto_kdb_doc.TaskProcessReq) (res *dto_kdb.LoadProcessRes, err error) {
	userInfo, err := utils.LoginInfo(k.GetCtx())
	if err != nil {
		return nil, err
	}
	kdb, err := k.kdbData.CheckKdbAuth(req.KdbId, userInfo.UserId, models.AuthTypeSuperAdmin)
	if err != nil {
		return nil, err
	}
	processRes, err := k.kdbDocDao.QueryProcess(kdb.Id, req.TaskId)
	if err != nil {
		return nil, err
	}
	res = &dto_kdb.LoadProcessRes{}
	var totalNum int64

	for _, p := range processRes {
		switch p.Status {
		case models.KdbDocWaiting:
			res.Waiting = p.Total
		case models.KdbDocFail:
			res.Fail = p.Total
		case models.KdbDocSuccess:
			res.Success = p.Total
		case models.KdbDocRunning:
			res.InProgress = p.Total
		default:
		}
		totalNum = totalNum + p.Total
	}
	res.Total = totalNum
	if res.Total == res.Success {
		res.TaskProcess = 100
	} else {
		res.TaskProcess = (res.Success * 100) / totalNum
	}
	return
}

func (k *KdbDocService) TaskRedo(req *dto_kdb_doc.TaskRedoReq) (res interface{}, err error) {
	docs, err := k.kdbDocDao.GetByTaskIdAndStatus(req.KdbId, req.TaskId, []int{models.KdbDocFail})
	if err != nil {
		return nil, err
	}
	docIds := make([]int64, 0, len(docs))
	for _, doc := range docs {
		docIds = append(docIds, doc.Id)
	}
	k.LogInfof("重做邮件，ids【%v】", slice.Join(docIds, ","))
	err = k.kdbDocDao.BatchUpdateStatus(docIds, models.KdbDocWaiting)
	if err != nil {
		return nil, err
	}
	return
}
