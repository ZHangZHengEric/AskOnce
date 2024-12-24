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
	"golang.org/x/sync/errgroup"
	"sync"
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

	schema, err := database_parse.GetSchema(databaseHandler)
	if err != nil {
		return nil, components.ErrorDbSchemaError
	}
	schemaJson, _ := json.Marshal(schema)

	add = &models.Datasource{
		Id:              helpers.GenIDStr(),
		Type:            info.DbType,
		Host:            info.DbHost,
		Port:            info.DbPort,
		Username:        info.DbUser,
		Password:        info.DbPwd,
		DatabaseName:    info.DbName,
		DatabaseComment: info.DbComment,
		JdbcParam:       "",
		Schema:          schemaJson,
		UserId:          userId,
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
func (f *DatasourceData) GetById(id string) (res *models.Datasource, err error) {
	res, err = f.datasourceDao.GetById(id)
	if err != nil {
		return nil, err
	}
	return
}

func (f *DatasourceData) UpdateDatasource(id string) (res *models.Datasource, err error) {
	res, err = f.datasourceDao.GetById(id)
	if err != nil {
		return nil, err
	}
	return
}
func (f *DatasourceData) DeleteByIds(ids []string) (err error) {
	if len(ids) == 0 {
		return
	}
	return f.datasourceDao.DeleteByIds(ids)
}

func (f *DatasourceData) GetSchemaAndValues(datasourceId string) (datasource *models.Datasource, schemas []database_parse.TableColumnInfo, err error) {
	datasource, err = f.datasourceDao.GetById(datasourceId)
	if err != nil {
		return
	}
	if datasource == nil {
		return nil, nil, components.ErrorDbSchemaError
	}
	databaseHandler, err := database_parse.GetDatabaseHandler(f.GetCtx(), database_parse.DatabaseConfig{
		Driver:   datasource.Type,
		Host:     datasource.Host,
		Port:     datasource.Port,
		Database: datasource.DatabaseName,
		User:     datasource.Username,
		Password: datasource.Password,
	})
	if err != nil {
		return nil, nil, components.ErrorDbConnError
	}
	schemas, err = database_parse.GetSchema(databaseHandler)
	if err != nil {
		return nil, nil, components.ErrorDbSchemaError
	}
	schemaJson, _ := json.Marshal(schemas)
	err = f.datasourceDao.Update(datasource.Id, map[string]interface{}{"schema": schemaJson})
	if err != nil {
		return nil, nil, components.ErrorDbSchemaError
	}
	valueInfoMap := make(map[string][]database_parse.ColumnValueInfo, 0)
	wg, _ := errgroup.WithContext(f.GetCtx())
	lock := sync.RWMutex{}
	for _, schema := range schemas {
		for _, column := range schema.ColumnInfos {
			if column.ColumnType == "varchar" || column.ColumnType == "text" || column.ColumnType == "longtext" {
				wg.Go(func() error {
					datas, err := databaseHandler.GetSampleData(schema.TableName, column.ColumnName)
					if err != nil {
						return err
					}
					lock.Lock()
					tmp := make([]database_parse.ColumnValueInfo, 0, len(datas))
					for _, d := range datas {
						if len(d) == 0 {
							continue
						}
						tmp = append(tmp, database_parse.ColumnValueInfo{
							Value: d,
						})
					}
					valueInfoMap[schema.TableName+column.ColumnName] = tmp
					lock.Unlock()
					return nil
				})
			}
		}
	}
	if err := wg.Wait(); err != nil {
		return nil, nil, err
	}
	for i := range schemas {
		for j := range schemas[i].ColumnInfos {
			key := schemas[i].TableName + schemas[i].ColumnInfos[j].ColumnName
			if v, ok := valueInfoMap[key]; ok {
				schemas[i].ColumnInfos[j].ColumnValues = v
			}
		}
	}
	return
}
