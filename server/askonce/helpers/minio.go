package helpers

import (
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/xiangtao94/golib/pkg/errors"
	"github.com/xiangtao94/golib/pkg/zlog"

	"askonce/conf"
	"net/url"
)

func GetMinioClient(ctx *gin.Context) (*minio.Client, error) {
	minioConf := conf.WebConf.MinioConf
	m, err := newMClientByAK(minioConf.Addr, minioConf.AK, minioConf.SK)
	if err != nil {
		zlog.Errorf(ctx, "minio init error: %s", err.Error())
		return nil, err
	}
	return m, nil
}

func newMClientByAK(endpoint string, accessKey, secretKey string) (*minio.Client, error) {
	endpointUrl, err := url.Parse(endpoint)
	if err != nil {
		return nil, errors.ErrorSystemError
	}
	minioClient, err := minio.New(endpointUrl.Host, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}
	return minioClient, nil
}
