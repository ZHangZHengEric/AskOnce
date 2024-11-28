package data

import (
	"askonce/components/dto"
	"askonce/models"
	"github.com/xiangtao94/golib/flow"
)

type UserData struct {
	flow.Service

	userDao *models.UserDao
}

func (u *UserData) OnCreate() {
	u.userDao = flow.Create(u.GetCtx(), new(models.UserDao))
}

func (u *UserData) QueryUserList(queryUserIds []string, queryName string, param dto.PageParam) (list []*models.User, cnt int64, err error) {
	return
}
