package models

import (
	"askonce/api/jobd"
	"github.com/xiangtao94/golib/flow"
	"github.com/xiangtao94/golib/pkg/orm"
	"gorm.io/datatypes"
)

type KdbDocSegment struct {
	Id         int64                                `gorm:"id; primaryKey;autoIncrement;comment:自增主键"`
	DocId      int64                                `gorm:"type:int(11);default:0;comment:文档id"`
	KdbId      int64                                `gorm:"type:int(11);default:0;comment:知识库id"`
	PageIndex  int                                  `gorm:"type:int(11);default:0;comment:页码下标"`
	StartIndex int                                  `gorm:"type:int(11);default:0;comment:开始下标"`
	EndIndex   int                                  `gorm:"type:int(11);default:0;comment:结束下标"`
	Text       string                               `gorm:"default:'';comment:文本内容"`
	Sub        datatypes.JSONSlice[jobd.TextDetail] `gorm:"column:sub;comment:内部详情"`
	orm.CrudModel
}

type KdbDocSegmentDao struct {
	flow.Dao
}

func (entity *KdbDocSegmentDao) OnCreate() {
	entity.SetTable("kdb_doc_segment")
	entity.SetPartitionNum(5)
}

func (entity *KdbDocSegmentDao) BatchInsert(kdbId int64, ms []*KdbDocSegment) error {
	if len(ms) == 0 {
		return nil
	}
	return entity.GetDB().Table(entity.GetPartitionTable(kdbId)).CreateInBatches(ms, 4000).Error
}

func (entity *KdbDocSegmentDao) DeleteByDocIds(kdbId int64, docIds []int64) (err error) {
	if len(docIds) == 0 {
		return
	}
	err = entity.GetDB().Table(entity.GetPartitionTable(kdbId)).Where("doc_id in ? ", docIds).Delete(&KdbDocSegment{}).Error
	if err != nil {
		return err
	}

	return
}

func (entity *KdbDocSegmentDao) GetByDataIds(kdbId int64, dataIds []int64) (res []*KdbDocSegment, err error) {
	if len(dataIds) == 0 {
		return
	}
	res = []*KdbDocSegment{}
	err = entity.GetDB().Table(entity.GetPartitionTable(kdbId)).Where("doc_id in ? ", dataIds).Find(&res).Error
	if err != nil {
		return nil, err
	}
	return
}

func (entity *KdbDocSegmentDao) GetByDataId(kdbId int64, dataId int64) (res []*KdbDocSegment, err error) {
	var tmp []*KdbDocSegment
	err = entity.GetDB().Table(entity.GetPartitionTable(kdbId)).Where("doc_id = ? ", dataId).Find(&tmp).Error
	if err != nil {
		return nil, err
	}
	return tmp, nil
}
