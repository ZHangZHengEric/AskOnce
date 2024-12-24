package service

import (
	"askonce/components"
	"askonce/components/dto/dto_kdb"
	"askonce/components/dto/dto_kdb_doc"
	"askonce/data"
	"askonce/es"
	"askonce/models"
	"askonce/utils"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/xiangtao94/golib/flow"
	"github.com/xiangtao94/golib/pkg/errors"
	"github.com/xiangtao94/golib/pkg/zlog"
	"sync"
)

type KdbDocService struct {
	flow.Service
	kdbData        *data.KdbData
	fileData       *data.FileData
	datasourceData *data.DatasourceData
	documentData   *data.DocumentData
	kdbDocData     *data.KdbDocData
	kdbDocDao      *models.KdbDocDao
}

func (k *KdbDocService) OnCreate() {
	k.kdbData = flow.Create(k.GetCtx(), new(data.KdbData))
	k.fileData = flow.Create(k.GetCtx(), new(data.FileData))
	k.datasourceData = flow.Create(k.GetCtx(), new(data.DatasourceData))
	k.documentData = flow.Create(k.GetCtx(), new(data.DocumentData))
	k.kdbDocData = flow.Create(k.GetCtx(), new(data.KdbDocData))
	k.kdbDocDao = flow.Create(k.GetCtx(), new(models.KdbDocDao))
}

func (k *KdbDocService) DocList(req *dto_kdb_doc.ListReq) (res *dto_kdb_doc.ListResp, err error) {
	userInfo, _ := utils.LoginInfo(k.GetCtx())
	kdb, err := k.kdbData.CheckKdbAuth(req.KdbId, userInfo.UserId, models.AuthTypeRead)
	if err != nil {
		return
	}
	res = &dto_kdb_doc.ListResp{
		List:  make([]dto_kdb_doc.InfoItem, 0),
		Total: 0,
	}
	list, cnt, err := k.kdbDocData.GetList(kdb.Id, req.QueryName, req.QueryStatus, req.PageParam)
	if err != nil {
		return nil, err
	}
	res.Total = cnt
	res.List = list
	return
}

func (k *KdbDocService) DocInfo(req *dto_kdb_doc.InfoReq) (res *dto_kdb_doc.InfoRes, err error) {
	userInfo, _ := utils.LoginInfo(k.GetCtx())
	kdb, err := k.kdbData.CheckKdbAuth(req.KdbId, userInfo.UserId, models.AuthTypeRead)
	if err != nil {
		return
	}
	info, err := k.kdbDocData.GetDoc(kdb.Id, req.DataId)
	if err != nil {
		return nil, err
	}
	res = &dto_kdb_doc.InfoRes{
		InfoItem: info,
	}
	return
}

func (k *KdbDocService) DocAdd(req *dto_kdb_doc.AddReq) (res *dto_kdb_doc.AddRes, err error) {
	userInfo, _ := utils.LoginInfo(k.GetCtx())
	kdb, err := k.kdbData.CheckKdbAuth(req.KdbId, userInfo.UserId, models.AuthTypeWrite)
	if err != nil {
		return
	}
	switch req.Type {
	case "text":
		if len(req.Text) == 0 {
			return nil, errors.NewError(10034, "文本内容为空！")
		}
		file, err := k.fileData.UploadByText(userInfo.UserId, req.Title, req.Text, "knowledge")
		if err != nil {
			return nil, err
		}
		_, err = k.kdbDocData.AddDocByFiles(kdb, []*models.File{file}, userInfo.UserId)
		if err != nil {
			return nil, err
		}
	case "file":
		file, err := k.fileData.UploadByFile(userInfo.UserId, req.File, "knowledge")
		if err != nil {
			return nil, err
		}
		_, err = k.kdbDocData.AddDocByFiles(kdb, []*models.File{file}, userInfo.UserId)
		if err != nil {
			return nil, err
		}
	case "database":
		datasource, err := k.datasourceData.Add(userInfo.UserId, req.ImportDataBase)
		if err != nil {
			return nil, err
		}
		err = k.kdbDocData.AddDocByDatasource(kdb, datasource, userInfo.UserId)
		if err != nil {
			return nil, err
		}
	}
	return
}

func (k *KdbDocService) DocDelete(req *dto_kdb_doc.DeleteReq) (res interface{}, err error) {
	userInfo, _ := utils.LoginInfo(k.GetCtx())
	kdb, err := k.kdbData.CheckKdbAuth(req.KdbId, userInfo.UserId, models.AuthTypeWrite)
	if err != nil {
		return
	}
	err = k.kdbDocData.DeleteDocs(kdb, []int64{req.DataId}, false)
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
	doc, err := k.kdbDocDao.GetById(req.DataId)
	if err != nil {
		return
	}
	go func(k *KdbDocService) {
		_ = k.DocBuild(kdb, doc)
	}(k.CopyWithCtx(k.GetCtx()).(*KdbDocService))
	return
}

func (k *KdbDocService) DocAddByZip(req *dto_kdb_doc.AddZipReq) (res *dto_kdb_doc.AddZipRes, err error) {
	userInfo, _ := utils.LoginInfo(k.GetCtx())
	kdb, err := k.kdbData.CheckKdbAuth(req.KdbId, userInfo.UserId, models.AuthTypeWrite)
	if err != nil {
		return
	}
	files, err := k.fileData.UploadByZip(userInfo.UserId, req.ZipUrl, "knowledge")
	if err != nil {
		return
	}
	taskId, err := k.kdbDocData.AddDocByFiles(kdb, files, userInfo.UserId)
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
	files := make([]*models.File, 0)
	for _, doc := range req.Docs {
		if len(doc.Text) == 0 {
			return nil, errors.NewError(10034, "文本内容为空！")
		}
		file, err := k.fileData.UploadByText(userInfo.UserId, doc.Title, doc.Text, "knowledge")
		if err != nil {
			return nil, err
		}
		file.Metadata = doc.Metadata
		files = append(files, file)
	}
	_, err = k.kdbDocData.AddDocByFiles(kdb, files, userInfo.UserId)
	return
}

func (k *KdbDocService) BuildWaitingDoc() (err error) {
	docs, err := k.kdbDocDao.GetListByStatus(models.KdbDocWaiting)
	if err != nil {
		return
	}
	if len(docs) == 0 {
		return
	}
	kdbIds := make([]int64, 0)
	for _, doc := range docs {
		kdbIds = append(kdbIds, doc.KdbId)
	}
	zlog.Infof(k.GetCtx(), "定时处理waitting数据 %v", kdbIds)
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

func (k *KdbDocService) BuildFailedDoc() (err error) {
	docs, err := k.kdbDocDao.GetFailedList()
	if err != nil {
		return
	}
	if len(docs) == 0 {
		return
	}
	kdbIds := make([]int64, 0)
	for _, doc := range docs {
		kdbIds = append(kdbIds, doc.KdbId)
	}
	zlog.Infof(k.GetCtx(), "定时处理failed数据 %v", kdbIds)
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
	_ = k.kdbDocDao.UpdateStatus(doc, models.KdbDocRunning)
	if doc.DataSource == models.DataSourceFile {
		err = k.docBuildDo(kdb, doc)
	} else if doc.DataSource == models.DataSourceDatabase {
		err = k.databaseBuildDo(kdb, doc)
	}
	if err != nil {
		k.LogErrorf("文档【%v】构建内存数据库失败 %s", doc.Id, err.Error())
		_ = k.kdbDocDao.UpdateStatus(doc, models.KdbDocFail)
		return
	}
	k.LogInfof("文档【%v】构建内存数据库成功", doc.Id)
	_ = k.kdbDocDao.UpdateStatus(doc, models.KdbDocSuccess)
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

func (k *KdbDocService) databaseBuildDo(kdb *models.Kdb, doc *models.KdbDoc) (err error) {
	k.LogInfof("开始获取表结构和数据，docId %v", doc.Id)
	datasource, schemaColumns, err := k.datasourceData.GetSchemaAndValues(doc.SourceId)
	if err != nil {
		return err
	}
	// 转换成需要存的document
	tables := make([]*es.DatabaseDocument, 0)
	tableColumns := make([]*es.DatabaseDocument, 0)
	tableColumnValues := make([]*es.DatabaseDocument, 0)
	for _, table := range schemaColumns {
		tables = append(tables, &es.DatabaseDocument{
			DocDocument: es.DocDocument{
				CommonDocument: es.CommonDocument{
					DocId:      doc.Id,
					DocContent: table.FormatTableInfo(), // 处理文字
				},
			},
			DatabaseName: datasource.DatabaseName,
			TableName:    table.TableName,
			TableComment: table.TableComment,
		})
		for _, column := range table.ColumnInfos {
			tableColumns = append(tableColumns, &es.DatabaseDocument{
				DocDocument: es.DocDocument{
					CommonDocument: es.CommonDocument{
						DocId:      doc.Id,
						DocContent: column.FormatColumnInfo(), // 处理文字
					},
				},
				DatabaseName:  datasource.DatabaseName,
				TableName:     table.TableName,
				ColumnName:    column.ColumnName,
				ColumnComment: column.ColumnComment,
				ColumnType:    column.ColumnType,
			})
			for _, v := range column.ColumnValues {
				tableColumnValues = append(tableColumnValues, &es.DatabaseDocument{
					DocDocument: es.DocDocument{
						CommonDocument: es.CommonDocument{
							DocId:      doc.Id,
							DocContent: v.FormatValueInfo(),
						},
					},
					DatabaseName: datasource.DatabaseName,
					TableName:    table.TableName,
					ColumnName:   column.ColumnName,
				})
			}
		}
	}
	k.LogInfof("开始文本转向量，docId %v", doc.Id)
	tableContents := make([]string, 0, len(tables))
	tableColumnContents := make([]string, 0, len(tableColumns))
	tableColumnValueContents := make([]string, 0, len(tableColumnValues))
	for _, table := range tables {
		tableContents = append(tableContents, table.DocContent)
	}
	embTables, err := k.documentData.TextEmbedding(tableContents)
	if err != nil {
		return err
	}
	for i := range tables {
		tables[i].Emb = embTables[i]
	}

	for _, tc := range tableColumns {
		tableColumnContents = append(tableColumnContents, tc.DocContent)
	}
	embTccs, err := k.documentData.TextEmbedding(tableColumnContents)
	if err != nil {
		return err
	}
	for i := range tableColumns {
		tableColumns[i].Emb = embTccs[i]
	}

	for _, tcv := range tableColumnValues {
		tableColumnValueContents = append(tableColumnValueContents, tcv.DocContent)
	}
	embTccvs, err := k.documentData.TextEmbedding(tableColumnValueContents)
	if err != nil {
		return err
	}
	for i := range tableColumnValues {
		tableColumnValues[i].Emb = embTccvs[i]
	}
	k.LogInfof("开始存向量库，docId %v", doc.Id)
	err = es.DatabaseDocumentInsert(k.GetCtx(), kdb.GetIndexName(), tables, tableColumns, tableColumnValues, nil)
	if err != nil {
		k.LogErrorf("数据源索引构建失败：%s", err.Error())
		return components.ErrorDbIndexError
	}
	return
}
