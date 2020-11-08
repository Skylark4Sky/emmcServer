package device

import (
	. "GoServer/middleWare/dataBases/mysql"
	. "GoServer/model"
	. "GoServer/model/device"
	. "GoServer/utils/respond"
	. "GoServer/utils/time"
	"github.com/gin-gonic/gin"
)

const (
	CHARGEING_TIME = 60 //忙时150秒
	LEISURE_TIME   = 120 //空闲时300秒
)

const (
	DEVICE_OFFLINE 	int8 = 0
	DEVICE_ONLINE 	int8 = 1
	DEVICE_WORKING 	int8 = 2
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

func createConnectLog(ctx *gin.Context, device_id uint64, accessway uint8, moduleSN string) {
	log := &ModuleConnectLog{}
	log.Create(device_id, accessway, moduleSN, ctx.ClientIP())
	CreateAsyncSQLTask(ASYNC_MODULE_CONNECT_LOG, log)
}

func (data *RequestData) Connect(ctx *gin.Context) interface{} {
	module := &ModuleInfo{}
	err := ExecSQL().Where("module_sn = ?", data.ModuleSN).First(module).Error
	var hasRecord = true

	if err != nil {
		if IsRecordNotFound(err) {
			hasRecord = false
		} else {
			return CreateMessage(SYSTEM_ERROR, err)
		}
	}

	var respond interface{} = nil

	if !hasRecord {
		//创建对应关系 第一次连接仅建立关系，版本检测下次连接时在进行处理
		var info CreateDeviceInfo
		curTimestampMs := GetTimestampMs()
		info.Module.Create(data.AccessWay, data.ModuleSN, data.ModuleVersion)
		info.Device.Create(data.AccessWay, data.DeviceSN, data.DeviceVersion, DEVICE_ONLINE)
		info.Log.Create(0, data.AccessWay, data.ModuleSN, ctx.ClientIP())
		info.Module.CreateTime = curTimestampMs
		info.Device.CreateTime = curTimestampMs
		CreateAsyncSQLTask(ASYNC_DEV_AND_MODULE_CREATE, info)
	} else {
		module.Update(data.ModuleVersion)
		moduleUpdateMap := map[string]interface{}{"module_version": module.ModuleVersion, "update_time": module.UpdateTime}
		CreateAsyncSQLTaskWithUpdateMap(ASYNC_UP_MODULE_INFO, module, moduleUpdateMap)
		device := &DeviceInfo{}
		device.Update(module.DeviceID, data.AccessWay, data.DeviceVersion, module.UpdateTime, DEVICE_ONLINE)
		device.DeviceSn = data.DeviceSN
		deviceUpdateMap := map[string]interface{}{"access_way": device.AccessWay, "device_version": device.DeviceVersion, "update_time": device.UpdateTime, "status": device.Status}
		CreateAsyncSQLTaskWithUpdateMap(ASYNC_UP_DEVICE_INFO, device, deviceUpdateMap)
		// 检测并返回固件版本
		// 返回版本升级格式
		//	data := &FirmwareInfo{
		//		URL:  "http://www.gisunlink.com/GiSunLink.ota.bin",
		//		Size: 476448,
		//	}
		createConnectLog(ctx, module.ID, data.AccessWay, data.ModuleSN)
	}
	return CreateMessage(SUCCESS, respond)
}
