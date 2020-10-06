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
	var log DeviceConnectLog
	log.Create(device_id, accessway, moduleSN, ctx.ClientIP())
	CreateAsyncSQLTask(ASYNC_DEV_CONNECT_LOG, log)
}

func (data *RequestData) Connect(ctx *gin.Context) interface{} {

	var device DeviceInfo

	err := ExecSQL().Where("module_sn = ?", data.ModuleSN).First(&device).Error

	var hasRecord = true

	if err != nil {
		if IsRecordNotFound(err) {
			hasRecord = false
		} else {
			return CreateMessage(SYSTEM_ERROR, err)
		}
	}

	if !hasRecord {
		// 建立新记录
		device.Create(data.AccessWay, 0, data.ModuleSN, data.DeviceSN, data.ModuleVersion, data.DeviceVersion)
		if err := ExecSQL().Create(&device).Error; err != nil {
			return CreateMessage(SYSTEM_ERROR, err)
		}
	} else {
		device.Update(data.ModuleVersion, data.DeviceVersion)
		updateMap := map[string]interface{}{"module_version": data.ModuleVersion, "device_version": data.DeviceVersion, "update_time": device.UpdateTime}
		if err := ExecSQL().Model(&device).Updates(updateMap).Error; err != nil {
			return CreateErrorMessage(SYSTEM_ERROR, err)
		}
	}

	createConnectLog(ctx, device.ID, data.AccessWay, data.ModuleSN)

	//	data := &FirmwareInfo{
	//		URL:  "http://www.gisunlink.com/GiSunLink.ota.bin",
	//		Size: 476448,
	//	}

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
