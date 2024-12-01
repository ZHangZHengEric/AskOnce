package models

import (
	"github.com/xiangtao94/golib/flow"
	"github.com/xiangtao94/golib/pkg/orm"
)

type KdbCover struct {
	Id           int64  `gorm:"id; primaryKey;autoIncrement;comment:自增主键"`
	Type         string `gorm:"type:varchar(52);default:'';comment:类型"`
	Url          string `gorm:"type:varchar(512);default:'';comment:url"`
	DefaultColor bool   `gorm:"type:tinyint(1);default:0;comment:默认颜色"`
	UserId       string `gorm:"type:varchar(128);default:'';comment:用户id"`
	orm.CrudModel
}

func (KdbCover) TableName() string {
	return "kdb_cover"
}

type KdbCoverDao struct {
	flow.Dao
}

func (entity *KdbCoverDao) OnCreate() {
	entity.SetTable(KdbCover{}.TableName())
}

func (entity *KdbCoverDao) Insert(add *KdbCover) (err error) {
	return entity.GetDB().Create(add).Error
}

func (entity *KdbCoverDao) GetAll() (res []*KdbCover, err error) {
	db := entity.GetDB()
	err = db.Where("1 = 1").Find(&res).Error
	return
}

func (entity *KdbCoverDao) GetById(id int64) (res *KdbCover, err error) {
	res = &KdbCover{}
	db := entity.GetDB()
	db = db.Table(entity.GetTable())
	err = db.Where("id = ?", id).First(&res).Error
	return
}
