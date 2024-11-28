package service

import (
	"askonce/api/jobd"
	"askonce/components"
	"askonce/components/defines"
	"askonce/components/dto"
	"askonce/helpers"
	"askonce/models"
	"askonce/utils"
	"fmt"
	"github.com/duke-git/lancet/v2/cryptor"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/minio/minio-go/v7"
	"github.com/xiangtao94/golib/flow"
	"github.com/xiangtao94/golib/pkg/orm"
	"github.com/xiangtao94/golib/pkg/zlog"
	"net/url"

	"strings"
)

type FileService struct {
	flow.Service
	jobdApi *jobd.JobdApi
}

func (c *FileService) OnCreate() {
	c.jobdApi = c.Create(new(jobd.JobdApi)).(*jobd.JobdApi)
}

// 允许的格式
var allowExtension = []string{"pdf", "doc", "docx", "txt", "ppt", "pptx", "xlsx", "xls", "json"}

func (c *FileService) FileUpload(req *dto.FileUploadReq) (res *dto.FileUploadRes, err error) {
	userInfo, _ := utils.LoginInfo(c.GetCtx())
	splitF := strings.Split(req.File.Filename, ".")
	// 文件原始名称
	fileOriginName := splitF[0]
	// 文件原始格式
	fileOriginExtension := splitF[1]
	if !slice.Contain(allowExtension, fileOriginExtension) {
		return nil, components.ErrorFormatError
	}
	// 文件上传名称
	uploadObjectName := fmt.Sprintf("%v_%s", helpers.GenID(), fileOriginName)
	// 文件上传目录
	uploadObjectDir := fmt.Sprintf("%s/%s", req.Source, userInfo.UserId)
	// 文件上传路径
	uploadObjectPath := fmt.Sprintf("%s/%s.%s", uploadObjectDir, uploadObjectName, fileOriginExtension)
	// 文件bucket
	bucketName := defines.BucketOrigin
	minioClient, err := helpers.GetMinioClient(c.GetCtx())
	if err != nil {
		return nil, components.ErrorFileClientError
	}
	f, err := req.File.Open()
	if err != nil {
		c.LogErrorf("file open fail, %v", err)
		return
	}
	defer f.Close()
	objectInfo, err := minioClient.PutObject(c.GetCtx(), bucketName, uploadObjectPath, f, req.File.Size, minio.PutObjectOptions{})
	if err != nil {
		zlog.Errorf(c.GetCtx(), "upload file fail: %+v", err)
		return nil, components.ErrorFileUploadError
	}
	filePath2 := url.PathEscape(uploadObjectPath)
	fileUrl := minioClient.EndpointURL().String() + "/" + bucketName + "/" + filePath2
	fileDao := flow.Create(c.GetCtx(), new(models.FileDao))
	add := &models.File{
		Id:        cryptor.HmacMd5(objectInfo.Key, "askonce"),
		Name:      uploadObjectName,
		Extension: fileOriginExtension,
		Path:      fileUrl,
		Source:    req.Source,
		UserId:    userInfo.UserId,
		CrudModel: orm.CrudModel{
			CreatedAt: objectInfo.LastModified,
			UpdatedAt: objectInfo.LastModified,
		},
	}
	err = fileDao.Insert(add)
	if err != nil {
		return nil, err
	}
	res = &dto.FileUploadRes{
		FileId: add.Id,
		Name:   add.Name,
		Size:   objectInfo.Size,
	}
	return
}
