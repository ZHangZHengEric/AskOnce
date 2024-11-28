package api

import (
	"askonce/service"
	"github.com/xiangtao94/golib/flow"
)

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

type ListController struct {
	flow.Controller
}

func (entity *ListController) Action(req *dto_kdb.ListReq) (interface{}, error) {
	s := entity.Create(new(service.KdbService)).(*service.KdbService)
	return s.List(req)
}

type InfoController struct {
	flow.Controller
}

func (entity *InfoController) Action(req *dto_kdb.InfoReq) (interface{}, error) {
	s := entity.Create(new(service.KdbService)).(*service.KdbService)
	return s.Info(req)
}
