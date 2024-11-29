package kdb

import (
	"askonce/components/dto"
	"askonce/components/dto/dto_kdb"
	"askonce/service"
	"github.com/xiangtao94/golib/flow"
)

type ListController struct {
	flow.Controller
}

func (entity *ListController) Action(req *dto_kdb.ListReq) (interface{}, error) {
	s := entity.Create(new(service.KdbService)).(*service.KdbService)
	return s.List(req)
}

type AddController struct {
	flow.Controller
}

func (entity *AddController) Action(req *dto_kdb.AddReq) (interface{}, error) {
	s := entity.Create(new(service.KdbService)).(*service.KdbService)
	return s.Add(req)
}

type UpdateController struct {
	flow.Controller
}

func (entity *UpdateController) Action(req *dto_kdb.UpdateReq) (interface{}, error) {
	s := entity.Create(new(service.KdbService)).(*service.KdbService)
	return s.Update(req)
}

type DeleteController struct {
	flow.Controller
}

func (entity *DeleteController) Action(req *dto_kdb.DeleteReq) (interface{}, error) {
	s := entity.Create(new(service.KdbService)).(*service.KdbService)
	return s.Delete(req)
}

type InfoController struct {
	flow.Controller
}

func (entity *InfoController) Action(req *dto_kdb.InfoReq) (interface{}, error) {
	s := entity.Create(new(service.KdbService)).(*service.KdbService)
	return s.Info(req)
}

type DeleteSelfController struct {
	flow.Controller
}

func (entity *DeleteSelfController) Action(req *dto_kdb.DeleteSelfReq) (interface{}, error) {
	s := entity.Create(new(service.KdbService)).(*service.KdbService)
	return s.DeleteRelation(req)
}

type CoversController struct {
	flow.Controller
}

func (entity *CoversController) Action(req *dto.EmptyReq) (interface{}, error) {
	s := entity.Create(new(service.KdbService)).(*service.KdbService)
	return s.Covers(req)
}

type AuthController struct {
	flow.Controller
}

func (entity *AuthController) Action(req *dto_kdb.AuthReq) (interface{}, error) {
	s := entity.Create(new(service.KdbService)).(*service.KdbService)
	return s.Auth(req)
}

type UserListController struct {
	flow.Controller
}

func (entity *UserListController) Action(req *dto_kdb.UserListReq) (interface{}, error) {
	s := entity.Create(new(service.KdbService)).(*service.KdbService)
	return s.UserList(req)
}

type UserQueryController struct {
	flow.Controller
}

func (entity *UserQueryController) Action(req *dto_kdb.UserQueryReq) (interface{}, error) {
	s := entity.Create(new(service.KdbService)).(*service.KdbService)
	return s.UserQuery(req)
}

type UserAddController struct {
	flow.Controller
}

func (entity *UserAddController) Action(req *dto_kdb.UserAddReq) (interface{}, error) {
	s := entity.Create(new(service.KdbService)).(*service.KdbService)
	return s.UserAdd(req)
}

type UserDeleteController struct {
	flow.Controller
}

func (entity *UserDeleteController) Action(req *dto_kdb.UserDeleteReq) (interface{}, error) {
	s := entity.Create(new(service.KdbService)).(*service.KdbService)
	return s.UserDelete(req)
}

type GenShareCodeController struct {
	flow.Controller
}

func (entity *GenShareCodeController) Action(req *dto_kdb.GenShareCodeReq) (interface{}, error) {
	s := entity.Create(new(service.KdbService)).(*service.KdbService)
	return s.GenShareCode(req)
}

type VerifyShareCodeController struct {
	flow.Controller
}

func (entity *VerifyShareCodeController) Action(req *dto_kdb.VerifyShareCodeReq) (interface{}, error) {
	s := entity.Create(new(service.KdbService)).(*service.KdbService)
	return s.VerifyShareCode(req)
}

type ShareCodeInfoController struct {
	flow.Controller
}

func (entity *ShareCodeInfoController) Action(req *dto_kdb.InfoShareCodeReq) (interface{}, error) {
	s := entity.Create(new(service.KdbService)).(*service.KdbService)
	return s.ShareCodeInfo(req)
}
