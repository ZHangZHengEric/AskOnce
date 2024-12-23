package models

import (
	"askonce/components"
	"askonce/components/dto"
	"github.com/pkg/errors"
	"github.com/xiangtao94/golib/flow"
	"github.com/xiangtao94/golib/pkg/orm"
	"gorm.io/datatypes"

	"gorm.io/gorm"

	"time"
)

const (
	KdbDocWaiting = 0
	KdbDocRunning = 1
	KdbDocFail    = 2
	KdbDocSuccess = 9
)

const (
	DataSourceFile     = "file"
	DataSourceDatabase = "database"
)

// KdbDoc  知识库文档
type KdbDoc struct {
	Id         int64             `gorm:"id; primaryKey;autoIncrement;comment:自增主键"`
	KdbId      int64             `gorm:"type:int(11);default:0;comment:知识库id"`
	TaskId     string            `gorm:"type:varchar(128);default:'';comment:任务id"`
	DocName    string            `gorm:"type:varchar(128);default:'';comment:文档名称"`
	DataSource string            `gorm:"type:varchar(52);default:'';comment:文档来源 file"`
	SourceId   string            `gorm:"type:varchar(128);default:0;comment:来源id"`
	Status     int               `gorm:"type:int(11);default:0;comment: 状态 0 初始化，没有导入，1 正在处理， 2 导入失败 9 导入成功"`
	RetryCount int               `gorm:"type:int(11);default:0;comment:重试次数"`
	UserId     string            `gorm:"type:varchar(128);default:'';comment:上传用户id"`
	Metadata   datatypes.JSONMap `gorm:"type:json;comment:元数据信息"`
	orm.CrudModel
}

func (KdbDoc) TableName() string {
	return "kdb_doc"
}

type KdbDocDao struct {
	flow.Dao
}

func (entity *KdbDocDao) OnCreate() {
	entity.SetTable(KdbDoc{}.TableName())
}

func (entity *KdbDocDao) Insert(add *KdbDoc) (err error) {
	return entity.GetDB().Create(add).Error
}

func (entity *KdbDocDao) BatchInsert(add []*KdbDoc) (err error) {
	if len(add) == 0 {
		return nil
	}
	return entity.GetDB().CreateInBatches(add, 2000).Error
}

// 更新
func (entity *KdbDocDao) Update(id int64, update map[string]interface{}) error {
	update["updated_at"] = time.Now()
	err := entity.GetDB().Where("id = ?", id).Updates(update).Error
	if err != nil {
		return components.ErrorMysqlError
	}
	return nil
}

// 更新状态
func (entity *KdbDocDao) UpdateStatus(doc *KdbDoc, status int) error {
	doc.Status = status
	doc.UpdatedAt = time.Now()
	if status == KdbDocFail {
		doc.RetryCount = doc.RetryCount + 1
	}
	err := entity.GetDB().Where("id = ?", doc.Id).Updates(doc).Error
	if err != nil {
		return components.ErrorMysqlError
	}
	return nil
}

func (entity *KdbDocDao) GetById(id int64) (res *KdbDoc, err error) {
	res = &KdbDoc{}
	db := entity.GetDB()
	db = db.Table(entity.GetTable())
	err = db.Where("id = ?", id).First(&res).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return
}

func (entity *KdbDocDao) GetByIds(ids []int64) (res []*KdbDoc, err error) {
	res = []*KdbDoc{}
	if len(ids) == 0 {
		return
	}
	db := entity.GetDB()
	db = db.Table(entity.GetTable())
	err = db.Where("id in ?", ids).Find(&res).Error
	return
}
func (entity *KdbDocDao) GetListByStatus(status int) (res []*KdbDoc, err error) {
	res = []*KdbDoc{}
	db := entity.GetDB()
	db = db.Table(entity.GetTable())
	err = db.Where("status = ?", status).Limit(10).Find(&res).Error
	return
}

func (entity *KdbDocDao) GetFailedList() (res []*KdbDoc, err error) {
	res = []*KdbDoc{}
	db := entity.GetDB()
	db = db.Table(entity.GetTable())
	err = db.Where("status = ? and retry_count < 3", KdbDocFail).Limit(10).Find(&res).Error
	return
}

func (entity *KdbDocDao) GetList(kdbId int64, queryName string, status []int, param dto.PageParam) (list []*KdbDoc, cnt int64, err error) {
	db := entity.GetDB().Model(&KdbDoc{})
	db = db.Where("kdb_id =   ?", kdbId)
	if len(queryName) > 0 {
		db = db.Where("doc_name like ?", "%"+queryName+"%")
	}
	if len(status) > 0 {
		db = db.Where("status in (?)", status)
	}
	db = db.Count(&cnt)
	if err = db.Offset((param.PageNo - 1) * param.PageSize).Limit(param.PageSize).Order("created_at desc").Find(&list).Error; err != nil {
		entity.LogErrorf("KdbDocDao GetList err:%s", err.Error())
		return nil, 0, components.ErrorMysqlError
	}
	return
}

func (entity *KdbDocDao) DeleteById(docId int64) (err error) {
	db := entity.GetDB()
	db = db.Table(entity.GetTable())
	err = db.Where("id =  ? ", docId).Delete(&KdbDoc{}).Error
	return err
}

func (entity *KdbDocDao) DeleteByIds(docIds []int64) (err error) {
	if len(docIds) == 0 {
		return
	}
	db := entity.GetDB()
	db = db.Table(entity.GetTable())
	err = db.Where("id in (?) ", docIds).Delete(&KdbDoc{}).Error
	return err
}

type Progress struct {
	Status int   `json:"status"`
	Total  int64 `json:"total"`
}

func (entity *KdbDocDao) QueryProcess(kdbId int64, taskId string) (res []*Progress, err error) {
	db := entity.GetDB()
	db = db.Table(entity.GetTable()).Model(&KdbDoc{})
	db = db.Where("kdb_id = ? and task_id = ?", kdbId, taskId)
	err = db.Select("status,count(1) as total").Group("status").Find(&res).Error
	return
}

func (entity *KdbDocDao) GetByTaskIdAndStatus(kdbId int64, taskId string, status []int) (res []*KdbDoc, err error) {
	db := entity.GetDB()
	db = db.Table(entity.GetTable()).Model(&KdbDoc{})
	err = db.Where("kdb_id = ? and task_id = ? and status in (?)", kdbId, taskId, status).Find(&res).Error
	return
}

func (entity *KdbDocDao) BatchUpdateStatus(ids []int64, status int) (err error) {
	if len(ids) == 0 {
		return nil
	}
	db := entity.GetDB()
	db = db.Table(entity.GetTable())
	db = db.Where("id in (?)", ids).Updates(map[string]interface{}{"status": status, "updated_at": time.Now()})
	return db.Error
}

func (entity *KdbDocDao) GetByKdbId(kdbId int64) (res []*KdbDoc, err error) {
	db := entity.GetDB()
	db = db.Table(entity.GetTable()).Model(&KdbDoc{})
	db = db.Where("kdb_id = ?", kdbId)
	err = db.Find(&res).Error
	return
}
