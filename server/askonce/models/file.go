package models

import (
	"github.com/pkg/errors"
	"github.com/xiangtao94/golib/flow"
	"github.com/xiangtao94/golib/pkg/orm"

	"gorm.io/gorm"
)

// File  文件
type File struct {
	Id         string `gorm:"id;primaryKey;comment:文件名称"`
	Name       string `gorm:"type:varchar(512);default:'';comment:文件名称"`
	OriginName string `gorm:"type:varchar(512);default:'';comment:文件原始名称"`
	Extension  string `gorm:"ype:varchar(52);default:'';comment:文件格式"`
	Path       string `gorm:"type:varchar(4096);default:'';comment:文件原始路径"`
	Size       int64  `gorm:"type:int(11);default:0;comment:文件大小"`
	Source     string `gorm:"type:varchar(52);default:'';comment:文件来源 "`
	UserId     string `gorm:"type:varchar(128);default:'';comment:用户id"`
	orm.CrudModel
}

func (File) TableName() string {
	return "file"
}

type FileDao struct {
	flow.Dao
}

func (entity *FileDao) OnCreate() {
	entity.SetTable(File{}.TableName())
}

func (entity *FileDao) Insert(add *File) (err error) {
	return entity.GetDB().Create(add).Error
}

func (entity *FileDao) GetById(id string) (res *File, err error) {
	res = &File{}
	err = entity.GetDB().Where("id = ?", id).First(&res).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return
}

func (entity *FileDao) GetByIds(ids []string) (res []*File, err error) {
	if len(ids) == 0 {
		return nil, nil
	}
	err = entity.GetDB().Where("id in ?", ids).Find(&res).Error
	return
}
