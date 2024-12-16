package service

import (
	"askonce/api/jobd"
	"askonce/components/dto"
	"askonce/data"
	"askonce/utils"
	"github.com/xiangtao94/golib/flow"
)

type FileService struct {
	flow.Service
	jobdApi *jobd.JobdApi
}

func (c *FileService) OnCreate() {
	c.jobdApi = c.Create(new(jobd.JobdApi)).(*jobd.JobdApi)
}

func (c *FileService) FileUpload(req *dto.FileUploadReq) (res *dto.FileUploadRes, err error) {
	userInfo, _ := utils.LoginInfo(c.GetCtx())
	file, err := flow.Create(c.GetCtx(), new(data.FileData)).UploadByFile(userInfo.UserId, req.File, req.Source)
	if err != nil {
		return nil, err
	}
	res = &dto.FileUploadRes{
		FileId: file.Id,
		Name:   file.Name,
		Size:   file.Size,
	}
	return
}
