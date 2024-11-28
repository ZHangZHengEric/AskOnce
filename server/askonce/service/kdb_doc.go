package service

import (
	"askonce/api/jobd"
	"askonce/components"
	"askonce/components/dto/dto_kdb"
	"askonce/components/dto/dto_kdb_doc"
	"askonce/data"
	"askonce/models"
	"askonce/utils"
	"fmt"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/xiangtao94/golib/flow"
	"github.com/xiangtao94/golib/pkg/orm"
	"path"
	"strings"
	"time"
)

type KdbDocService struct {
	flow.Service
	kdbData    *data.KdbData
	kdbDocData *data.KdbDocData
	userData   *data.UserData
}

func (k *KdbDocService) OnCreate() {
	k.kdbData = flow.Create(k.GetCtx(), new(data.KdbData))
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

	for _, doc := range docs {
		tmpp := strings.Split(doc.Path, ".")
		t := dto_kdb.DataListItem{
			KdbDataId:  doc.Id,
			Type:       doc.DataSource,
			DataName:   doc.DocName,
			DataPath:   doc.Path,
			DataSuffix: tmpp[len(tmpp)-1],
			Status:     doc.Status,
			CreateTime: doc.CreatedAt.Format(time.DateTime),
		}
		res.List = append(res.List, t)
	}
	return
}

var allowFormat = []string{"pdf", "doc", "docx", "txt", "ppt", "pptx", "xlsx", "xls"}
var needToPdfFormat = []string{"ppt", "pptx", "xlsx", "xls"}

func (k *KdbDocService) DataAdd(req *dto_kdb_doc.AddReq) (res any, err error) {
	userInfo, _ := utils.LoginInfo(k.GetCtx())
	kdb, err := k.kdbData.CheckKdbAuth(req.KdbId, userInfo.UserId, models.AuthTypeWrite)
	if err != nil {
		return
	}
	suffixName := strings.Split(req.File.Filename, ".")[1]
	if !slice.Contain(allowFormat, suffixName) {
		return nil, components.ErrorFormatError
	}
	var needToPdf bool
	if slice.Contain(needToPdfFormat, suffixName) {
		needToPdf = true
	}
	fileApi := k.Create(new(file.MinioApi)).(*file.MinioApi)
	dataSplit := []*models.KnowledgeDataSplit{}
	fileName := req.File.Filename
	objectFullPath := path.Join(fmt.Sprintf("kdb%v", kdb.Id), path.Clean(fileName))

	_, fileUrl, err := fileApi.UploadFile("rag-backend", objectFullPath, req.File)
	if err != nil {
		return nil, err
	}
	var fileToTextRes = &jobd.FileToTextRes{}
	if needToPdf {
		fileUrl2, err := k.GoUnoApi.FileToPdf(req.File)
		if err != nil {
			k.LogErrorf("FileToPdf，error: %v", err.Error())
			return nil, err
		}
		fileToTextRes, err = k.jobdApi.FileToText(fileUrl2)
		if err != nil {
			k.LogErrorf("调用文件转text报错，error: %v", err.Error())
			return nil, err
		}
	} else {
		fileToTextRes, err = k.jobdApi.FileToText(fileUrl)
		if err != nil {
			k.LogErrorf("调用文件转text报错，error: %v", err.Error())
			return nil, err
		}
	}
	fileContent := fileToTextRes.Text
	// 切分
	splitRes, err := k.ragApi.TextSplit(fileContent, data.DefaultProcessRuleDetail)
	if err != nil {
		return
	}
	for _, item := range splitRes.Text {
		dataSplit = append(dataSplit, &models.KnowledgeDataSplit{
			KdbId: req.KdbId,
			Text:  item,
			CrudModel: orm.CrudModel{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		})
	}
	db := helpers.MysqlClient.WithContext(k.GetCtx())
	tx := db.Begin()
	k.kdbDataDao.SetDB(tx)
	k.kdbDataSplitDao.SetDB(tx)
	k.kdbDataContentDao.SetDB(tx)
	add := &models.KnowledgeData{
		KdbId: kdb.Id,
		Type:  "file",
		Name:  fileName,
		Path:  fileUrl,
		CrudModel: orm.CrudModel{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	err = k.kdbDataDao.Insert(add)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	err = k.kdbDataContentDao.Insert(&models.KnowledgeDataContent{
		DataId:  add.Id,
		KdbId:   kdb.Id,
		Content: fileContent,
		CrudModel: orm.CrudModel{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	})
	for i := range dataSplit {
		dataSplit[i].DataId = add.Id
	}

	err = k.kdbDataSplitDao.BatchInsert(kdb.Id, dataSplit)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	go func(entity *KdbService) {
		dataIds := make([]int64, 0)
		for _, item := range dataSplit {
			dataIds = append(dataIds, item.DataId)
		}
		dataIds = slice.Unique(dataIds)
		err = entity.kdbDataDao.BatchUpdate(dataIds, map[string]interface{}{"status": 0})
		if err != nil {
			entity.LogErrorf("kdbDataDao update error: %v", err.Error())
		}
		err = entity.knowledgeData.InsertEs(kdb.GetIndexName(), dataSplit)
		if err != nil {
			entity.LogErrorf("es插入失败")
			err = entity.kdbDataDao.BatchUpdate(dataIds, map[string]interface{}{"status": models.KnowledgeDataFail})
		} else {
			err = entity.kdbDataDao.BatchUpdate(dataIds, map[string]interface{}{"status": models.KnowledgeDataSuccess})
		}
	}(k.CopyWithCtx(k.GetCtx()).(*KdbService))
	return
}

func (k *KdbDocService) DataDelete(req *dto_kdb_doc.DeleteReq) (res any, err error) {

	userInfo, _ := utils.LoginInfo(k.GetCtx())
	kdb, err := k.kdbData.CheckKdbAuth(req.KdbId, userInfo.UserId, models.AuthTypeWrite)
	if err != nil {
		return
	}

	err = k.kdbDataDao.DeleteByKdbIdAndDataId(kdb.Id, req.KdbDataId)
	if err != nil {
		return nil, err
	}
	err = k.kdbDataContentDao.DeleteByDataIds([]int64{req.KdbDataId})
	if err != nil {
		return nil, err
	}
	err = k.kdbDataSplitDao.DeleteByDataIds(kdb.Id, []int64{req.KdbDataId})
	if err != nil {
		return nil, err
	}

	err = k.knowledgeData.DeleteEs(kdb.GetIndexName(), req.KdbDataId)
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

	kdata, err := k.kdbDataDao.GetById(req.KdbDataId)
	if err != nil {
		return nil, err
	}
	if kdata.Status != models.KnowledgeDataFail {
		return nil, components.ErrorKdbDataRedoError
	}
	dataSplit, err := k.kdbDataSplitDao.GetByDataId(req.KdbDataId)
	if err != nil {
		return nil, err
	}
	dataIds := make([]int64, 0)
	for _, item := range dataSplit {
		dataIds = append(dataIds, item.DataId)
	}
	dataIds = slice.Unique(dataIds)
	err = k.kdbDataDao.BatchUpdate(dataIds, map[string]interface{}{"status": 0})
	if err != nil {
		k.LogErrorf("kdbDataDao update error: %v", err.Error())
	}
	err = k.knowledgeData.DeleteEs(kdb.GetIndexName(), req.KdbDataId)
	if err != nil {
		k.LogErrorf("es删除失败 error: %v", err.Error())
		return nil, err
	}
	err = k.knowledgeData.InsertEs(kdb.GetIndexName(), dataSplit)
	if err != nil {
		k.LogErrorf("es插入失败")
		err = k.kdbDataDao.BatchUpdate(dataIds, map[string]interface{}{"status": models.KnowledgeDataFail})
	} else {
		err = k.kdbDataDao.BatchUpdate(dataIds, map[string]interface{}{"status": models.KnowledgeDataSuccess})
	}
	return
}
