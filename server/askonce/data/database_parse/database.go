// Package handler -----------------------------
// @file      : database.go
// @author    : xiangtao
// @contact   : xiangtao1994@gmail.com
// @time      : 2024/12/22 14:56
// -------------------------------------------
package database_parse

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"sync"
	"time"
)

// DatabaseHandler 定义通用数据库操作接口
type DatabaseHandler interface {
	SetCtx(ctx *gin.Context) error
	GetCtx() *gin.Context
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

type TableColumnInfo struct {
	TableInfo
	ColumnInfos []ColumnInfo `json:"column_infos" db:"-"`
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

func GetDatabaseHandler(ctx *gin.Context, config DatabaseConfig) (DatabaseHandler, error) {
	var handler DatabaseHandler
	switch config.Driver {
	case "mysql":
		mysqlDB, err := gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?timeout=%s", config.User, config.Password, config.Host, config.Port, config.Database, 3*time.Second)),
			&gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
		if err != nil {
			return nil, err
		}
		handler = &MySQLHandler{DB: mysqlDB}
	case "postgresql":

		postgresDB, err := gorm.Open(postgres.New(postgres.Config{
			DSN:                  fmt.Sprintf("host=%s user=%S password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", config.Host, config.User, config.Password, config.Database, config.Port), // data source name, refer https://github.com/jackc/pgx
			PreferSimpleProtocol: true,                                                                                                                                                                         // disables implicit prepared statement usage. By default pgx automatically uses the extended protocol
		}), &gorm.Config{})
		if err != nil {
			return nil, err
		}
		handler = &PostgreSQLHandler{DB: postgresDB}
	default:
		return nil, fmt.Errorf("unknown database driver: %s", config.Driver)
	}
	handler.SetCtx(ctx)
	return handler, nil
}

// 获取所有的表结构
func GetSchema(handler DatabaseHandler) (tableColumns []TableColumnInfo, err error) {
	tableColumns = make([]TableColumnInfo, 0)
	tables, err := handler.GetTables()
	if err != nil {
		return nil, err
	}
	wg, _ := errgroup.WithContext(handler.GetCtx())
	lock := sync.RWMutex{}
	for _, table := range tables {
		wg.Go(func() error {
			lock.Lock()
			columns, err := handler.GetColumns(table.TableName)
			if err != nil {
				return err
			}
			tableColumns = append(tableColumns, TableColumnInfo{
				TableInfo:   table,
				ColumnInfos: columns,
			})
			lock.Unlock()
			return nil
		})
	}
	if err := wg.Wait(); err != nil {
		return nil, err
	}
	return tableColumns, nil
}