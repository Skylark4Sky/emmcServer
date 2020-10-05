package model

import (
	. "GoServer/middleWare/dataBases/mysql"
	"GoServer/model/device"
	"GoServer/model/user"
	. "GoServer/utils/log"
	"go.uber.org/zap"
)

type AsynSQLTask struct {
	Entity interface{}
}

func (task *AsynSQLTask) ExecTask() error {
	switch entity := task.Entity.(type) {
	case user.UserLoginLog:
		if err := ExecSQL().Create(&entity).Error; err != nil {
			SystemLog("add login log Error", zap.Error(err))
		}
	case user.UserAuth:
		if err := ExecSQL().Model(&entity).Update("update_time", entity.UpdateTime).Error; err != nil {
			SystemLog("update auth time Error:", zap.Error(err))
		}
	case device.DeviceConnectLog:
		if err := ExecSQL().Create(&entity).Error; err != nil {
			SystemLog("add device connect log Error", zap.Error(err))
		}
	}
	return nil
}
