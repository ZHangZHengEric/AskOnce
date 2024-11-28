package files

import (
	"askonce/components/dto"
	"askonce/service"
	"github.com/xiangtao94/golib/flow"
)

type FileUploadController struct {
	flow.Controller
}

func (entity *FileUploadController) Action(req *dto.FileUploadReq) (interface{}, error) {
	s := flow.Create(entity.GetCtx(), new(service.FileService))
	return s.FileUpload(req)
}
