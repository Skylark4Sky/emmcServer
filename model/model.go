package model

import (
	. "GoServer/middleWare/dataBases/mysql"
	. "GoServer/middleWare/dataBases/redis"
	"GoServer/model/device"
	"GoServer/model/user"
	. "GoServer/utils/log"
	. "GoServer/utils/threadWorker"
	"go.uber.org/zap"
	"reflect"
)

type AsyncSQLTaskType uint64

const (
	UNKNOWN_ASYNC_SQL_TASK         AsyncSQLTaskType = iota
	ASYNC_CREATE_THIRD_USER                         //建立第三方用户数据 1
	ASYNC_CREATE_NORMAL_USER                        //建立用户数据 2
	ASYNC_USER_LOGIN_LOG                            //用户登录日志 3
	ASYNC_UP_USER_AUTH_TIME                         //更新用户授权时间 4
	ASYNC_MODULE_CONNECT_LOG                        //模组连接日志 5
	ASYNC_UP_MODULE_INFO                            //更新模组版本 6
	ASYNC_UP_DEVICE_INFO                            //更新设备版本 7
	ASYNC_DEV_AND_MODULE_CREATE                     //建立设备与模组关系 8
	ASYNC_UPDATA_WEUSER_LOCAL                       //更新用户地址 9
	ASYNC_UPDATA_WEUSER_INFO                        //更新用户资料 10
	ASYNC_UPDATA_USER_EXTRA                         //更新用户扩展资料 11
	ASYNC_CREATE_USER_AUTH                          //建立授权记录 12
	ASYNC_CREATE_USER_REGISTER_LOG                  //建立用户注册日志 13
	ASYNC_CREATE_USER_EXTRA                         //建立用户信息扩展记录 14
	ASYNC_CREATE_USER_LOCATION                      //建立用户地址记录 15
	ASYNC_UPDATE_DEVICE_STATUS                      //更新设备状态 16
	ASYNC_CREATE_COM_CHARGE_TASK					//建立充电记录
)

type TaskFunc func(task  *AsyncSQLTask)

type AsyncSQLTask struct {
	Type      AsyncSQLTaskType
	WhereSQL  string
	RecordID  int64
	Lock	  *RedisLock
	Func 	  TaskFunc
	MapParam map[string]interface{}
	Entity    interface{}
}

func CreateSQLAndRetLastID(entity interface{}) (uint64, error) {
	var id []uint64

	tx := ExecSQL().Begin()
	if err := tx.Create(entity).Error; err != nil {
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

func CreateAsyncSQLTask(asyncType AsyncSQLTaskType, entity interface{}) {
	var task AsyncSQLTask
	task.Entity = entity
	task.Type = asyncType
	var work Job = &task
	InsertAsyncTask(work)
}

//根据map更新
func CreateAsyncSQLTaskWithUpdateMap(asyncType AsyncSQLTaskType, entity interface{}, mapParam map[string]interface{}) {
	var task AsyncSQLTask
	task.Entity = entity
	task.Type = asyncType
	task.MapParam = mapParam
	var work Job = &task
	InsertAsyncTask(work)
}

func CreateAsyncSQLTaskWithCallback(asyncType AsyncSQLTaskType, entity interface{},lock *RedisLock, taskFunc TaskFunc) {
	var task AsyncSQLTask
	task.Entity = entity
	task.Type = asyncType
	task.Func = taskFunc
	task.Lock = lock
	var work Job = &task
	InsertAsyncTask(work)
}

//根据uid 更新
func CreateAsyncSQLTaskWithRecordID(asyncType AsyncSQLTaskType, recordID int64, entity interface{}) {
	var task AsyncSQLTask
	task.Entity = entity
	task.Type = asyncType
	task.WhereSQL = "uid = ?"
	task.RecordID = recordID
	var work Job = &task
	InsertAsyncTask(work)
}

func transactionCreateUserInfo(entity *user.CreateUserInfo, hasAuth bool) error {
	var id []uint64
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

	var userID uint64 = id[0]

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

func updateDeviceIDToRedisByDeviceSN(deviceSN string, deviceID uint64) {
	if deviceSN != "" && deviceID != 0 {
		Redis().InitWithInsertDeviceIDToken(deviceSN, deviceID)
	}
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
		var id []uint64
		if err := tx.Raw("select LAST_INSERT_ID() as id").Pluck("id", &id).Error; err != nil {
			SystemLog("get LastID Error", zap.Error(err))
			tx.Rollback()
			return err
		}

		var ModuleID uint64 = id[0]
		log.ModuleID = ModuleID
		if err := tx.Create(&log).Error; err != nil {
			SystemLog("add module connect log Error", zap.Error(err))
			tx.Rollback()
			return err
		}
		tx.Commit()
	} else {
		//事务建立 模组 和 设备信息
		var id []uint64
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

		var DeviceID uint64 = id[0]
		device.ID = DeviceID
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

		var ModuleID uint64 = id[0]
		log.ModuleID = ModuleID
		if err := tx.Create(&log).Error; err != nil {
			SystemLog("add module connect log Error", zap.Error(err))
			tx.Rollback()
			return err
		}
		tx.Commit()
	}

	updateDeviceIDToRedisByDeviceSN(device.DeviceSn, device.ID)

	return nil
}

func findComChargeTaskRecord(entity *device.DeviceCom) (bool,error) {
	err := ExecSQL().Where("device_id = ? AND charge_id = ? AND com_id = ?", entity.DeviceID,entity.ChargeID, entity.ComID).Order("create_time desc").First(&entity).Error
	var hasRecord = true

	if err != nil {
		if IsRecordNotFound(err) {
			hasRecord = false
		} else {
			SystemLog("createComChargeTaskRecord select Error", zap.Error(err))
			return true,err
		}
	}
	return hasRecord,nil
}

func createComChargeTaskRecord(entity *device.DeviceCom) error {

	taskRecord := &device.DeviceCom {
		DeviceID: entity.DeviceID,
		ChargeID: entity.ChargeID,
		ComID: entity.ComID,
	}

	hasRecord,err := findComChargeTaskRecord(taskRecord)
	if (err != nil) {
		return err
	}

	//存在记录
	if hasRecord {
		taskRecord.MaxEnergy = entity.MaxEnergy
		taskRecord.MaxTime = entity.MaxTime
		taskRecord.MaxElectricity = entity.MaxElectricity
		if err := ExecSQL().Update(taskRecord).Error; err != nil {
			structTpey := reflect.Indirect(reflect.ValueOf(taskRecord)).Type()
			SystemLog("Updatte ", structTpey, " Error ", zap.Any("SQL", taskRecord), zap.Error(err))
		}
	} else { //不存在记录
		if err := ExecSQL().Create(entity).Error; err != nil {
			structTpey := reflect.Indirect(reflect.ValueOf(entity)).Type()
			SystemLog("Create ", structTpey, " Error ", zap.Any("SQL", entity), zap.Error(err))
		}
	}
	return nil
}

func (task *AsyncSQLTask) ExecTask() error {
	switch task.Type {

	case ASYNC_CREATE_THIRD_USER:
		entity := task.Entity.(user.CreateUserInfo)
		transactionCreateUserInfo(&entity, true)
	case ASYNC_CREATE_NORMAL_USER:
		entity := task.Entity.(user.CreateUserInfo)
		transactionCreateUserInfo(&entity, false)
	case ASYNC_UP_USER_AUTH_TIME, ASYNC_UP_DEVICE_INFO, ASYNC_UP_MODULE_INFO:
		if err := ExecSQL().Model(task.Entity).Updates(task.MapParam).Error; err != nil {
			SystemLog("update Data Error:", zap.Any("SQL", task.Entity), zap.Error(err))
		}
		if ASYNC_UP_DEVICE_INFO == task.Type {
			entity := task.Entity.(*device.DeviceInfo)
			updateDeviceIDToRedisByDeviceSN(entity.DeviceSn, entity.ID)
		}
	case ASYNC_DEV_AND_MODULE_CREATE:
		entity := task.Entity.(device.CreateDeviceInfo)
		transactionCreateDevInfo(&entity)
	case ASYNC_UPDATA_WEUSER_LOCAL, ASYNC_UPDATA_WEUSER_INFO, ASYNC_UPDATA_USER_EXTRA:
		if err := ExecSQL().Model(task.Entity).Where(task.WhereSQL, task.RecordID).Updates(task.Entity).Error; err != nil {
			structTpey := reflect.Indirect(reflect.ValueOf(task.Entity)).Type()
			SystemLog("Update ", structTpey, " Error ", zap.Any("SQL", task.Entity), zap.Error(err))
		}
	case ASYNC_USER_LOGIN_LOG, ASYNC_CREATE_USER_REGISTER_LOG, ASYNC_MODULE_CONNECT_LOG,
		ASYNC_CREATE_USER_AUTH, ASYNC_CREATE_USER_EXTRA, ASYNC_CREATE_USER_LOCATION:
		if err := ExecSQL().Create(task.Entity).Error; err != nil {
			structTpey := reflect.Indirect(reflect.ValueOf(task.Entity)).Type()
			SystemLog("Create ", structTpey, " Error ", zap.Any("SQL", task.Entity), zap.Error(err))
		}
	case ASYNC_UPDATE_DEVICE_STATUS:
		if(task.Func != nil) {
			task.Func(task)
		}
	case ASYNC_CREATE_COM_CHARGE_TASK:
		entity := task.Entity.(*device.DeviceCom)
		createComChargeTaskRecord(entity)
	}
	return nil
}
