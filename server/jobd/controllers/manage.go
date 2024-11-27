package controllers

import (
	"github.com/xiangtao94/golib/flow"
	"jobd/components/dto"
	"jobd/service"
)

type AddTaskTypeInfoCtl struct {
	flow.Controller
}

func (entity *AddTaskTypeInfoCtl) Action(req *dto.AddTaskTypeInfoReq) (res interface{}, err error) {
	s := flow.Create(entity.GetCtx(), new(service.TaskManageService))
	return s.AddTaskTypeInfo(req)
}
