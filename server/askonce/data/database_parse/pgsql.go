// Package database_parse -----------------------------
// @file      : pgsql.go
// @author    : xiangtao
// @contact   : xiangtao1994@gmail.com
// @time      : 2024/12/22 14:57
// -------------------------------------------
package database_parse

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PostgreSQLHandler 实现 DatabaseHandler 接口
type PostgreSQLHandler struct {
	DB  *gorm.DB
	ctx *gin.Context
}

func (h *PostgreSQLHandler) Ping() error {
	// Ping数据库
	sqlDB, err := h.DB.DB()
	if err != nil {
		return err
	}
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("数据库Ping失败: %v", err)
	}
	return nil
}

func (h *PostgreSQLHandler) SetCtx(ctx *gin.Context) error {
	h.ctx = ctx
	return nil
}

func (h *PostgreSQLHandler) GetCtx() (ctx *gin.Context) {
	return h.ctx
}

func (h *PostgreSQLHandler) GetTables() ([]TableInfo, error) {
	query := " SELECT t.table_name, COALESCE(d.description, '') AS table_comment  FROM information_schema.tables t LEFT JOIN pg_class c ON c.relname = t.table_name LEFT JOIN pg_namespace n ON n.oid = c.relnamespace LEFT JOIN pg_description d ON c.oid = d.objoid  WHERE t.table_schema = 'public' AND n.nspname = 'public';`"
	var tables []TableInfo
	err := h.DB.Select(&tables, query).Error
	return tables, err
}

func (h *PostgreSQLHandler) GetColumns(table string) ([]ColumnInfo, error) {
	query := fmt.Sprintf(`
		SELECT column_name, data_type, col_description((SELECT oid FROM pg_class WHERE relname='%s'), ordinal_position::int) AS column_comment 
		FROM information_schema.columns 
		WHERE table_name='%s';`, table, table)
	var columns []ColumnInfo
	err := h.DB.Model(ColumnInfo{}).Raw(query).Scan(&columns).Error
	if err != nil {
		return nil, err
	}
	return columns, nil
}

func (h *PostgreSQLHandler) GetSampleData(table, column string) ([]string, error) {
	query := fmt.Sprintf("SELECT %s FROM %s LIMIT 1000;", column, table)
	var samples []string
	err := h.DB.Select(&samples, query).Error
	return samples, err
}

func (h *PostgreSQLHandler) Close() error {
	return nil
}
