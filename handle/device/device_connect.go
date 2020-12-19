package device

import (
	. "GoServer/middleWare/dataBases/mysql"
	. "GoServer/model/asyncTask"
	. "GoServer/model/device"
	. "GoServer/utils/respond"
	. "GoServer/utils/time"
	"github.com/gin-gonic/gin"
)

const (
	CHARGEING_TIME = 60  //忙时150秒
	LEISURE_TIME   = 120 //空闲时300秒
)

const (
	DEVICE_OFFLINE int8 = 0
	DEVICE_ONLINE  int8 = 1
	DEVICE_WORKING int8 = 2
)

type RequestData struct {
	AccessWay     uint8  `form:"type" json:"type" binding:"required"`
	ModuleSN      string `form:"module_sn" json:"module_sn" binding:"required"`
	ModuleVersion string `form:"module_version" json:"module_version" binding:"required"`
	DeviceSN      string `form:"device_sn" json:"device_sn" binding:"required"`
	DeviceVersion string `form:"device_version" json:"device_version" binding:"required"`
	Token         string `form:"token" json:"token" binding:"required"`
}

type FirmwareInfo struct {
	Size int64  `json:"size"`
	URL  string `json:"url"`
}

func createConnectLog(ctx *gin.Context, module_id, userID uint64, accessway uint8, moduleSN string) {
	log := &ModuleConnectLog{}
	log.Create(module_id, accessway, moduleSN, ctx.ClientIP())
	log.UID = userID
	NewAsyncTaskWithParam(ASYNC_MODULE_CONNECT_LOG, log)
}

func (data *RequestData) Connect(ctx *gin.Context) interface{} {
	var info CreateDeviceInfo

	info.Type = HAS_DEVICE_WITH_MODULE

	if err := ExecSQL().Where("device_sn = ?", data.DeviceSN).First(&info.Device).Error; err != nil {
		if IsRecordNotFound(err) {
			info.Type = (info.Type &^ (DEVICE_BUILD_BIT))
		} else {
			return CreateMessage(SYSTEM_ERROR, err)
		}
	}

	if err := ExecSQL().Where("module_sn = ?", data.ModuleSN).First(&info.Module).Error; err != nil {
		if IsRecordNotFound(err) {
			info.Type = (info.Type &^ (MODULE_BUILD_BIT))
		} else {
			return CreateMessage(SYSTEM_ERROR, err)
		}
	}

	var respond interface{} = nil
	curTimestampMs := GetTimestampMs()

	switch info.Type {
	case NO_DEVICE_WITH_MODULE:
		//创建对应关系 第一次连接仅建立关系，版本检测下次连接时在进行处理
		info.Module.Create(data.AccessWay, data.ModuleSN, data.ModuleVersion)
		info.Device.Create(data.AccessWay, data.DeviceSN, data.DeviceVersion, DEVICE_ONLINE)
		info.Log.Create(0, data.AccessWay, data.ModuleSN, ctx.ClientIP())
		info.Module.CreateTime = curTimestampMs
		info.Device.CreateTime = curTimestampMs
		NewAsyncTaskWithParam(ASYNC_DEV_AND_MODULE_CREATE, info)
	case DEVICE_BUILD_BIT: //需单独创建Module
		info.Module.Create(data.AccessWay, data.ModuleSN, data.ModuleVersion)
		info.Log.Create(0, data.AccessWay, data.ModuleSN, ctx.ClientIP())
		info.Module.CreateTime = curTimestampMs
		NewAsyncTaskWithParam(ASYNC_DEV_AND_MODULE_CREATE, info)
	case MODULE_BUILD_BIT: //需单独创建Device
		info.Device.Create(data.AccessWay, data.DeviceSN, data.DeviceVersion, DEVICE_ONLINE)
		info.Log.Create(0, data.AccessWay, data.ModuleSN, ctx.ClientIP())
		info.Device.CreateTime = curTimestampMs
		NewAsyncTaskWithParam(ASYNC_DEV_AND_MODULE_CREATE, info)
	case HAS_DEVICE_WITH_MODULE:
		info.Module.Update(data.ModuleVersion)
		moduleUpdateMap := map[string]interface{}{"module_version": info.Module.ModuleVersion, "update_time": info.Module.UpdateTime}
		moduleTask := NewTask()
		moduleTask.Param = moduleUpdateMap
		moduleTask.RunTaskWithTypeAndEntity(ASYNC_UP_MODULE_INFO, info.Module)

		info.Device.Update(info.Module.DeviceID, data.AccessWay, data.DeviceVersion, info.Module.UpdateTime, DEVICE_ONLINE)
		info.Device.DeviceSn = data.DeviceSN
		deviceUpdateMap := map[string]interface{}{"access_way": info.Device.AccessWay, "device_version": info.Device.DeviceVersion, "update_time": info.Device.UpdateTime, "status": info.Device.Status}
		deviceTask := NewTask()
		deviceTask.Param = deviceUpdateMap
		deviceTask.RunTaskWithTypeAndEntity(ASYNC_UP_DEVICE_INFO, info.Device)
		// 检测并返回固件版本
		// 返回版本升级格式
		//	data := &FirmwareInfo{
		//		URL:  "http://www.gisunlink.com/GiSunLink.ota.bin",
		//		Size: 476448,
		//	}
		createConnectLog(ctx, info.Module.ID, info.Module.UID, data.AccessWay, data.ModuleSN)
	}
	return CreateMessage(SUCCESS, respond)
}
