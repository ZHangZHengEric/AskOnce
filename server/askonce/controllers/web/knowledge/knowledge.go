package knowledge

import (
	"askonce/components/dto"
	"askonce/components/dto/dto_knowledge"
	"github.com/xiangtao94/golib/flow"
)

type ListController struct {
	flow.Controller
}

func (entity *ListController) Action(req *dto_knowledge.ListReq) (interface{}, error) {
	s := entity.Create(new(service.KnowledgeService)).(*service.KnowledgeService)
	return s.List(req)
}

type AddController struct {
	flow.Controller
}

func (entity *AddController) Action(req *dto_knowledge.AddReq) (interface{}, error) {
	s := entity.Create(new(service.KnowledgeService)).(*service.KnowledgeService)
	return s.Add(req)
}

type UpdateController struct {
	flow.Controller
}

func (entity *UpdateController) Action(req *dto_knowledge.UpdateReq) (interface{}, error) {
	s := entity.Create(new(service.KnowledgeService)).(*service.KnowledgeService)
	return s.Update(req)
}

type DeleteController struct {
	flow.Controller
}

func (entity *DeleteController) Action(req *dto_knowledge.DeleteReq) (interface{}, error) {
	s := entity.Create(new(service.KnowledgeService)).(*service.KnowledgeService)
	return s.Delete(req)
}

type DataListController struct {
	flow.Controller
}

func (entity *DataListController) Action(req *dto_knowledge.DataListReq) (interface{}, error) {
	s := entity.Create(new(service.KnowledgeService)).(*service.KnowledgeService)
	return s.DataList(req)
}

type DataAddController struct {
	flow.Controller
}

func (entity *DataAddController) Action(req *dto_knowledge.DataAddReq) (interface{}, error) {
	s := entity.Create(new(service.KnowledgeService)).(*service.KnowledgeService)
	return s.DataAdd(req)
}

type DataBatchAddController struct {
	flow.Controller
}

func (entity *DataBatchAddController) Action(req *dto_knowledge.DataBatchAddReq) (interface{}, error) {
	s := entity.Create(new(service.KnowledgeService)).(*service.KnowledgeService)
	return s.DataBatchAdd(req)
}

type DataDeleteController struct {
	flow.Controller
}

func (entity *DataDeleteController) Action(req *dto_knowledge.DataDeleteReq) (interface{}, error) {
	s := entity.Create(new(service.KnowledgeService)).(*service.KnowledgeService)
	return s.DataDelete(req)
}

type DetailController struct {
	flow.Controller
}

func (entity *DetailController) Action(req *dto_knowledge.DetailReq) (interface{}, error) {
	s := entity.Create(new(service.KnowledgeService)).(*service.KnowledgeService)
	return s.Detail(req)
}

type SearchController struct {
	flow.Controller
}

func (entity *SearchController) Action(req *dto_knowledge.SearchReq) (interface{}, error) {
	s := entity.Create(new(service.KnowledgeService)).(*service.KnowledgeService)
	return s.Search(req)
}

type CoversController struct {
	flow.Controller
}

func (entity *CoversController) Action(req *dto.EmptyReq) (interface{}, error) {
	s := entity.Create(new(service.KnowledgeService)).(*service.KnowledgeService)
	return s.Covers(req)
}

type AuthController struct {
	flow.Controller
}

func (entity *AuthController) Action(req *dto_knowledge.AuthReq) (interface{}, error) {
	s := entity.Create(new(service.KnowledgeService)).(*service.KnowledgeService)
	return s.Auth(req)
}

type UserListController struct {
	flow.Controller
}

func (entity *UserListController) Action(req *dto_knowledge.UserListReq) (interface{}, error) {
	s := entity.Create(new(service.KnowledgeService)).(*service.KnowledgeService)
	return s.UserList(req)
}

type UserQueryController struct {
	flow.Controller
}

func (entity *UserQueryController) Action(req *dto_knowledge.UserQueryReq) (interface{}, error) {
	s := entity.Create(new(service.KnowledgeService)).(*service.KnowledgeService)
	return s.UserQuery(req)
}

type UserAddController struct {
	flow.Controller
}

func (entity *UserAddController) Action(req *dto_knowledge.UserAddReq) (interface{}, error) {
	s := entity.Create(new(service.KnowledgeService)).(*service.KnowledgeService)
	return s.UserAdd(req)
}

type UserDeleteController struct {
	flow.Controller
}

func (entity *UserDeleteController) Action(req *dto_knowledge.UserDeleteReq) (interface{}, error) {
	s := entity.Create(new(service.KnowledgeService)).(*service.KnowledgeService)
	return s.UserDelete(req)
}

type DeleteSelfController struct {
	flow.Controller
}

func (entity *DeleteSelfController) Action(req *dto_knowledge.DeleteSelfReq) (interface{}, error) {
	s := entity.Create(new(service.KnowledgeService)).(*service.KnowledgeService)
	return s.DeleteSelf(req)
}

type DataRedoController struct {
	flow.Controller
}

func (entity *DataRedoController) Action(req *dto_knowledge.DataRedoReq) (interface{}, error) {
	s := entity.Create(new(service.KnowledgeService)).(*service.KnowledgeService)
	return s.DataRedo(req)
}

type GenShareCodeController struct {
	flow.Controller
}

func (entity *GenShareCodeController) Action(req *dto_knowledge.GenShareCodeReq) (interface{}, error) {
	s := entity.Create(new(service.KnowledgeService)).(*service.KnowledgeService)
	return s.GenShareCode(req)
}

type VerifyShareCodeController struct {
	flow.Controller
}

func (entity *VerifyShareCodeController) Action(req *dto_knowledge.VerifyShareCodeReq) (interface{}, error) {
	s := entity.Create(new(service.KnowledgeService)).(*service.KnowledgeService)
	return s.VerifyShareCode(req)
}

type ShareCodeInfoController struct {
	flow.Controller
}

func (entity *ShareCodeInfoController) Action(req *dto_knowledge.InfoShareCodeReq) (interface{}, error) {
	s := entity.Create(new(service.KnowledgeService)).(*service.KnowledgeService)
	return s.ShareCodeInfo(req)
}

type SearchAdminController struct {
	flow.Controller
}

func (entity *SearchAdminController) Action(req *dto_knowledge.SearchAdminReq) (interface{}, error) {
	s := entity.Create(new(service.KnowledgeService)).(*service.KnowledgeService)
	return s.SearchAdmin(req)
}
