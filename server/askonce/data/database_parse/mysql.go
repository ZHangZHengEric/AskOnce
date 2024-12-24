// Package database -----------------------------
// @file      : mysql.go
// @author    : xiangtao
// @contact   : xiangtao1994@gmail.com
// @time      : 2024/12/22 14:56
// -------------------------------------------
package database_parse

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xiangtao94/golib/pkg/zlog"
	"gorm.io/gorm"
)

// MySQLHandler 实现 DatabaseHandler 接口
type MySQLHandler struct {
	DB  *gorm.DB
	ctx *gin.Context
}

func (h *MySQLHandler) Ping() error {
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

func (h *MySQLHandler) SetCtx(ctx *gin.Context) error {
	h.ctx = ctx
	return nil
}

func (h *MySQLHandler) GetCtx() (ctx *gin.Context) {
	return h.ctx
}

func (h *MySQLHandler) GetTables() ([]TableInfo, error) {
	var tables []TableInfo
	err := h.DB.Model(TableInfo{}).Raw("SELECT TABLE_NAME, TABLE_COMMENT FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA = DATABASE();").Scan(&tables).Error
	return tables, err
}

func (h *MySQLHandler) GetColumns(table string) ([]ColumnInfo, error) {
	query := fmt.Sprintf(`
		SELECT COLUMN_NAME,        SUBSTRING_INDEX(COLUMN_TYPE, '(', 1) AS COLUMN_TYPE, 
		       COLUMN_COMMENT 
		FROM INFORMATION_SCHEMA.COLUMNS 
		WHERE TABLE_NAME = '%s' AND TABLE_SCHEMA = DATABASE();`, table)
	var columns []ColumnInfo
	err := h.DB.Model(ColumnInfo{}).Raw(query).Scan(&columns).Error
	if err != nil {
		return nil, err
	}
	return columns, nil
}

type ColumnStats struct {
	DistinctCount int64 `gorm:"column:distinct_count"`
	TotalCount    int64 `gorm:"column:total_count"`
}

func (h *MySQLHandler) GetSampleData(tableName, columnName string) ([]string, error) {

	// 统计去重值和未去重值的数量
	var stats ColumnStats
	h.DB.Model(&ColumnStats{}).Raw(fmt.Sprintf("SELECT COUNT(DISTINCT %s) as distinct_count, COUNT(%s) as total_count FROM %s", columnName, columnName, tableName)).Scan(&stats)

	// 判断条件
	if stats.DistinctCount*2 >= stats.TotalCount || stats.DistinctCount >= 10000 {
		return nil, nil
	}
	zlog.Infof(h.ctx, "表【%s】列【%s】符合条件，开始获取value值", tableName, columnName)
	// 计算 Top N 的值
	var topValues []string
	n := 1000 // 可调整
	h.DB.Raw(fmt.Sprintf(`
		SELECT %s, COUNT(*) AS cnt 
		FROM %s 
		GROUP BY %s 
		ORDER BY cnt DESC 
		LIMIT ?
	`, columnName, tableName, columnName), n).Scan(&topValues)
	return topValues, nil
}

func (h *MySQLHandler) Close() error {
	return nil
}
