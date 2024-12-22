// Package handler -----------------------------
// @file      : database.go
// @author    : xiangtao
// @contact   : xiangtao1994@gmail.com
// @time      : 2024/12/22 14:56
// -------------------------------------------
package database_parse

// DatabaseHandler 定义通用数据库操作接口
type DatabaseHandler interface {
	GetTables() ([]TableInfo, error)
	GetColumns(table string) ([]ColumnInfo, error)
	GetSampleData(table, column string) ([]string, error)
}

// 定义结构体用于接收结果
type ColumnInfo struct {
	ColumnName    string `gorm:"column:COLUMN_NAME" json:"column_name" db:"COLUMN_NAME"`
	ColumnType    string `gorm:"column:COLUMN_TYPE" json:"column_type" db:"COLUMN_TYPE"`
	ColumnComment string `gorm:"column:COLUMN_COMMENT" json:"column_comment" db:"COLUMN_COMMENT"`
}

type TableInfo struct {
	TableName    string `json:"table_name" gorm:"column:TABLE_NAME" db:"TABLE_NAME"`
	TableComment string `json:"table_comment" gorm:"column:TABLE_COMMENT" db:"TABLE_COMMENT"`
}
