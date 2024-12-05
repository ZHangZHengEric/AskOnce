package search

import (
	"askonce/components/dto"
	"askonce/components/dto/dto_search"
	"askonce/service"
	"github.com/xiangtao94/golib/flow"
	"github.com/xiangtao94/golib/pkg/sse"
	"time"
)

type CaseController struct {
	flow.Controller
}

func (entity *CaseController) Action(req *dto_search.CaseReq) (interface{}, error) {
	s := entity.Create(new(service.SearchService)).(*service.SearchService)
	return s.Case(req)
}

type SessionController struct {
	flow.Controller
}

func (entity *SessionController) Action(req *dto.EmptyReq) (interface{}, error) {
	s := entity.Create(new(service.SearchService)).(*service.SearchService)
	return s.Session(req)
}

type AskController struct {
	flow.Controller
}

func (entity *AskController) ShouldRender() bool {
	return false
}

func (entity *AskController) Action(req *dto_search.AskReq) (interface{}, error) {
	s := entity.Create(new(service.SearchService)).(*service.SearchService)
	err := s.Ask(req)
	if err != nil {
		time.Sleep(50 * time.Millisecond)
		sse.RenderStreamError(entity.GetCtx(), err)
	}
	return nil, nil
}

type ChatAskSyncController struct {
	flow.Controller
}

func (entity *ChatAskSyncController) Action(req *dto_search.ChatAskReq) (interface{}, error) {
	s := entity.Create(new(service.SearchService)).(*service.SearchService)
	return s.AskSync(req)
}

type HisController struct {
	flow.Controller
}

func (entity *HisController) Action(req *dto_search.HisReq) (interface{}, error) {
	s := entity.Create(new(service.SearchService)).(*service.SearchService)
	return s.History(req)
}

type ReferController struct {
	flow.Controller
}

func (entity *ReferController) Action(req *dto_search.ReferReq) (interface{}, error) {
	s := entity.Create(new(service.SearchService)).(*service.SearchService)
	return s.Reference(req)
}

type OutlineController struct {
	flow.Controller
}

func (entity *OutlineController) Action(req *dto_search.OutlineReq) (interface{}, error) {
	s := entity.Create(new(service.SearchService)).(*service.SearchService)
	return s.Outline(req)
}

type UnlikeController struct {
	flow.Controller
}

func (entity *UnlikeController) Action(req *dto_search.UnlikeReq) (interface{}, error) {
	s := entity.Create(new(service.SearchService)).(*service.SearchService)
	return s.Unlike(req)
}

type RelationController struct {
	flow.Controller
}

func (entity *RelationController) Action(req *dto_search.RelationReq) (interface{}, error) {
	s := entity.Create(new(service.SearchService)).(*service.SearchService)
	return s.Relation(req)
}

type ProcessController struct {
	flow.Controller
}

func (entity *ProcessController) Action(req *dto_search.ProcessReq) (interface{}, error) {
	s := entity.Create(new(service.SearchService)).(*service.SearchService)
	return s.Process(req)
}

type KdbListController struct {
	flow.Controller
}

func (entity *KdbListController) Action(req *dto_search.KdbListReq) (interface{}, error) {
	s := entity.Create(new(service.SearchService)).(*service.SearchService)
	return s.KdbList(req)
}
