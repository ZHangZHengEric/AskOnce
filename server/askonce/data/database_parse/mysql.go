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

	//// 统计去重值和未去重值的数量
	//var stats ColumnStats
	//h.DB.Model(&ColumnStats{}).Raw(fmt.Sprintf("SELECT COUNT(DISTINCT %s) as distinct_count, COUNT(%s) as total_count FROM %s", columnName, columnName, tableName)).Scan(&stats)

	zlog.Infof(h.ctx, "表【%s】列【%s】符合条件，开始获取value值", tableName, columnName)
	// 计算 Top N 的值
	n := 1000 // 可调整
	type topValue struct {
		Value string `gorm:"column:value"`
		Cnt   int64  `gorm:"column:cnt"`
	}
	var topValues []topValue
	sqlI := fmt.Sprintf(`
		SELECT %s as value, COUNT(*) AS cnt 
		FROM %s 
		GROUP BY %s 
		ORDER BY cnt DESC 
		LIMIT ?
	`, wrapColumnName(columnName), tableName, wrapColumnName(columnName))
	err := h.DB.Model(topValue{}).Raw(sqlI, n).Scan(&topValues).Error
	if err != nil {
		return nil, err
	}
	topValuestr := []string{}
	for _, t := range topValues {
		if t.Value != "" {
			topValuestr = append(topValuestr, t.Value)
		}
	}
	zlog.Infof(h.ctx, "表【%s】列【%s】符合条件，获取value值完成，数量【%v】", tableName, columnName, len(topValuestr))
	return topValuestr, nil
}
func wrapColumnName(columnName string) string {
	return fmt.Sprintf("`%s`", columnName)
}
func (h *MySQLHandler) Close() error {
	return nil
}
