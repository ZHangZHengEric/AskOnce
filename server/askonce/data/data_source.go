// Package data -----------------------------
// @file      : data_soruce.go
// @author    : xiangtao
// @contact   : xiangtao1994@gmail.com
// @time      : 2024/12/22 16:18
// -------------------------------------------
package data

import (
	"askonce/components/dto/dto_kdb_doc"
	"askonce/helpers"
	"askonce/models"
	"github.com/xiangtao94/golib/flow"
	"github.com/xiangtao94/golib/pkg/orm"
	"time"
)

type DatasourceData struct {
	flow.Data
	datasourceDao *models.DatasourceDao
}

func (f *DatasourceData) OnCreate() {
	f.datasourceDao = flow.Create(f.GetCtx(), new(models.DatasourceDao))
}

func (f *DatasourceData) Add(userId string, info dto_kdb_doc.ImportDataBase) (add *models.Datasource, err error) {
	add = &models.Datasource{
		Id:           helpers.GenIDStr(),
		Type:         info.DbType,
		Host:         info.DbHost,
		Port:         info.DbPort,
		Username:     info.DbUser,
		Password:     info.DbPwd,
		DatabaseName: info.DbName,
		JdbcParam:    "",
		UserId:       userId,
		CrudModel: orm.CrudModel{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	err = f.datasourceDao.Insert(add)
	return
}
