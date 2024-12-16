package data

import (
	"archive/zip"
	"askonce/api/jobd"
	"askonce/components"
	"askonce/components/defines"
	"askonce/helpers"
	"askonce/models"
	"askonce/utils"
	"bytes"
	"fmt"
	"github.com/duke-git/lancet/v2/cryptor"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/minio/minio-go/v7"
	"github.com/xiangtao94/golib/flow"
	"github.com/xiangtao94/golib/pkg/errors"
	"github.com/xiangtao94/golib/pkg/orm"
	"github.com/xiangtao94/golib/pkg/zlog"
	"golang.org/x/sync/errgroup"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"mime/multipart"
	"net/url"
	"os"
	"path"
	"strings"
	"sync"
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
var allowExtension = []string{".pdf", ".doc", ".docx", ".txt", ".ppt", ".pptx", ".xlsx", ".xls", ".json"}

func (f *FileData) UploadByText(userId string, fileName string, content string, source string) (add *models.File, err error) {
	uploadObjectPath, uploadObjectName, fileOriginName, fileOriginExt, err := parseFileName(fileName, source, userId)
	if err != nil {
		return nil, err
	}
	// 文件bucket
	bucketName := defines.BucketOrigin
	minioClient, err := helpers.GetMinioClient(f.GetCtx())
	if err != nil {
		return nil, components.ErrorFileClientError
	}
	zlog.Infof(f.GetCtx(), "upload object %s to minio bucket %s", uploadObjectPath, bucketName)
	objectInfo, err := minioClient.PutObject(f.GetCtx(), bucketName, uploadObjectPath, strings.NewReader(content), int64(len(content)), minio.PutObjectOptions{})
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
		Extension:  fileOriginExt,
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

func (f *FileData) UploadByFile(userId string, file *multipart.FileHeader, source string) (add *models.File, err error) {
	fileR, err := file.Open()
	if err != nil {
		f.LogErrorf("file open fail, %v", err)
		return
	}
	defer fileR.Close()
	uploadObjectPath, uploadObjectName, fileOriginName, fileOriginExt, err := parseFileName(file.Filename, source, userId)
	if err != nil {
		return nil, err
	}
	// 文件bucket
	bucketName := defines.BucketOrigin
	minioClient, err := helpers.GetMinioClient(f.GetCtx())
	if err != nil {
		return nil, components.ErrorFileClientError
	}
	zlog.Infof(f.GetCtx(), "upload object %s to minio bucket %s", uploadObjectPath, bucketName)
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
		Extension:  fileOriginExt,
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

func (f *FileData) UploadByZip(userId string, zipUrl string, source string) (files []*models.File, err error) {
	network, err := utils.IsNetWorkUrlOrLocal(zipUrl)
	if err != nil {
		return nil, err
	}
	var zipData []byte
	if network {
		zipData, err = utils.DownloadZip(zipUrl)
		if err != nil {
			return nil, err
		}
	} else {
		zipData, err = os.ReadFile(zipUrl)
	}
	reader, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
	if err != nil {
		return nil, errors.NewError(errors.DEFAULT_ERROR, fmt.Sprintf("failed to open zip reader: %s", err))
	}
	files = make([]*models.File, 0)
	wg, _ := errgroup.WithContext(f.GetCtx())
	lock := sync.RWMutex{}
	fileChunck := slice.Chunk(reader.File, 100)
	minioClient, err := helpers.GetMinioClient(f.GetCtx())
	if err != nil {
		return nil, err
	}
	for _, zipFiles := range fileChunck {
		for _, zipFile := range zipFiles {
			if zipFile.FileInfo().IsDir() {
				continue
			}
			wg.Go(func() error {
				file, err := f.uploadByZipDo(minioClient, userId, zipFile, source)
				if err != nil {
					return err
				}

				lock.Lock()
				files = append(files, file)
				lock.Unlock()
				return nil
			})
		}
		if err = wg.Wait(); err != nil {
			return nil, err
		}
	}
	err = f.fileDao.BatchInsert(files)
	return
}

func (f *FileData) uploadByZipDo(minioClient *minio.Client, userId string, file *zip.File, source string) (add *models.File, err error) {
	fileR, err := file.Open()
	if err != nil {
		f.LogErrorf("file open fail, %v", err)
		return
	}
	defer fileR.Close()
	var fileName string
	if file.Flags == 0 {
		i := bytes.NewReader([]byte(file.Name))
		decoder := transform.NewReader(i, simplifiedchinese.GB18030.NewDecoder())
		content, _ := ioutil.ReadAll(decoder)
		fileName = string(content)
	} else {
		fileName = file.Name
	}
	uploadObjectPath, uploadObjectName, fileOriginName, fileOriginExt, err := parseFileName(fileName, source, userId)
	if err != nil {
		return nil, err
	} // 文件bucket
	bucketName := defines.BucketOrigin
	if err != nil {
		return nil, components.ErrorFileClientError
	}
	zlog.Infof(f.GetCtx(), "upload object %s to minio bucket %s", uploadObjectPath, bucketName)
	objectInfo, err := minioClient.PutObject(f.GetCtx(), bucketName, uploadObjectPath, fileR, file.FileInfo().Size(), minio.PutObjectOptions{})
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
		Extension:  fileOriginExt,
		Path:       fileUrl,
		Size:       objectInfo.Size,
		Source:     source,
		UserId:     userId,
		CrudModel: orm.CrudModel{
			CreatedAt: objectInfo.LastModified,
			UpdatedAt: objectInfo.LastModified,
		},
	}
	return
}

func (f *FileData) DeleteByFileIds(fileIds []string) (err error) {
	if len(fileIds) == 0 {
		return
	}
	files, err := f.fileDao.GetByIds(fileIds)
	if err != nil {
		return err
	}
	minioClient, err := helpers.GetMinioClient(f.GetCtx())
	if err != nil {
		return err
	}
	bucketName := defines.BucketOrigin
	objectsCh := make(chan minio.ObjectInfo)
	go func() {
		defer close(objectsCh)
		for _, file := range files {
			parsedURL, _ := url.Parse(file.Path)
			p1, _ := url.PathUnescape(parsedURL.Path)
			p, _ := strings.CutPrefix(p1, "/"+bucketName+"/")
			objectsCh <- minio.ObjectInfo{
				Key: p,
			}
		}
	}()
	errorCh := minioClient.RemoveObjects(f.GetCtx(), bucketName, objectsCh, minio.RemoveObjectsOptions{})
	for e := range errorCh {
		if e.Err != nil {
			return fmt.Errorf("minio remove objects fail, %v", e.Err)
		}
	}
	err = f.fileDao.DeleteByFileIds(fileIds)
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

func parseFileName(fileName, source, userId string) (uploadObjectPath, uploadObjectName, fileOriginName, fileOriginExt string, err error) {
	pathN := path.Base(fileName)
	// 文件原始格式
	fileOriginExt = path.Ext(pathN)
	// 文件原始名称
	fileOriginName = pathN[0 : len(pathN)-len(fileOriginExt)]
	if !slice.Contain(allowExtension, fileOriginExt) {
		err = components.ErrorFormatError
		return
	}
	// 文件上传名称
	uploadObjectName = fmt.Sprintf("%v_%s", helpers.GenID(), fileOriginName)
	// 文件上传目录
	uploadObjectDir := fmt.Sprintf("%s/%s", source, userId)
	// 文件上传路径
	uploadObjectPath = fmt.Sprintf("%s/%s%s", uploadObjectDir, uploadObjectName, fileOriginExt)
	return
}
