package models

import (
	"github.com/pkg/errors"
	"github.com/xiangtao94/golib/flow"
	"github.com/xiangtao94/golib/pkg/orm"
	"gorm.io/gorm"
)

type KdbDocContent struct {
	DocId   int64  `gorm:"type:int(11);default:0;comment:文档id"`
	KdbId   int64  `gorm:"type:int(11);default:0;comment:知识库id"`
	Content string `gorm:"column:content"` //  文本内容地址
	orm.CrudModel
}

type KdbDocContentDao struct {
	flow.Dao
}

func (entity *KdbDocContentDao) OnCreate() {
	entity.SetTable("kdb_doc_content")
}

func (entity *KdbDocContentDao) Insert(add *KdbDocContent) (err error) {
	return entity.GetDB().Create(add).Error
}

func (entity *KdbDocContentDao) GetByDataIds(dataIds []int64) (res []*KdbDocContent, err error) {
	if len(dataIds) == 0 {
		return
	}
	res = []*KdbDocContent{}
	err = entity.GetDB().Where("doc_id in ? ", dataIds).Find(&res).Error
	if err != nil {
		return nil, err
	}
	return
}

func (entity *KdbDocContentDao) DeleteByDocIds(docIds []int64) (err error) {
	if len(docIds) == 0 {
		return
	}
	err = entity.GetDB().Where("doc_id in ? ", docIds).Delete(&KdbDocContent{}).Error
	if err != nil {
		return err
	}
	return
}
func (entity *KdbDocContentDao) GetByDataId(id int64) (res *KdbDocContent, err error) {
	res = &KdbDocContent{}
	db := entity.GetDB()
	err = db.Where("doc_id = ?", id).First(&res).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return
}

func (entity *KdbDocContentDao) BatchInsert(add []*KdbDocContent) (err error) {
	if len(add) == 0 {
		return nil
	}
	return entity.GetDB().CreateInBatches(add, 2000).Error
}
