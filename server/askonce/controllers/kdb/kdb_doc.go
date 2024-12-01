package kdb

import (
	"askonce/components/dto/dto_kdb_doc"
	"askonce/service"
	"github.com/xiangtao94/golib/flow"
)

type DocListController struct {
	flow.Controller
}

func (entity *DocListController) Action(req *dto_kdb_doc.ListReq) (interface{}, error) {
	s := entity.Create(new(service.KdbDocService)).(*service.KdbDocService)
	return s.DocList(req)
}

type DocAddController struct {
	flow.Controller
}

func (entity *DocAddController) Action(req *dto_kdb_doc.AddReq) (interface{}, error) {
	s := entity.Create(new(service.KdbDocService)).(*service.KdbDocService)
	return s.DocAdd(req)
}

type DocDeleteController struct {
	flow.Controller
}

func (entity *DocDeleteController) Action(req *dto_kdb_doc.DeleteReq) (interface{}, error) {
	s := entity.Create(new(service.KdbDocService)).(*service.KdbDocService)
	return s.DocDelete(req)
}

type DocRedoController struct {
	flow.Controller
}

func (entity *DocRedoController) Action(req *dto_kdb_doc.RedoReq) (interface{}, error) {
	s := entity.Create(new(service.KdbDocService)).(*service.KdbDocService)
	return s.DataRedo(req)
}

type RecallController struct {
	flow.Controller
}

func (entity *RecallController) Action(req *dto_kdb_doc.RecallReq) (interface{}, error) {
	s := entity.Create(new(service.SearchService)).(*service.SearchService)
	return s.Recall(req)
}
