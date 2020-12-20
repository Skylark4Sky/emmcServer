package asyncTask

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

type AsyncTaskType uint64

const (
	UNKNOWN_ASYNC_TASK                    AsyncTaskType = iota
	ASYNC_CREATE_THIRD_USER                             //建立第三方用户数据 1
	ASYNC_CREATE_NORMAL_USER                            //建立用户数据 2
	ASYNC_USER_LOGIN_LOG                                //用户登录日志 3
	ASYNC_UP_USER_AUTH_TIME                             //更新用户授权时间 4
	ASYNC_MODULE_CONNECT_LOG                            //模组连接日志 5
	ASYNC_UP_MODULE_INFO                                //更新模组版本 6
	ASYNC_UP_DEVICE_INFO                                //更新设备版本 7
	ASYNC_DEV_AND_MODULE_CREATE                         //建立设备与模组关系 8
	ASYNC_UPDATA_WEUSER_LOCAL                           //更新用户地址 9
	ASYNC_UPDATA_WEUSER_INFO                            //更新用户资料 10
	ASYNC_UPDATA_USER_EXTRA                             //更新用户扩展资料 11
	ASYNC_CREATE_USER_AUTH                              //建立授权记录 12
	ASYNC_CREATE_USER_REGISTER_LOG                      //建立用户注册日志 13
	ASYNC_CREATE_USER_EXTRA                             //建立用户信息扩展记录 14
	ASYNC_CREATE_USER_LOCATION                          //建立用户地址记录 15
	ASYNC_UPDATE_DEVICE_STATUS                          //更新设备状态 16
	ASYNC_CREATE_COM_CHARGE_TASK                        //建立充电记录
	ASYNC_CREATE_COM_CHARGE_TASK_ACK                    //设备上报开始充电
	ASYNC_STOP_COM_CHARGE_TASK                          //退出充电
	ASYNC_STOP_COM_CHARGE_TASK_ACK                      //设备响应上报退出充电
	ASYNC_INITIATIVE_EXIT_COM_CHARGE_TASK               //主动退出充电
	ASYNC_UPDATE_CHARGE_TASK_DATA                       //定时同步数据
	ASYNC_CHK_CHARGE_TASK_DATA                          //检查数据
)

type AsyncTaskFunc func(task *AsyncTaskEntity)

type AsyncTaskEntity struct {
	Type   AsyncTaskType
	Lock   *RedisLock
	Func   AsyncTaskFunc
	Param  map[string]interface{}
	Entity interface{}
}

func NewTask() *AsyncTaskEntity {
	return &AsyncTaskEntity{
		Type:   UNKNOWN_ASYNC_TASK,
		Lock:   nil,
		Func:   nil,
		Param:  nil,
		Entity: nil,
	}
}

func NewAsyncTaskWithParam(typeVal AsyncTaskType, entity interface{}) {
	task := NewTask()
	task.RunTaskWithTypeAndEntity(typeVal, entity)
}

func (task *AsyncTaskEntity) RunTaskWithTypeAndEntity(typeVal AsyncTaskType, entity interface{}) {
	task.SetType(typeVal)
	task.SetEntity(entity)
	task.InsertWorkerQueue()
}

func (task *AsyncTaskEntity) SetType(typeVal AsyncTaskType) {
	if typeVal != UNKNOWN_ASYNC_TASK {
		task.Type = typeVal
	}
}

func (task *AsyncTaskEntity) SetEntity(entity interface{}) {
	if entity != nil {
		task.Entity = entity
	}
}

func (task *AsyncTaskEntity) SetParam(param map[string]interface{}) {
	if param != nil {
		task.Param = param
	}
}

func (task *AsyncTaskEntity) SetTaskFunc(taskFunc AsyncTaskFunc) {
	if taskFunc != nil {
		task.Func = taskFunc
	}
}

func (task *AsyncTaskEntity) SetLock(lock *RedisLock) {
	if lock != nil {
		task.Lock = lock
	}
}

func (task *AsyncTaskEntity) InsertWorkerQueue() {
	var work Job = task
	InsertAsyncTask(work)
}

func (task *AsyncTaskEntity) ExecTask() error {
	if task != nil {
		if task.Func != nil {
			task.Func(task)
			return nil
		} else {
			switch task.Type {
			case ASYNC_CREATE_THIRD_USER, ASYNC_CREATE_NORMAL_USER:
				entity := task.Entity.(user.CreateUserInfo)
				if task.Type == ASYNC_CREATE_NORMAL_USER {
					user.CreateNewUser(&entity, false)
				} else {
					user.CreateNewUser(&entity, true)
				}
			case ASYNC_UP_USER_AUTH_TIME, ASYNC_UP_DEVICE_INFO, ASYNC_UP_MODULE_INFO:
				if err := ExecSQL().Model(task.Entity).Updates(task.Param).Error; err != nil {
					SystemLog("update Data Error:", zap.Any("SQL", task.Entity), zap.Error(err))
				}
				if ASYNC_UP_DEVICE_INFO == task.Type {
					entity := task.Entity.(device.DeviceInfo)
					Redis().UpdateDeviceIDToRedisByDeviceSN(entity.DeviceSn, entity.ID, entity.UID)
				}
			case ASYNC_DEV_AND_MODULE_CREATE:
				entity := task.Entity.(device.CreateDeviceInfo)
				if err := device.CreateDeviceAndModuleInfo(&entity); err == nil {
					entity.Log.ModuleID = entity.Module.ID
					entity.Log.UID = entity.Module.UID
					if err = ExecSQL().Create(&entity.Log).Error; err != nil {
						SystemLog("Create Module Connect Error ", zap.Any("SQL", entity.Log), zap.Error(err))
					}
					//刷新DeviceRedis
					Redis().UpdateDeviceIDToRedisByDeviceSN(entity.Device.DeviceSn, entity.Device.ID, entity.Device.UID)
				}
			case ASYNC_UPDATA_WEUSER_LOCAL, ASYNC_UPDATA_WEUSER_INFO, ASYNC_UPDATA_USER_EXTRA:
				if WhereSQL, ok := task.Param["WhereSQL"]; ok {
					if err := ExecSQL().Model(task.Entity).Where(WhereSQL.(string)).Updates(task.Entity).Error; err != nil {
						structTpey := reflect.Indirect(reflect.ValueOf(task.Entity)).Type()
						SystemLog("Update ", structTpey, " Error ", zap.Any("SQL", task.Entity), zap.Error(err))
					}
				}
			case ASYNC_USER_LOGIN_LOG, ASYNC_CREATE_USER_REGISTER_LOG, ASYNC_MODULE_CONNECT_LOG,
				ASYNC_CREATE_USER_AUTH, ASYNC_CREATE_USER_EXTRA, ASYNC_CREATE_USER_LOCATION:
				if err := ExecSQL().Create(task.Entity).Error; err != nil {
					structTpey := reflect.Indirect(reflect.ValueOf(task.Entity)).Type()
					SystemLog("Create ", structTpey, " Error ", zap.Any("SQL", task.Entity), zap.Error(err))
				}
			}
		}
	}
	return nil
}
