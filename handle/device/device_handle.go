package device

import (
	. "GoServer/middleWare/dataBases/mysql"
	. "GoServer/model"
	. "GoServer/model/device"
	. "GoServer/utils/log"
	. "GoServer/utils/respond"
	. "GoServer/utils/time"
	"github.com/gin-gonic/gin"
)

type RequestData struct {
	AccessWay     AccesswayType `form:"type" json:"type" binding:"required"`
	ModuleSN      string        `form:"module_sn" json:"module_sn" binding:"required"`
	ModuleVersion string        `form:"module_version" json:"module_version" binding:"required"`
	DeviceSN      string        `form:"device_sn" json:"device_sn" binding:"required"`
	DeviceVersion string        `form:"device_version" json:"device_version" binding:"required"`
	Token         string        `form:"token" json:"token" binding:"required"`
}

type FirmwareInfo struct {
	Size int64  `json:"size"`
	URL  string `json:"url"`
}

func createConnectLog(ctx *gin.Context, device_id int64, accessway AccesswayType, moduleSN string) {
	var log ModuleConnectLog
	log.Create(device_id, accessway, moduleSN, ctx.ClientIP())
	CreateAsyncSQLTask(ASYNC_MODULE_CONNECT_LOG, log)
}

func (data *RequestData) Connect(ctx *gin.Context) interface{} {
	var module ModuleInfo
	err := ExecSQL().Where("module_sn = ?", data.ModuleSN).First(&module).Error
	var hasRecord = true

	if err != nil {
		if IsRecordNotFound(err) {
			hasRecord = false
		} else {
			return CreateMessage(SYSTEM_ERROR, err)
		}
	}

	if !hasRecord {
		//创建对应关系
		var info CreateDeviceInfo
		curTimestampMs := GetTimestampMs()
		info.Module.Create(data.AccessWay, data.ModuleSN, data.ModuleVersion)
		info.Device.Create(data.DeviceSN, data.DeviceVersion)
		info.Log.Create(0, data.AccessWay, data.ModuleSN, ctx.ClientIP())
		info.Module.CreateTime = curTimestampMs
		info.Device.CreateTime = curTimestampMs
		CreateAsyncSQLTask(ASYNC_DEV_AND_MODULE_CREATE, info)
		return CreateMessage(SUCCESS, nil)
	} else {
		module.Update(data.ModuleVersion)
		CreateAsyncSQLTask(ASYNC_UP_MODULE_VERSION, module)
		var device DeviceInfo
		device.Update(module.DeviceID,data.DeviceVersion,module.UpdateTime)
		CreateAsyncSQLTask(ASYNC_UP_DEVICE_VERSION, device)
		// 检测并返回固件版本
		// 返回版本升级格式
		//	data := &FirmwareInfo{
		//		URL:  "http://www.gisunlink.com/GiSunLink.ota.bin",
		//		Size: 476448,
		//	}
		createConnectLog(ctx, module.ID, data.AccessWay, data.ModuleSN)
		return CreateMessage(SUCCESS, nil)
	}
	return CreateMessage(SUCCESS, nil)
}

func CreateDeviceTransferLog(transfer *DeviceTransferLog) {
	if transfer == nil {
		return
	}

	transfer.CreateTime = GetTimestampMs()
	if err := ExecSQL().Create(&transfer).Error; err != nil {
		SystemLog("CreateDeviceTransferLog", err)
	}
}
