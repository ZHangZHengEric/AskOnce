package service

import (
	"askonce/components/dto"
	"askonce/components/dto/dto_config"
	"askonce/models"
	"askonce/utils"
	"encoding/json"
	"github.com/xiangtao94/golib/flow"
)

type ConfigService struct {
	flow.Service
}

func (entity *ConfigService) Detail(req *dto_config.DetailReq) (res *dto_config.ConfigResp, err error) {
	userInfo, err := utils.LoginInfo(entity.GetCtx())
	if err != nil {
		return nil, err
	}
	userDao := entity.Create(new(models.UserDao)).(*models.UserDao)
	user, err := userDao.GetByUserId(userInfo.UserId)
	if err != nil {
		return nil, err
	}
	userSetting := user.Setting.Data()
	res = &dto_config.ConfigResp{
		Language:  "zh-cn",
		ModelType: "",
	}
	if userSetting.ModelType == "" {
		return
	}
	res.Language = userSetting.Language
	res.ModelType = userSetting.ModelType
	return
}

func (entity *ConfigService) Save(req *dto_config.SaveReq) (res interface{}, err error) {
	userInfo, err := utils.LoginInfo(entity.GetCtx())
	if err != nil {
		return nil, err
	}
	userDao := entity.Create(new(models.UserDao)).(*models.UserDao)

	user, err := userDao.GetByUserId(userInfo.UserId)
	if err != nil {
		return nil, err
	}
	userSetting := user.Setting.Data()
	userSetting.Language = req.Language
	userSetting.ModelType = req.ModelType
	tmpStr, _ := json.Marshal(&userSetting)
	err = userDao.Update(userInfo.UserId, map[string]interface{}{"setting": tmpStr})
	if err != nil {
		return nil, err
	}
	return
}

func (entity *ConfigService) Dict(req *dto.EmptyReq) (res map[string][]dto_config.Dict, err error) {
	res = make(map[string][]dto_config.Dict)
	res["language"] = []dto_config.Dict{
		{
			Name:   "中文",
			EnName: "Chinese",
			Value:  "zh-cn",
		},
		{
			Name:   "英文",
			EnName: "English",
			Value:  "en-us",
		},
	}
	res["models"] = []dto_config.Dict{
		{
			Name:   "OpenAI",
			EnName: "OpenAI",
			Value:  "chatgpt-4-minio",
		},
	}
	return
}
