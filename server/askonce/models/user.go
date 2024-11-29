package models

import (
	"askonce/components"
	"askonce/components/dto"
	"errors"
	"github.com/xiangtao94/golib/flow"
	"github.com/xiangtao94/golib/pkg/orm"
	"gorm.io/datatypes"

	"gorm.io/gorm"
	"time"
)

// User  用户信息定义表
type User struct {
	Id            int64                           `gorm:"id; primaryKey;autoIncrement;comment:自增主键"`
	UserId        string                          `gorm:"type:varchar(128);default:'';comment:用户id"`
	UserName      string                          `gorm:"type:varchar(512);default:'';comment:登陆账户名"`
	Password      string                          `gorm:"type:varchar(512);default:'';comment:登陆密码"`
	Setting       datatypes.JSONType[UserSetting] `gorm:"comment:用户配置"`
	LastLoginTime time.Time                       `gorm:"comment:最后登陆时间"`
	orm.CrudModel
}

func (User) TableName() string {
	return "user"
}

type UserSetting struct {
	Language  string `json:"language"`
	ModelType string ` json:"modelType"`
}

type UserDao struct {
	flow.Dao
}

func (entity *UserDao) OnCreate() {
	entity.SetTable(User{}.TableName())
}

func (entity *UserDao) Insert(user *User) (err error) {
	return entity.GetDB().Create(user).Error
}

// 更新
func (entity *UserDao) Update(userId string, update map[string]interface{}) error {
	update["updated_at"] = time.Now()
	err := entity.GetDB().Where("user_id = ?", userId).Updates(update).Error
	if err != nil {
		entity.LogErrorf("OpenSecretDao update error, teamId:%s, err:%s", userId, err)
		return components.ErrorMysqlError
	}
	return nil
}

func (entity *UserDao) GetByUserId(userId string) (res *User, err error) {
	res = &User{}
	db := entity.GetDB()
	db = db.Table(entity.GetTable())
	err = db.Where("user_id = ?", userId).First(&res).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return
}

func (entity *UserDao) GetList(userIds []string, queryName string, param dto.PageParam) (list []*User, cnt int64, err error) {
	db := entity.GetDB().Model(&User{})
	if len(userIds) > 0 {
		db = db.Where("user_id in ?", userIds)
	}
	if len(queryName) > 0 {
		db = db.Where("(phone like ? or account like ?)", "%"+queryName+"%", "%"+queryName+"%")
	}
	db = db.Count(&cnt)
	if err = db.Offset((param.PageNo - 1) * param.PageSize).Limit(param.PageSize).Order("created_at desc").Find(&list).Error; err != nil {
		entity.LogErrorf("UserDao GetList err:%s", err.Error())
		return nil, 0, components.ErrorMysqlError
	}
	return
}

func (entity *UserDao) GetByAccount(userName string) (res *User, err error) {
	db := entity.GetDB()
	err = db.Where("user_name = ?", userName).First(&res).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return
}

func (entity *UserDao) GetByUserIds(userIds []string) (res []*User, err error) {
	if len(userIds) == 0 {
		return
	}
	db := entity.GetDB()
	err = db.Where("user_id in ?", userIds).Find(&res).Error
	return
}
