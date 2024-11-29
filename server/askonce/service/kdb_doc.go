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
	"os"
	"time"
)

type KdbDocService struct {
	flow.Service
	kdbData    *data.KdbData
	fileData   *data.FileData
	kdbDocData *data.KdbDocData
	userData   *data.UserData
}

func (k *KdbDocService) OnCreate() {
	k.kdbData = flow.Create(k.GetCtx(), new(data.KdbData))
	k.fileData = flow.Create(k.GetCtx(), new(data.FileData))
	k.kdbDocData = flow.Create(k.GetCtx(), new(data.KdbDocData))
	k.userData = flow.Create(k.GetCtx(), new(data.UserData))
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
	docs, cnt, err := k.kdbDocData.GetDocList(kdb.Id, req.QueryName, req.PageParam)
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
	if req.Type == "text" {
		fileName := ""
		if len(req.Title) > 0 {
			fileName = fmt.Sprintf("%s.txt", req.Title)
		} else {
			fileName = fmt.Sprintf("%s.txt", helpers.GenID())
		}
		if len([]rune(req.Text)) < 1024 {
			needSplit = false
		}
		tmpFile, err := os.CreateTemp("", fileName)
		if err != nil {
			return nil, err
		}
		defer func() {
			tmpFile.Close()
			os.Remove(tmpFile.Name())
		}()
		if _, err := tmpFile.Write([]byte(req.Text)); err != nil {
			return nil, components.ErrorFileUploadError
		}
		tmpFileHeader, err := utils.GetFileHeader(tmpFile)
		if err != nil {
			return nil, err
		}
		req.File = tmpFileHeader
	}
	file, err := k.fileData.Upload(userInfo.UserId, req.File, "knowledge")
	if err != nil {
		return nil, err
	}
	doc, err := k.kdbDocData.AddDocFormFile(kdb.Id, userInfo.UserId, file, needSplit)
	if err != nil {
		return nil, components.ErrorFileUploadError
	}
	go func(k *KdbDocService) {
		err = k.kdbDocData.DocBuild([]int64{doc.Id})
		if err != nil {
			k.LogErrorf("文档构建内存数据库失败 %s", err.Error())
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
	_, err = k.kdbData.CheckKdbAuth(req.KdbId, userInfo.UserId, models.AuthTypeWrite)
	if err != nil {
		return
	}
	go func(k *KdbDocService) {
		err = k.kdbDocData.DocBuild([]int64{req.DocId})
		if err != nil {
			k.LogErrorf("文档构建内存数据库失败 %s", err.Error())
		}
	}(k.CopyWithCtx(k.GetCtx()).(*KdbDocService))
	return
}
