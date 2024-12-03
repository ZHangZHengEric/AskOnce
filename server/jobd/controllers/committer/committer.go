package committer

import (
	"github.com/xiangtao94/golib/flow"
	"jobd/components/dto"
	"jobd/service"
)

type DoTaskCtl struct {
	flow.Controller
}

func (entity *DoTaskCtl) Action(req *dto.DoTaskReq) (res interface{}, err error) {
	s := flow.Create(entity.GetCtx(), new(service.CommitterService))
	return s.DoTask(req)
}

type DoTaskStreamCtl struct {
	flow.Controller
}

func (entity *DoTaskStreamCtl) ShouldRender() bool {
	return false
}

func (entity *DoTaskStreamCtl) Action(req *dto.DoTaskReq) (res interface{}, err error) {
	s := flow.Create(entity.GetCtx(), new(service.CommitterService))
	err = s.DoTaskStream(req)
	return nil, err
}

type CommitCtl struct {
	flow.Controller
}

func (entity *CommitCtl) Action(req *dto.CommitReq) (res interface{}, err error) {
	s := flow.Create(entity.GetCtx(), new(service.CommitterService))
	return s.Commit(req)
}

type GetInfoCtl struct {
	flow.Controller
}

func (entity *GetInfoCtl) Action(req *dto.GetInfoReq) (res interface{}, err error) {
	s := flow.Create(entity.GetCtx(), new(service.CommitterService))
	return s.GetInfo(req)
}

type BlockGetInfoCtl struct {
	flow.Controller
}

func (entity *BlockGetInfoCtl) Action(req *dto.GetInfoReq) (res interface{}, err error) {
	s := flow.Create(entity.GetCtx(), new(service.CommitterService))
	return s.BlockGetInfo(req)
}
