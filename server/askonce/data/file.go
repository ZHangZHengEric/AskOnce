package data

import (
	"askonce/api/jobd"
	"askonce/components"
	"askonce/components/defines"
	"askonce/helpers"
	"askonce/models"
	"fmt"
	"github.com/duke-git/lancet/v2/cryptor"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/minio/minio-go/v7"
	"github.com/xiangtao94/golib/flow"
	"github.com/xiangtao94/golib/pkg/orm"
	"github.com/xiangtao94/golib/pkg/zlog"
	"mime/multipart"
	"net/url"
	"strings"
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
func (f *FileData) ConvertFileToText(fileId string) (fileName string, output string, err error) {
	// 获取文件
	file, err := f.fileDao.GetById(fileId)
	if err != nil {
		return
	}
	if file == nil { // 文件不存在
		return "", "", components.ErrorFileNoExist
	}
	fileToText, err := f.jobdApi.FileToText(file.Path)
	if err != nil {
		return
	}
	return file.Name, fileToText.Text, nil
}

// 允许的格式
var allowExtension = []string{"pdf", "doc", "docx", "txt", "ppt", "pptx", "xlsx", "xls", "json"}

func (f *FileData) Upload(userId string, file *multipart.FileHeader, source string) (add *models.File, err error) {
	splitF := strings.Split(file.Filename, ".")
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
	uploadObjectDir := fmt.Sprintf("%s/%s", source, userId)
	// 文件上传路径
	uploadObjectPath := fmt.Sprintf("%s/%s.%s", uploadObjectDir, uploadObjectName, fileOriginExtension)
	// 文件bucket
	bucketName := defines.BucketOrigin
	minioClient, err := helpers.GetMinioClient(f.GetCtx())
	if err != nil {
		return nil, components.ErrorFileClientError
	}
	fileR, err := file.Open()
	if err != nil {
		f.LogErrorf("file open fail, %v", err)
		return
	}
	defer fileR.Close()
	objectInfo, err := minioClient.PutObject(f.GetCtx(), bucketName, uploadObjectPath, fileR, file.Size, minio.PutObjectOptions{})
	if err != nil {
		zlog.Errorf(f.GetCtx(), "upload file fail: %+v", err)
		return nil, components.ErrorFileUploadError
	}
	filePath2 := url.PathEscape(uploadObjectPath)
	fileUrl := minioClient.EndpointURL().String() + "/" + bucketName + "/" + filePath2
	add = &models.File{
		Id:         cryptor.HmacMd5(objectInfo.Key, "askonce"),
		Name:       uploadObjectName,
		OriginName: fileOriginName,
		Extension:  fileOriginExtension,
		Path:       fileUrl,
		Size:       objectInfo.Size,
		Source:     source,
		UserId:     userId,
		CrudModel: orm.CrudModel{
			CreatedAt: objectInfo.LastModified,
			UpdatedAt: objectInfo.LastModified,
		},
	}
	err = f.fileDao.Insert(add)
	if err != nil {
		return nil, err
	}
	return
}

func (f *FileData) GetFileByFileIds(fileIds []string) (res map[string]*models.File, err error) {
	res = make(map[string]*models.File)
	files, err := f.fileDao.GetByIds(fileIds)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		res[file.Id] = file
	}
	return
}
