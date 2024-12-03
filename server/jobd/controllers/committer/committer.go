package committer

import (
	"encoding/json"
	"fmt"
	"github.com/xiangtao94/golib/flow"
	"jobd/components/defines"
	"jobd/components/dto"
	"jobd/service"
	"net/http"
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
	if err != nil {
		resp := dto.DoTaskResp{}
		resp.Output = err.Error()
		resp.Status = defines.STATUS_EXEC_FAILED
		resp.SessionId = req.SessionId
		resp.TaskType = req.TaskType
		str, _ := json.Marshal(resp)
		flusher, _ := entity.GetCtx().Writer.(http.Flusher)
		fmt.Fprintf(entity.GetCtx().Writer, "%s\n", str)
		flusher.Flush()
	}
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
