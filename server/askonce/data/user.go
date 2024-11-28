package data

import (
	"askonce/components"
	"askonce/components/dto"
	"askonce/helpers"
	"askonce/models"
	"encoding/base64"
	"github.com/xiangtao94/golib/flow"
	"github.com/xiangtao94/golib/pkg/orm"
	"gorm.io/datatypes"
	"time"
)

type UserData struct {
	flow.Service

	userDao *models.UserDao
}

func (u *UserData) OnCreate() {
	u.userDao = flow.Create(u.GetCtx(), new(models.UserDao))
}

func (u *UserData) QueryUserList(queryUserIds []string, queryName string, param dto.PageParam) (list []*models.User, cnt int64, err error) {
	list, cnt, err = u.userDao.GetList(queryUserIds, queryName, param)
	if err != nil {
		return
	}
	return
}

func (u *UserData) LoginUserByAccount(account string, password string) (*models.User, error) {
	user, err := u.userDao.GetByAccount(account)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, components.ErrorLoginError
	}
	if user.Password != base64.StdEncoding.EncodeToString([]byte(password)) {
		return nil, components.ErrorLoginError
	}
	now := time.Now()
	user.LastLoginTime = now
	err = u.userDao.Update(user.UserId, map[string]interface{}{"last_login_time": now})
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserData) RegisterUserAccount(account string, password string) (newUser *models.User, err error) {
	now := time.Now()
	userId := helpers.GenUserID()
	db := helpers.MysqlClient.WithContext(u.GetCtx())
	tx := db.Begin()
	u.userDao.SetDB(tx)
	insertUser := &models.User{
		UserId:   userId,
		Password: base64.StdEncoding.EncodeToString([]byte(password)),
		UserName: account,
		Setting: datatypes.NewJSONType(models.UserSetting{
			Language:  "zh-cn",
			ModelType: "chatgpt-4-minio",
		}),
		LastLoginTime: now,
		CrudModel: orm.CrudModel{
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
	if err = u.userDao.Insert(insertUser); err != nil {
		tx.Rollback()
		return
	}
	tx.Commit()
	return insertUser, nil
}
