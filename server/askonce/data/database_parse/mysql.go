// Package database -----------------------------
// @file      : mysql.go
// @author    : xiangtao
// @contact   : xiangtao1994@gmail.com
// @time      : 2024/12/22 14:56
// -------------------------------------------
package database_parse

import (
	"fmt"
	"gorm.io/gorm"
)

// MySQLHandler 实现 DatabaseHandler 接口
type MySQLHandler struct {
	DB *gorm.DB
}

func (h *MySQLHandler) Ping() error {
	// Ping数据库
	sqlDB, err := h.DB.DB()
	if err != nil {
		return err
	}
	defer sqlDB.Close()
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("数据库Ping失败: %v", err)
	}
	return nil
}

func (h *MySQLHandler) GetTables() ([]TableInfo, error) {
	var tables []TableInfo
	err := h.DB.Model(TableInfo{}).Raw("SELECT TABLE_NAME, TABLE_COMMENT FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA = DATABASE();").Scan(&tables).Error
	return tables, err
}

func (h *MySQLHandler) GetColumns(table string) ([]ColumnInfo, error) {
	query := fmt.Sprintf(`
		SELECT COLUMN_NAME, COLUMN_TYPE, COLUMN_COMMENT 
		FROM INFORMATION_SCHEMA.COLUMNS 
		WHERE TABLE_NAME = '%s' AND TABLE_SCHEMA = DATABASE();`, table)
	var columns []ColumnInfo
	err := h.DB.Model(ColumnInfo{}).Raw(query).Scan(&columns).Error
	if err != nil {
		return nil, err
	}
	return columns, nil
}

func (h *MySQLHandler) GetSampleData(table, column string) ([]string, error) {
	var samples []string
	query := fmt.Sprintf("SELECT `%s` FROM `%s` LIMIT 1000;", column, table)
	err := h.DB.Raw(query).Scan(&samples).Error
	return samples, err
}

func (h *MySQLHandler) Close() error {
	return nil
}
