// Package data -----------------------------
// @file      : data_soruce.go
// @author    : xiangtao
// @contact   : xiangtao1994@gmail.com
// @time      : 2024/12/22 16:18
// -------------------------------------------
package data

import (
	"askonce/components"
	"askonce/components/dto/dto_kdb_doc"
	"askonce/data/database_parse"
	"askonce/helpers"
	"askonce/models"
	"encoding/json"
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
	databaseHandler, err := database_parse.GetDatabaseHandler(f.GetCtx(), database_parse.DatabaseConfig{
		Driver:   info.DbType,
		Host:     info.DbHost,
		Port:     info.DbPort,
		Database: info.DbName,
		User:     info.DbUser,
		Password: info.DbPwd,
	})
	if err != nil {
		return nil, components.ErrorDbConnError
	}
	defer databaseHandler.Close()
	err = databaseHandler.Ping()
	if err != nil {
		return nil, components.ErrorDbConnError
	}
	schema, err := database_parse.GetSchema(databaseHandler)
	if err != nil {
		return nil, components.ErrorDbSchemaError
	}
	schemaJson, _ := json.Marshal(schema)

	add = &models.Datasource{
		Id:           helpers.GenIDStr(),
		Type:         info.DbType,
		Host:         info.DbHost,
		Port:         info.DbPort,
		Username:     info.DbUser,
		Password:     info.DbPwd,
		DatabaseName: info.DbName,
		JdbcParam:    "",
		Schema:       schemaJson,
		UserId:       userId,
		CrudModel: orm.CrudModel{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	err = f.datasourceDao.Insert(add)
	return
}

func (f *DatasourceData) GetByIds(ids []string) (res map[string]*models.Datasource, err error) {
	if len(ids) == 0 {
		return
	}
	res = make(map[string]*models.Datasource)
	dat, err := f.datasourceDao.GetByIds(ids)
	if err != nil {
		return nil, err
	}
	for _, d := range dat {
		res[d.Id] = d
	}
	return
}

func (f *DatasourceData) DeleteByIds(ids []string) (err error) {
	if len(ids) == 0 {
		return
	}
	return f.datasourceDao.DeleteByIds(ids)
}
