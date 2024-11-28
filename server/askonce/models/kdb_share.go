package models

import (
	"github.com/pkg/errors"
	"github.com/xiangtao94/golib/flow"
	"github.com/xiangtao94/golib/pkg/orm"
	"gorm.io/gorm"
)

type KdbShare struct {
	Id        int64  `gorm:"id; primaryKey;autoIncrement;comment:自增主键"`
	ShareCode string ` gorm:"type:varchar(52);default:'';comment:分享code"`
	KdbId     int64  `gorm:"type:int(11);default:0;comment:知识库id"`
	UserId    string `gorm:"type:varchar(128);default:'';comment:用户id"`
	AuthType  int    `gorm:"type:int(11);default:0;comment:用户权限 0 可阅读 1 可编辑 9 超管"`
	orm.CrudModel
}

func (KdbShare) TableName() string {
	return "kdb_share"
}

type KdbShareDao struct {
	flow.Dao
}

func (entity *KdbShareDao) OnCreate() {
	entity.SetTable(KdbShare{}.TableName())
}

func (entity *KdbShareDao) Insert(user *KdbShare) (err error) {
	return entity.GetDB().Create(user).Error
}

func (entity *KdbShareDao) GetByShareCode(shareCode string) (res *KdbShare, err error) {
	db := entity.GetDB()
	db = db.Table(entity.GetTable())
	err = db.Where("share_code = ? ", shareCode).First(&res).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return
}
