// Package es -----------------------------
// @file      : common_doc.go
// @author    : xiangtao
// @contact   : xiangtao1994@gmail.com
// @time      : 2024/12/13 19:05
// -------------------------------------------
package es

import (
	"askonce/helpers"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/gin-gonic/gin"
	"github.com/xiangtao94/golib/pkg/env"
	"os"
	"path/filepath"
)

type DocDocument struct {
	CommonDocument       // 索引列 内容是 content
	DocSegmentId   int64 `json:"doc_segment_id"`
	Start          int   `json:"start"`
	End            int   `json:"end"`
}

func (cd *DocDocument) ConvertAny() any {
	return cd
}

func DocIndexCreate(ctx *gin.Context, indexName string) (err error) {
	envPath := filepath.Join(env.GetConfDirPath(), "mount/es")
	mappingPath := filepath.Join(envPath, "doc_mapping.json")
	mappingFile, err := os.ReadFile(mappingPath)
	if err != nil {
		return fmt.Errorf("read doc mapping file failed: %v", err)
	}
	mapping := &types.TypeMapping{}
	err = mapping.UnmarshalJSON(mappingFile)
	if err != nil {
		return
	}
	settingsPath := filepath.Join(envPath, "doc_setting.json")
	setting := &types.IndexSettings{}
	settingsFile, err := os.ReadFile(settingsPath)
	if err != nil {
		return fmt.Errorf("read doc mapping file failed: %v", err)
	}
	err = setting.UnmarshalJSON(settingsFile)
	if err != nil {
		return
	}
	err = helpers.EsClient.CreateIndex(ctx, indexName, mapping, setting)
	if err != nil {
		return
	}
	return nil
}

func DocIndexDelete(ctx *gin.Context, indexName string) (err error) {
	err = helpers.EsClient.DeleteIndex(ctx, indexName)
	if err != nil {
		return
	}
	return nil
}
