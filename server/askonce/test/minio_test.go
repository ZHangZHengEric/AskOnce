// Package helpers -----------------------------
// @file      : minio_test.go
// @author    : xiangtao
// @contact   : xiangtao1994@gmail.com
// @time      : 2024/12/16 11:18
// -------------------------------------------
package test

import (
	"askonce/components/defines"
	"askonce/helpers"
	"fmt"
	"github.com/minio/minio-go/v7"
	"net/url"
	"strings"
	"testing"
)

func TestMinio(t *testing.T) {
	Init()
	minioClient, err := helpers.GetMinioClient(Ctx)
	if err != nil {
		return
	}
	bucketName := defines.BucketOrigin
	objectsCh := make(chan minio.ObjectInfo)
	go func() {
		defer close(objectsCh)
		parsedURL, _ := url.Parse("http://36.133.44.114:20002/origin/knowledge%2F681298043804%2F841996222692_123.txt")
		p1, _ := url.PathUnescape(parsedURL.Path)
		p, _ := strings.CutPrefix(p1, "/"+bucketName+"/")
		objectsCh <- minio.ObjectInfo{
			Key: p,
		}
	}()
	errorCh := minioClient.RemoveObjectsWithResult(Ctx, bucketName, objectsCh, minio.RemoveObjectsOptions{})
	for e := range errorCh {
		if e.Err != nil {
			fmt.Printf("minio remove objects fail, %v", e.Err)
		}
	}
}
