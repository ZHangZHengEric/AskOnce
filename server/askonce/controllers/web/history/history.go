package history

import (
	"askonce/components/dto/dto_history"
	"github.com/xiangtao94/golib/flow"
)

type AskController struct {
	flow.Controller
}

func (entity *AskController) Action(req *dto_history.AskReq) (interface{}, error) {
	s := entity.Create(new(service.HistoryService)).(*service.HistoryService)
	return s.Ask(req)
}
