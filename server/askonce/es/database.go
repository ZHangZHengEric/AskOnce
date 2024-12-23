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
	"os"
	"path/filepath"
)

type TableDocument struct {
	CommonDocument        // 索引列 内容是 TableName + TableComment +  TableSchema
	DatabaseName   string `json:"database_name"`
	TableName      string `json:"table_name"`
	TableComment   string `json:"table_comment"`
}

func (cd *TableDocument) ConvertAny() any {
	return cd
}

type TableColumnDocument struct {
	CommonDocument        // 索引列 内容是 ColumnName + ColumnComment +  ColumnType
	DatabaseName   string `json:"database_name"`
	TableName      string `json:"table_name"`
	ColumnName     string `json:"column_name"`
	ColumnComment  string `json:"column_comment"`
	ColumnType     string `json:"column_type"`
}

func (cd *TableColumnDocument) ConvertAny() any {
	return cd
}

type TableColumnValueDocument struct {
	CommonDocument        // 索引列 内容是 值
	DatabaseName   string `json:"database_name"`
	TableName      string `json:"table_name"`
	ColumnName     string `json:"column_name"`
}

func (cd *TableColumnValueDocument) ConvertAny() any {
	return cd
}

type DatabaseFaqDocument struct {
	CommonDocument //  索引列 内容是 faq

	DatabaseName string `json:"database_name"`
	SqlKey       string `json:"sql_key"`
}

func (cd *DatabaseFaqDocument) ConvertAny() any {
	return cd
}

var indexSuffix = []string{"table", "column", "column_value", "faq"}

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
			return err
		}
	}
	return nil
}

func DatabaseDocumentInsert(ctx *gin.Context, indexName string,
	tables []*TableDocument,
	columns []*TableColumnDocument,
	columnValues []*TableColumnValueDocument,
	faqs []*DatabaseFaqDocument,
) (err error) {
	if err := CommonBatchInsert(ctx, fmt.Sprintf("%s_%s", indexName, "table"), tables); err != nil {
		return err
	}
	if err := CommonBatchInsert(ctx, fmt.Sprintf("%s_%s", indexName, "column"), columns); err != nil {
		return err
	}
	if err := CommonBatchInsert(ctx, fmt.Sprintf("%s_%s", indexName, "column_value"), columnValues); err != nil {
		return err
	}
	if err := CommonBatchInsert(ctx, fmt.Sprintf("%s_%s", indexName, "faq"), faqs); err != nil {
		return err
	}
	return
}
