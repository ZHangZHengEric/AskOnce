package models

import (
	"askonce/components"
	"askonce/components/dto"
	"fmt"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/pkg/errors"
	"github.com/xiangtao94/golib/flow"
	"github.com/xiangtao94/golib/pkg/orm"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"time"
)

const (
	KdbTypePublic  = "public"
	KdbTypePrivate = "private"
)

const (
	DataTypeDoc = "doc"      // 通用文件库
	DataTypeEml = "eml"      // 邮件库
	DataTypeDB  = "database" // sql数据库
)

// Kdb  知识库表
type Kdb struct {
	Id       int64                              `gorm:"id; primaryKey;autoIncrement;comment:自增主键"`
	Name     string                             `gorm:"type:varchar(128);default:'';comment:知识库名"`
	Intro    string                             `gorm:"type:varchar(1024);default:'';comment:介绍"`
	Setting  datatypes.JSONType[dto.KdbSetting] `gorm:"type:json;comment:知识库设置"`
	DataType string                             `gorm:"type:varchar(52);default:'';comment:数据类型 common"`
	Type     string                             `gorm:"type:varchar(52);default:'';comment:类型  private 私有知识库 public 共有知识库"`
	Creator  string                             `gorm:"type:varchar(128);default:'';comment: 创建人id"`
	orm.CrudModel
}

func (k Kdb) TableName() string {
	return "kdb"
}

func (k Kdb) GetIndexName() string {
	return fmt.Sprintf("askonce_%v", k.Id)
}

type KdbDao struct {
	flow.Dao
}

func (entity *KdbDao) OnCreate() {
	entity.SetTable(Kdb{}.TableName())
}

func (entity *KdbDao) Insert(add *Kdb) (err error) {
	return entity.GetDB().Create(add).Error
}

// 更新
func (entity *KdbDao) Update(id int64, update map[string]interface{}) error {
	update["updated_at"] = time.Now()
	err := entity.GetDB().Where("id = ?", id).Updates(update).Error
	if err != nil {
		return components.ErrorMysqlError
	}
	return nil
}

func (entity *KdbDao) GetById(id int64) (res *Kdb, err error) {
	res = &Kdb{}
	db := entity.GetDB()
	db = db.Table(entity.GetTable())
	err = db.Where("id = ?", id).First(&res).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return
}

func (entity *KdbDao) GetByNameAndCreator(name string, creator string) (res *Kdb, err error) {
	res = &Kdb{}
	db := entity.GetDB()
	db = db.Table(entity.GetTable())
	err = db.Where("name = ? and creator = ? ", name, creator).First(&res).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return
}

func (entity *KdbDao) GetByIds(ids []int64) (res []*Kdb, err error) {
	res = []*Kdb{}
	if len(ids) == 0 {
		return
	}
	ids = slice.Unique(ids)
	db := entity.GetDB()
	db = db.Table(entity.GetTable())
	err = db.Where("id in ?", ids).Find(&res).Error
	return
}
func (entity *KdbDao) GetList(kdbIds []int64, queryName string, param dto.PageParam) (list []*Kdb, cnt int64, err error) {
	db := entity.GetDB().Model(&Kdb{})
	db = db.Where("id in (?)", kdbIds)
	if len(queryName) > 0 {
		db = db.Where("name like ?", "%"+queryName+"%")
	}
	db = db.Count(&cnt)
	if err = db.Offset((param.PageNo - 1) * param.PageSize).Limit(param.PageSize).Order("created_at desc").Find(&list).Error; err != nil {
		entity.LogErrorf("KdbDao GetList err:%s", err.Error())
		return nil, 0, components.ErrorMysqlError
	}
	return
}

func (entity *KdbDao) GetPubIds() (res []int64, err error) {
	db := entity.GetDB()
	err = db.Model(&Kdb{}).Where("type = ?", KdbTypePublic).Select("id").Scan(&res).Error
	return
}

func (entity *KdbDao) GetPub() (res []*Kdb, err error) {
	db := entity.GetDB()
	err = db.Where("type = ?", KdbTypePublic).Find(&res).Error
	return
}

func (entity *KdbDao) DeleteById(id int64) (err error) {
	db := entity.GetDB()
	db = db.Table(entity.GetTable())
	err = db.Where("id = ?", id).Delete(&Kdb{}).Error
	return
}
