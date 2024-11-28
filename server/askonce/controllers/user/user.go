package user

import (
	"askonce/components/dto"
	"askonce/components/dto/dto_user"
	"askonce/service"
	"github.com/xiangtao94/golib/flow"
)

type LoginAccountController struct {
	flow.Controller
}

func (entity *LoginAccountController) Action(req *dto_user.LoginAccountReq) (interface{}, error) {
	s := entity.Create(new(service.UserService)).(*service.UserService)
	return s.LoginAccount(req)
}

type RegisterAccountController struct {
	flow.Controller
}

func (entity *RegisterAccountController) Action(req *dto_user.RegisterAccountReq) (interface{}, error) {
	s := entity.Create(new(service.UserService)).(*service.UserService)
	return s.RegisterAccount(req)
}

type LoginInfoController struct {
	flow.Controller
}

func (entity *LoginInfoController) Action(req *dto.EmptyReq) (interface{}, error) {
	s := entity.Create(new(service.UserService)).(*service.UserService)
	return s.LoginInfo(req)
}
