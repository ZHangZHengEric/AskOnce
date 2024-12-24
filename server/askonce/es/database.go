// Package es -----------------------------
// @file      : database.go
// @author    : xiangtao
// @contact   : xiangtao1994@gmail.com
// @time      : 2024/12/23 12:26
// -------------------------------------------
package es

import (
	"askonce/helpers"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/xiangtao94/golib/pkg/env"
	"github.com/xiangtao94/golib/pkg/zlog"
	"os"
	"path/filepath"
)

type DatabaseDocument struct {
	DocDocument
	DatabaseName  string `json:"database_name,omitempty"`
	TableName     string `json:"table_name,omitempty"`
	TableComment  string `json:"table_comment,omitempty"`
	ColumnName    string `json:"column_name,omitempty"`
	ColumnComment string `json:"column_comment,omitempty"`
	ColumnType    string `json:"column_type,omitempty"`
	Sql           string `json:"sql,omitempty"`
}

func (cd *DatabaseDocument) ConvertAny() any {
	return cd
}

const (
	IndexSuffixTable       = "table"
	IndexSuffixColumn      = "column"
	IndexSuffixColumnValue = "column_value"
	IndexSuffixFaq         = "faq"
)

var indexSuffix = []string{IndexSuffixTable, IndexSuffixColumn, IndexSuffixColumnValue, IndexSuffixFaq}

func DatabaseIndexCreate(ctx *gin.Context, indexName string) (err error) {
	envPath := filepath.Join(env.GetConfDirPath(), "mount/es")
	settingsPath := filepath.Join(envPath, "database_setting.json")
	setting := &types.IndexSettings{}
	settingsFile, err := os.ReadFile(settingsPath)
	if err != nil {
		return fmt.Errorf("read common mapping file failed: %v", err)
	}
	err = setting.UnmarshalJSON(settingsFile)
	if err != nil {
		return
	}

	mappingPath := filepath.Join(envPath, "database_mapping.json")
	mappingFile, err := os.ReadFile(mappingPath)
	if err != nil {
		return fmt.Errorf("read common mapping file failed: %v", err)
	}
	// 需要创建四个mapping
	mappingAll := map[string]jsoniter.RawMessage{}
	err = jsoniter.Unmarshal(mappingFile, &mappingAll)
	if err != nil {
		return fmt.Errorf("parse databse mapping file failed: %v", err)
	}
	if len(mappingAll) != 4 {
		return fmt.Errorf("parse databse mapping  failed")
	}
	for _, v := range indexSuffix {
		indexNameNew := fmt.Sprintf("%s_%s", indexName, v)
		tmp := mappingAll[v]
		if tmp == nil {
			return fmt.Errorf("parse %s database mapping  failed", v)
		}
		mapping := &types.TypeMapping{}
		err = mapping.UnmarshalJSON(tmp)
		if err != nil {
			return err
		}
		err = helpers.EsClient.CreateIndex(ctx, indexNameNew, mapping, setting)
		if err != nil {
			return err
		}
	}
	return nil
}

func DatabaseIndexDelete(ctx *gin.Context, indexName string) (err error) {
	for _, v := range indexSuffix {
		indexNameNew := fmt.Sprintf("%s_%s", indexName, v)
		err = helpers.EsClient.DeleteIndex(ctx, indexNameNew)
		if err != nil {
			zlog.Errorf(ctx, "delete index %s failed: %v", indexName, err)
		}
	}
	return nil
}
func DatabaseDocumentDelete(ctx *gin.Context, indexName string, docIds []int64) (err error) {
	for _, v := range indexSuffix {
		indexNameNew := fmt.Sprintf("%s_%s", indexName, v)
		err = CommonDocumentDelete(ctx, indexNameNew, docIds)
		if err != nil {
			zlog.Errorf(ctx, "delete index %s failed: %v", indexName, err)
		}
	}
	return nil
}

func DatabaseDocumentInsert(ctx *gin.Context, indexName string,
	tables []*DatabaseDocument,
	columns []*DatabaseDocument,
	columnValues []*DatabaseDocument,
	faqs []*DatabaseDocument,
) (err error) {
	if err := CommonBatchInsert(ctx, fmt.Sprintf("%s_%s", indexName, IndexSuffixTable), tables); err != nil {
		return err
	}
	if err := CommonBatchInsert(ctx, fmt.Sprintf("%s_%s", indexName, IndexSuffixColumn), columns); err != nil {
		return err
	}
	if err := CommonBatchInsert(ctx, fmt.Sprintf("%s_%s", indexName, IndexSuffixColumnValue), columnValues); err != nil {
		return err
	}
	if err := CommonBatchInsert(ctx, fmt.Sprintf("%s_%s", indexName, IndexSuffixFaq), faqs); err != nil {
		return err
	}
	return
}
