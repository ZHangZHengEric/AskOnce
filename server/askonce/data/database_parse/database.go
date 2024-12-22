// Package handler -----------------------------
// @file      : database.go
// @author    : xiangtao
// @contact   : xiangtao1994@gmail.com
// @time      : 2024/12/22 14:56
// -------------------------------------------
package database_parse

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
)

// DatabaseHandler 定义通用数据库操作接口
type DatabaseHandler interface {
	Ping() error
	GetTables() ([]TableInfo, error)
	GetColumns(table string) ([]ColumnInfo, error)
	GetSampleData(table, column string) ([]string, error)
	Close() error
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

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Driver   string
	Host     string
	Port     int
	Database string
	User     string
	Password string
}

func GetDatabaseHandler(config DatabaseConfig) (DatabaseHandler, error) {
	var handler DatabaseHandler
	defer handler.Close()
	switch config.Driver {
	case "mysql":
		mysqlDB, err := gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", config.User, config.Password, config.Host, config.Port, config.Database)),
			&gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
		if err != nil {
			return nil, err
		}
		handler = &MySQLHandler{DB: mysqlDB}
	case "postgresql":
		postgresDB, err := sqlx.Connect(config.Driver, fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", config.User, config.Password, config.Host, config.Port, config.Database))
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}
		handler = &PostgreSQLHandler{DB: postgresDB}
	default:
		return nil, fmt.Errorf("unknown database driver: %s", config.Driver)
	}
	return handler, nil
}
