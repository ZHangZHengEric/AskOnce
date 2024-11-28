package models

import (
	"askonce/components"
	"github.com/pkg/errors"
	"github.com/xiangtao94/golib/flow"
	"github.com/xiangtao94/golib/pkg/orm"

	"gorm.io/gorm"
	"time"
)

const (
	AuthTypeSuperAdmin = 9
	AuthTypeWrite      = 1
	AuthTypeRead       = 0
)

// KdbUser  用户知识库关系表
type KdbUser struct {
	Id       int64  `gorm:"id; primaryKey;autoIncrement;comment:自增主键"`
	KdbId    int64  `gorm:"type:int(11);default:0;comment:知识库id"`
	UserId   string `gorm:"type:varchar(128);default:'';comment:用户id"`
	AuthType int    `gorm:"type:int(11);default:0;comment:用户权限 0 可阅读 1 可编辑 9 超管"`
	orm.CrudModel
}

func (KdbUser) TableName() string {
	return "kdb_user"
}

type KdbUserDao struct {
	flow.Dao
}

func (entity *KdbUserDao) OnCreate() {
	entity.SetTable(KdbUser{}.TableName())
}

func (entity *KdbUserDao) Insert(add *KdbUser) (err error) {
	return entity.GetDB().Create(add).Error
}

// 更新
func (entity *KdbUserDao) Update(userId int64, update map[string]interface{}) error {
	update["updated_at"] = time.Now()
	err := entity.GetDB().Where("user_id = ?", userId).Updates(update).Error
	if err != nil {
		return components.ErrorMysqlError
	}
	return nil
}

func (entity *KdbUserDao) GetByKdbIdAndUserId(kdbId int64, userId string) (res *KdbUser, err error) {
	res = &KdbUser{}
	db := entity.GetDB()
	db = db.Table(entity.GetTable())
	err = db.Where("kdb_id = ? and user_id = ?", kdbId, userId).First(&res).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return
}

func (entity *KdbUserDao) GetByKdbId(kdbId int64) (res []*KdbUser, err error) {
	res = make([]*KdbUser, 0)
	db := entity.GetDB()
	db = db.Table(entity.GetTable())
	err = db.Where("kdb_id = ? ", kdbId).Find(&res).Error
	return
}

func (entity *KdbUserDao) GetByUserId(userId string) (res []*KdbUser, err error) {
	res = make([]*KdbUser, 0)
	db := entity.GetDB()
	db = db.Table(entity.GetTable())
	err = db.Where("user_id = ?", userId).Find(&res).Error
	return
}

func (entity *KdbUserDao) DeleteByKdbIdAndUserId(kdbId int64, userId string) (err error) {
	db := entity.GetDB()
	db = db.Table(entity.GetTable())
	err = db.Where("kdb_id = ? and user_id = ?", kdbId, userId).Delete(&KdbUser{}).Error
	return
}

func (entity *KdbUserDao) DeleteById(id int64) (err error) {
	db := entity.GetDB()
	db = db.Table(entity.GetTable())
	err = db.Where("id = ?", id).Delete(&KdbUser{}).Error
	return
}
