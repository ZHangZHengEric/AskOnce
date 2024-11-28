package file

import (
	"askonce/conf"
	"bytes"
	"github.com/minio/minio-go/v7"
	"github.com/minio/pkg/mimedb"
	"github.com/xiangtao94/golib/flow"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
)

type GoUnoApi struct {
	flow.Api
}

func (a *GoUnoApi) OnCreate() {
	a.Client = conf.WebConf.Api["gouno"]
}

func (a *GoUnoApi) FileToPdf(fileReq *multipart.FileHeader) (fileUrl string, err error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", fileReq.Filename)
	if err != nil {
		return
	}
	src, err := fileReq.Open()
	defer src.Close()
	// 拷贝文件数据到part
	io.Copy(part, src)
	writer.Close()
	req, err := http.NewRequest("POST", a.Client.Domain+"/unoconv/pdf", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	a.Client.HTTPClient = new(http.Client)
	resp, err := a.Client.HTTPClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	fileApi := a.Create(new(MinioApi)).(*MinioApi)
	filePrefix := strings.Split(fileReq.Filename, ".")[0]
	_, err = fileApi.minioClient.PutObject(a.GetCtx(), "tmp", filePrefix+".pdf", resp.Body, resp.ContentLength, minio.PutObjectOptions{
		ContentType: mimedb.TypeByExtension(filepath.Ext(fileReq.Filename)),
	})
	if err != nil {
		return
	}
	objectFullPath := filePrefix + ".pdf"
	filePath2 := url.PathEscape(objectFullPath)
	fileUrl = fileApi.minioClient.EndpointURL().String() + "/" + "tmp" + "/" + filePath2
	return fileUrl, nil
}
