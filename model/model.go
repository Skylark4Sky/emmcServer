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
	ASYNC_CREATE_THIRD_USER                         //建立第三方用户数据
	ASYNC_CREATE_NORMAL_USER                        //建立用户数据
	ASYNC_USER_LOGIN_LOG                            //用户登录日志
	ASYNC_UP_USER_AUTH_TIME                         //更新用户授权时间
	ASYNC_MODULE_CONNECT_LOG                        //模组连接日志
	ASYNC_UP_MODULE_VERSION							//更新模组版本
	ASYNC_UP_DEVICE_VERSION							//更新设备版本
	ASYNC_DEV_AND_MODULE_CREATE						//建立设备与模组关系
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

func transactionCreateUserInfo(entity *user.CreateUserInfo, hasAuth bool) error {
	var id []int64
	tx := ExecSQL().Begin()
	if err := tx.Create(&entity.Base).Error; err != nil {
		SystemLog("add UserBase Error", zap.Error(err))
		tx.Rollback()
		return err
	}

	if err := tx.Raw("select LAST_INSERT_ID() as id").Pluck("id", &id).Error; err != nil {
		SystemLog("get LastID Error", zap.Error(err))
		tx.Rollback()
		return err
	}

	var userID int64 = id[0]

	if hasAuth == true {
		entity.Auth.UID = userID
		if err := tx.Create(&entity.Auth).Error; err != nil {
			SystemLog("add UserAuth Error", zap.Error(err))
			tx.Rollback()
			return err
		}
	}

	entity.Log.UID = userID
	if err := tx.Create(&entity.Log).Error; err != nil {
		SystemLog("add UserRegisterLog Error", zap.Error(err))
		tx.Rollback()
		return err
	}

	entity.Extra.UID = userID
	if err := tx.Create(&entity.Extra).Error; err != nil {
		SystemLog("add UserExtra Error", zap.Error(err))
		tx.Rollback()
		return err
	}

	entity.Location.UID = userID
	if err := tx.Create(&entity.Location).Error; err != nil {
		SystemLog("add UserLocation Error", zap.Error(err))
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func transactionCreateDevInfo(entity *device.CreateDeviceInfo) error {

	device := entity.Device
	module := entity.Module
	log := entity.Log

	err := ExecSQL().Where("device_sn = ?", device.DeviceSn).First(&device).Error
	var hasRecord = true

	if err != nil {
		if IsRecordNotFound(err) {
			hasRecord = false
		} else {
			SystemLog("transactionCreateDevInfo select Error", zap.Error(err))
			return err
		}
	}

	if hasRecord {
		//设备已存在，单独建立模组信息
		tx := ExecSQL().Begin()

		module.DeviceID = device.ID
		if err := tx.Create(&module).Error; err != nil {
			SystemLog("add ModuleInfo Error", zap.Error(err))
			tx.Rollback()
			return err
		}
		var id []int64
		if err := tx.Raw("select LAST_INSERT_ID() as id").Pluck("id", &id).Error; err != nil {
			SystemLog("get LastID Error", zap.Error(err))
			tx.Rollback()
			return err
		}

		var ModuleID int64 = id[0]
		log.ModuleID = ModuleID
		if err := tx.Create(&log).Error; err != nil {
			SystemLog("add module connect log Error", zap.Error(err))
			tx.Rollback()
			return err
		}
		tx.Commit()
	} else {
		//事务建立 模组 和 设备信息
		var id []int64
		tx := ExecSQL().Begin()
		if err := tx.Create(&device).Error; err != nil {
			SystemLog("add DeviceInfo Error", zap.Error(err))
			tx.Rollback()
			return err
		}

		if err := tx.Raw("select LAST_INSERT_ID() as id").Pluck("id", &id).Error; err != nil {
			SystemLog("get LastID Error", zap.Error(err))
			tx.Rollback()
			return err
		}

		var DeviceID int64 = id[0]
		module.DeviceID = DeviceID
		if err := tx.Create(&module).Error; err != nil {
			SystemLog("add ModuleInfo Error", zap.Error(err))
			tx.Rollback()
			return err
		}

		if err := tx.Raw("select LAST_INSERT_ID() as id").Pluck("id", &id).Error; err != nil {
			SystemLog("get LastID Error", zap.Error(err))
			tx.Rollback()
			return err
		}

		var ModuleID int64 = id[0]
		log.ModuleID = ModuleID
		if err := tx.Create(&log).Error; err != nil {
			SystemLog("add module connect log Error", zap.Error(err))
			tx.Rollback()
			return err
		}
		tx.Commit()
	}
	return nil
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

	case ASYNC_CREATE_THIRD_USER:
		entity := task.Entity.(user.CreateUserInfo)
		transactionCreateUserInfo(&entity, true)
		break
	case ASYNC_CREATE_NORMAL_USER:
		entity := task.Entity.(user.CreateUserInfo)
		transactionCreateUserInfo(&entity, false)
		break
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
	case ASYNC_MODULE_CONNECT_LOG:
		entity := task.Entity.(device.ModuleConnectLog)
		if err := ExecSQL().Create(&entity).Error; err != nil {
			SystemLog("add module connect log Error", zap.Error(err))
		}
		break
	case ASYNC_UP_MODULE_VERSION:
		entity := task.Entity.(device.ModuleInfo)
		updateMap := map[string]interface{}{"module_version": entity.ModuleVersion, "update_time": entity.UpdateTime}
		if err := ExecSQL().Model(&entity).Updates(updateMap).Error; err != nil {
			SystemLog("update module version Error", zap.Error(err))
		}
		break
	case ASYNC_UP_DEVICE_VERSION:
		entity := task.Entity.(device.DeviceInfo)
		updateMap := map[string]interface{}{"device_version": entity.DeviceVersion, "update_time": entity.UpdateTime}
		if err := ExecSQL().Model(&entity).Updates(updateMap).Error; err != nil {
			SystemLog("update device version Error", zap.Error(err))
		}
		break
	case ASYNC_DEV_AND_MODULE_CREATE:
		entity := task.Entity.(device.CreateDeviceInfo)
		transactionCreateDevInfo(&entity)
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
