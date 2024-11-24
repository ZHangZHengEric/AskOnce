package controllers

import (
	"dispatch/components/dto"
	"dispatch/service"
	"github.com/xiangtao94/golib/flow"
)

type AddTaskTypeInfoCtl struct {
	flow.Controller
}

func (entity *AddTaskTypeInfoCtl) Action(req *dto.AddTaskTypeInfoReq) (res interface{}, err error) {
	s := flow.Create(entity.GetCtx(), new(service.TaskManageService))
	return s.AddTaskTypeInfo(req)
}
