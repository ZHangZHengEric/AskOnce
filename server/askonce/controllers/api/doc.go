package api

import (
	"askonce/components/dto"
	"github.com/xiangtao94/golib/flow"
)

type DocProcessRuleController struct {
	flow.Controller
}

func (entity *DocProcessRuleController) Action(req *dto.EmptyReq) (interface{}, error) {
	s := entity.Create(new(service.KdbDocService)).(*service.KdbDocService)
	return s.GetProcessRule(req)
}

type DocProcessReviewController struct {
	flow.Controller
}

func (entity *DocProcessReviewController) Action(req *dto_kdb.ProcessReviewReq) (interface{}, error) {
	s := entity.Create(new(service.KdbDocService)).(*service.KdbDocService)
	return s.ProcessReview(req)
}

type DocAddController struct {
	flow.Controller
}

func (entity *DocAddController) Action(req *dto_kdb.DocAddReq) (interface{}, error) {
	s := entity.Create(new(service.KdbDocService)).(*service.KdbDocService)
	return s.DocAdd(req)
}

type DocInfoController struct {
	flow.Controller
}

func (entity *DocInfoController) Action(req *dto_kdb.DocInfoReq) (interface{}, error) {
	s := entity.Create(new(service.KdbDocService)).(*service.KdbDocService)
	return s.DocInfo(req)
}

type DocListController struct {
	flow.Controller
}

func (entity *DocListController) Action(req *dto_kdb.DocListReq) (interface{}, error) {
	s := entity.Create(new(service.KdbDocService)).(*service.KdbDocService)
	return s.DocList(req)
}

type DocDeleteController struct {
	flow.Controller
}

func (entity *DocDeleteController) Action(req *dto_kdb.DocDeleteReq) (interface{}, error) {
	s := entity.Create(new(service.KdbDocService)).(*service.KdbDocService)
	return s.DocDelete(req)
}

type DocRenameController struct {
	flow.Controller
}

func (entity *DocRenameController) Action(req *dto_kdb.DocRenameReq) (interface{}, error) {
	s := entity.Create(new(service.KdbDocService)).(*service.KdbDocService)
	return s.DocRename(req)
}

type DocEnableController struct {
	flow.Controller
}

func (entity *DocEnableController) Action(req *dto_kdb.DocStatusSettingReq) (interface{}, error) {
	s := entity.Create(new(service.KdbDocService)).(*service.KdbDocService)
	return s.DocEnable(req)
}

type DocDisableController struct {
	flow.Controller
}

func (entity *DocDisableController) Action(req *dto_kdb.DocStatusSettingReq) (interface{}, error) {
	s := entity.Create(new(service.KdbDocService)).(*service.KdbDocService)
	return s.DocDisable(req)
}

type DocSegmentController struct {
	flow.Controller
}

func (entity *DocSegmentController) Action(req *dto_kdb.DocSegmentReq) (interface{}, error) {
	s := entity.Create(new(service.KdbDocService)).(*service.KdbDocService)
	return s.DocSegments(req)
}
