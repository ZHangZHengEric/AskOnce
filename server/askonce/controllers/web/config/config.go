package config

import (
	"askonce/components/dto"
	"askonce/components/dto/dto_config"
	"github.com/xiangtao94/golib/flow"
)

type DetailController struct {
	flow.Controller
}

func (entity *DetailController) Action(req *dto_config.DetailReq) (interface{}, error) {
	s := entity.Create(new(service.ConfigService)).(*service.ConfigService)
	return s.Detail(req)
}

type SaveController struct {
	flow.Controller
}

func (entity *SaveController) Action(req *dto_config.SaveReq) (interface{}, error) {
	s := entity.Create(new(service.ConfigService)).(*service.ConfigService)
	return s.Save(req)
}

type DictController struct {
	flow.Controller
}

func (entity *DictController) Action(req *dto.EmptyReq) (interface{}, error) {
	s := entity.Create(new(service.ConfigService)).(*service.ConfigService)
	return s.Dict(req)
}
