package device

import (
	. "GoServer/middleWare/dataBases/mysql"
	. "GoServer/model"
	. "GoServer/model/device"
	. "GoServer/utils/log"
	. "GoServer/utils/respond"
	. "GoServer/utils/threadWorker"
	. "GoServer/utils/time"
	"github.com/gin-gonic/gin"
)

type RequestParam struct {
	ClientID string `form:"clientID" json:"type" binding:"required"`
}

type RequestData struct {
	AccessWay AccesswayType `form:"type" json:"type" binding:"required"`
	DeviceSN  string        `form:"deviceNo" json:"deviceNo" binding:"required"`
	ModuleSN  string        `form:"token" json:"token" binding:"required"`
	Version   string        `form:"version" json:"version" binding:"required"`
}

type FirmwareInfo struct {
	Size int64  `json:"size"`
	URL  string `json:"url"`
}

func createConnectLog(ctx *gin.Context, device_id int64, accessway AccesswayType, moduleSN string) {
	var task AsynSQLTask
	var log DeviceConnectLog

	log.Create(device_id, accessway, moduleSN, ctx.ClientIP())
	task.Entity = log
	var work Job = &task
	InsertAsynTask(work)
}

func (data *RequestData) Connect(ctx *gin.Context, clientID string) interface{} {

	var device DeviceInfo

	err := ExecSQL().Where("module_sn = ?", clientID).First(&device).Error

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
		device.Create(data.AccessWay, 0, clientID, data.DeviceSN, data.Version)
		if err := ExecSQL().Create(&device).Error; err != nil {
			return CreateMessage(SYSTEM_ERROR, err)
		}
	}

	createConnectLog(ctx, device.ID, data.AccessWay, clientID)

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
