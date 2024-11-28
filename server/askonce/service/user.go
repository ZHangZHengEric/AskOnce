package service

import (
	"askonce/components"
	"askonce/components/defines"
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
	// 3.种缓存
	sessionCache := entity.Create(new(data.SessionCache)).(*data.SessionCache)
	session := data.LoginInfoSession{
		UserId:    user.UserId,
		Account:   req.Account,
		LoginTime: user.LastLoginTime.Unix(),
	}
	userSessionKey := sessionCache.GenSessionValue()
	err = sessionCache.SetSession(session, userSessionKey, defines.COOKIE_DEFAULT_AGE)
	if err != nil {
		return nil, err
	}
	// set cookie
	sessionCache.AddUserSessionSet(session.UserId, userSessionKey, int64(defines.COOKIE_DEFAULT_AGE))
	domain := utils.GetCookieDomain(entity.GetCtx().Request.Host)
	entity.LogInfof("set cookie domain:%+v", domain)
	entity.GetCtx().SetCookie(defines.COOKIE_KEY, userSessionKey, defines.COOKIE_DEFAULT_AGE, defines.COOKIE_PATH, domain, false, false)
	res = &dto_user.LoginRes{
		UserId:          user.UserId,
		AtomechoSession: userSessionKey,
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
