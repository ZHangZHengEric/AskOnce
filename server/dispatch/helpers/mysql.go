package helpers

import (
	"dispatch/conf"
	"github.com/xiangtao94/golib/pkg/orm"
	"gorm.io/gorm"
)

var (
	MysqlClient *gorm.DB
)

func InitMysql() {
	var err error
	for name, dbConf := range conf.WebConf.Mysql {
		switch name {
		case "default":
			MysqlClient, err = orm.InitMysqlClient(dbConf)
		}
		if err != nil {
			panic("mysql connect error: %v" + err.Error())
		}
	}
}

func CloseMysql() {
}
