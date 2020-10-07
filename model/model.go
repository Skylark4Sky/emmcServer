package model

import (
	. "GoServer/middleWare/dataBases/mysql"
	"GoServer/model/device"
	"GoServer/model/user"
	. "GoServer/utils/log"

	. "GoServer/utils/threadWorker"
	"go.uber.org/zap"
)

type AsyncSQLTaskType uint64

const (
	UNKNOWN_ASYNC_SQL_TASK         AsyncSQLTaskType = iota
	ASYNC_USER_LOGIN_LOG                            //用户登录日志
	ASYNC_UP_USER_AUTH_TIME                         //更新用户授权时间
	ASYNC_DEV_CONNECT_LOG                           //连接日志
	ASYNC_UPDATA_WEUSER_LOCAL                       //更新用户地址
	ASYNC_UPDATA_WEUSER_INFO                        //更新用户资料
	ASYNC_UPDATA_USER_EXTRA                         //更新用户扩展资料
	ASYNC_CREATE_USER_AUTH                          //建立授权记录
	ASYNC_CREATE_USER_REGISTER_LOG                  //建立用户注册日志
	ASYNC_CREATE_USER_EXTRA                         //建立用户信息扩展记录
	ASYNC_CREATE_USER_LOCATION                      //建立用户地址记录
)

type AsyncSQLTask struct {
	Type     AsyncSQLTaskType
	RecordID int64
	Entity   interface{}
}

func CreateAsyncSQLTask(asyncType AsyncSQLTaskType, entity interface{}) {
	var task AsyncSQLTask
	task.Entity = entity
	task.Type = asyncType
	var work Job = &task
	InsertAsyncTask(work)
}

func CreateAsyncSQLTaskWithRecordID(asyncType AsyncSQLTaskType, recordID int64, entity interface{}) {
	var task AsyncSQLTask
	task.Entity = entity
	task.Type = asyncType
	task.RecordID = recordID
	var work Job = &task
	InsertAsyncTask(work)
}

func (task *AsyncSQLTask) ExecTask() error {
	switch task.Type {
	case ASYNC_USER_LOGIN_LOG:
		entity := task.Entity.(user.UserLoginLog)
		if err := ExecSQL().Create(&entity).Error; err != nil {
			SystemLog("add login log Error", zap.Error(err))
		}
		break
	case ASYNC_UP_USER_AUTH_TIME:
		entity := task.Entity.(user.UserAuth)
		if err := ExecSQL().Model(&entity).Update("update_time", entity.UpdateTime).Error; err != nil {
			SystemLog("update auth time Error:", zap.Error(err))
		}
		break
	case ASYNC_DEV_CONNECT_LOG:
		entity := task.Entity.(device.DeviceConnectLog)
		if err := ExecSQL().Create(&entity).Error; err != nil {
			SystemLog("add device connect log Error", zap.Error(err))
		}
		break
	case ASYNC_UPDATA_WEUSER_LOCAL:
		entity := task.Entity.(user.UserLocation)
		if err := ExecSQL().Model(&entity).Where("uid = ?", task.RecordID).Updates(entity).Error; err != nil {
			SystemLog("update UserLocation Error:", zap.Error(err))
		}
		break
	case ASYNC_UPDATA_WEUSER_INFO:
		entity := task.Entity.(user.UserBase)
		if err := ExecSQL().Model(&entity).Where("uid = ?", task.RecordID).Updates(entity).Error; err != nil {
			SystemLog("update userBase Error:", zap.Error(err))
		}
		break
	case ASYNC_UPDATA_USER_EXTRA:
		entity := task.Entity.(user.UserExtra)
		if err := ExecSQL().Model(&entity).Where("uid = ?", task.RecordID).Updates(entity).Error; err != nil {
			SystemLog("update userExtra Error:", zap.Error(err))
		}
		break
	case ASYNC_CREATE_USER_AUTH:
		entity := task.Entity.(user.UserAuth)
		if err := ExecSQL().Create(&entity).Error; err != nil {
			SystemLog("add UserAuth Error", zap.Error(err))
		}
		break
	case ASYNC_CREATE_USER_REGISTER_LOG:
		entity := task.Entity.(user.UserRegisterLog)
		if err := ExecSQL().Create(&entity).Error; err != nil {
			SystemLog("add UserRegisterLog Error", zap.Error(err))
		}
		break
	case ASYNC_CREATE_USER_EXTRA:
		entity := task.Entity.(user.UserExtra)
		if err := ExecSQL().Create(&entity).Error; err != nil {
			SystemLog("add UserExtra Error", zap.Error(err))
		}
		break
	case ASYNC_CREATE_USER_LOCATION:
		entity := task.Entity.(user.UserLocation)
		if err := ExecSQL().Create(&entity).Error; err != nil {
			SystemLog("add UserLocation Error", zap.Error(err))
		}
		break
	}
	return nil
}

func CreateSQLAndRetLastID(value interface{}) (int64, error) {
	var id []int64

	tx := ExecSQL().Begin()
	if err := tx.Debug().Create(value).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	if err := tx.Raw("select LAST_INSERT_ID() as id").Pluck("id", &id).Error; err != nil {
		tx.Rollback()
		return 0, err
	}
	tx.Commit()

	return id[0], nil
}
