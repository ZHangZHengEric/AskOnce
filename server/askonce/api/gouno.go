package api

import (
	"askonce/components/defines"
	"askonce/conf"
	"askonce/helpers"
	"bytes"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/pkg/mimedb"
	"github.com/xiangtao94/golib/flow"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"
)

type GoUnoApi struct {
	flow.Api
}

func (a *GoUnoApi) OnCreate() {
	a.Client = conf.WebConf.Api["gouno"]
}

func (a *GoUnoApi) HtmlToDocx(fileName string, content string) (fileUrl string, err error) {
	filePrefix := strings.Split(fileName, ".")[0]
	if filePrefix == "" {
		filePrefix = "askonce_doc_" + time.Now().String()
		fileName = fmt.Sprintf("%s.html", filePrefix)
	} else {
		filePrefix = filePrefix + time.Now().String()
		fileName = fmt.Sprintf("%s.html", filePrefix)
	}
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return
	}
	src := strings.NewReader(content)
	// 拷贝文件数据到part
	io.Copy(part, src)
	writer.Close()
	req, err := http.NewRequest("POST", a.Client.Domain+"/unoconv/docx", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	a.Client.HTTPClient = new(http.Client)
	resp, err := a.Client.HTTPClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	objectFullPath := filePrefix + ".docx"
	minioClient, err := helpers.GetMinioClient(a.GetCtx())
	_, err = minioClient.PutObject(a.GetCtx(), defines.BucketTmp, objectFullPath, resp.Body, resp.ContentLength, minio.PutObjectOptions{
		ContentType: mimedb.TypeByExtension("docx"),
	})
	if err != nil {
		return
	}
	filePath2 := url.PathEscape(objectFullPath)
	fileUrl = minioClient.EndpointURL().String() + "/" + defines.BucketTmp + "/" + filePath2
	return fileUrl, nil
}

func (a *GoUnoApi) FileToDocx(fileReq *multipart.FileHeader) (fileUrl string, err error) {
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
	req, err := http.NewRequest("POST", a.Client.Domain+"/unoconv/docx", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	a.Client.HTTPClient = new(http.Client)
	resp, err := a.Client.HTTPClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	minioClient, err := helpers.GetMinioClient(a.GetCtx())
	filePrefix := strings.Split(fileReq.Filename, ".")[0]
	_, err = minioClient.PutObject(a.GetCtx(), defines.BucketTmp, filePrefix+".docx", resp.Body, resp.ContentLength, minio.PutObjectOptions{
		ContentType: mimedb.TypeByExtension(filepath.Ext(fileReq.Filename)),
	})
	if err != nil {
		return
	}
	objectFullPath := filePrefix + ".docx"
	filePath2 := url.PathEscape(objectFullPath)
	fileUrl = minioClient.EndpointURL().String() + "/" + defines.BucketTmp + "/" + filePath2
	return fileUrl, nil
}
