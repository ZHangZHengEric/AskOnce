package service

import (
	"askonce/components"
	"askonce/components/dto"
	"askonce/components/dto/dto_user"
	"askonce/data"
	"askonce/models"
	"askonce/utils"
	"github.com/xiangtao94/golib/flow"
)

type UserService struct {
	flow.Service
}

func (entity *UserService) LoginAccount(req *dto_user.LoginAccountReq) (res *dto_user.LoginRes, err error) {
	user, err := entity.Create(new(data.UserData)).(*data.UserData).LoginUserByAccount(req.Account, req.Password)
	if err != nil {
		return nil, err
	}
	res = &dto_user.LoginRes{
		UserId:          user.UserId,
		AtomechoSession: user.Session,
		Account:         req.Account,
	}
	return
}

func (entity *UserService) RegisterAccount(req *dto_user.RegisterAccountReq) (res *dto_user.LoginRes, err error) {
	userDao := entity.Create(new(models.UserDao)).(*models.UserDao)
	user, err := userDao.GetByAccount(req.Account)
	if err != nil {
		return nil, err
	}
	if user != nil {
		return nil, components.ErrorUserExistError
	}
	_, err = entity.Create(new(data.UserData)).(*data.UserData).RegisterUserAccount(req.Account, req.Password)
	if err != nil {
		return nil, err
	}
	return
}

func (entity *UserService) LoginInfo(req *dto.EmptyReq) (res *dto_user.LoginInfoRes, err error) {
	userInfo, err := utils.LoginInfo(entity.GetCtx())
	if err != nil {
		return nil, err
	}
	userDao := entity.Create(new(models.UserDao)).(*models.UserDao)
	user, err := userDao.GetByUserId(userInfo.UserId)
	if err != nil {
		return nil, err
	}
	res = &dto_user.LoginInfoRes{
		UserId:  user.UserId,
		Account: user.UserName,
	}
	return
}
