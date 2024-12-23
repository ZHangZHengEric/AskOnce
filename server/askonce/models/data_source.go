// Package models -----------------------------
// @file      : data_source.go
// @author    : xiangtao
// @contact   : xiangtao1994@gmail.com
// @time      : 2024/12/22 16:14
// -------------------------------------------
package models

import (
	"askonce/components"
	"github.com/pkg/errors"
	"github.com/xiangtao94/golib/flow"
	"github.com/xiangtao94/golib/pkg/orm"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

// Datasource represents the data source configuration in the system.
type Datasource struct {
	Id           string         `gorm:"id;primaryKey;comment:主键id"`
	Type         string         `gorm:"type:varchar(52);not null;column:type" json:"type"`                            // 数据源类型
	Host         string         `gorm:"type:varchar(255);not null;column:host" json:"host"`                           // 数据库地址
	Port         int            `gorm:"not null;column:port" json:"port"`                                             // 数据库端口
	Username     string         `gorm:"type:varchar(255);not null;column:username" json:"username"`                   // 数据库用户名
	Password     string         `gorm:"type:varchar(255);not null;column:password" json:"password"`                   // 数据库密码
	DatabaseName string         `gorm:"type:varchar(255);not null;column:database_name" json:"database_name"`         // 数据库名称
	JdbcParam    string         `gorm:"type:varchar(500);default:null;column:jdbc_param" json:"jdbc_param,omitempty"` // JDBC 连接字符串（可选）
	UserId       string         `gorm:"type:varchar(128);default:'';column:user_id" json:"user_id"`                   // 用户ID
	Schema       datatypes.JSON `gorm:"type:json;;column:schema" json:"schema"`                                       // 结构
	orm.CrudModel
}

// TableName overrides the default table name.
func (Datasource) TableName() string {
	return "data_source"
}

type DatasourceDao struct {
	flow.Dao
}

func (entity *DatasourceDao) OnCreate() {
	entity.SetTable(Datasource{}.TableName())
}

func (entity *DatasourceDao) Insert(add *Datasource) (err error) {
	return entity.GetDB().Create(add).Error
}

func (entity *DatasourceDao) GetById(id string) (res *Datasource, err error) {
	res = &Datasource{}
	err = entity.GetDB().Where("id = ?", id).First(&res).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return
}

func (entity *DatasourceDao) GetByIds(ids []string) (res []*Datasource, err error) {
	if len(ids) == 0 {
		return nil, nil
	}
	err = entity.GetDB().Where("id in ?", ids).Find(&res).Error
	return
}

func (entity *DatasourceDao) DeleteByIds(ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	return entity.GetDB().Where("id in ?", ids).Delete(&Datasource{}).Error
}

// 更新
func (entity *DatasourceDao) Update(id string, update map[string]interface{}) error {
	update["updated_at"] = time.Now()
	err := entity.GetDB().Where("id = ?", id).Updates(update).Error
	if err != nil {
		return components.ErrorMysqlError
	}
	return nil
}
