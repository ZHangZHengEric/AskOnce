package data

import (
	"askonce/api/jobd"
	"askonce/components"
	"askonce/models"
	"github.com/xiangtao94/golib/flow"
)

type FileData struct {
	flow.Data
	fileDao *models.FileDao
	jobdApi *jobd.JobdApi
}

func (f *FileData) OnCreate() {
	f.fileDao = flow.Create(f.GetCtx(), new(models.FileDao))
	f.jobdApi = flow.Create(f.GetCtx(), new(jobd.JobdApi))
}

// 文件转文本
func (d *FileData) ConvertFileToText(fileId string) (fileName string, output string, err error) {
	// 获取文件
	file, err := d.fileDao.GetById(fileId)
	if err != nil {
		return
	}
	if file == nil { // 文件不存在
		return "", "", components.ErrorFileNoExist
	}
	fileToText, err := d.jobdApi.FileToText(file.Path)
	if err != nil {
		return
	}
	return file.Name, fileToText.Text, nil
}
