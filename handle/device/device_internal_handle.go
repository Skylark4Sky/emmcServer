package device

import (
	. "GoServer/middleWare/dataBases/mysql"
	. "GoServer/middleWare/dataBases/redis"
	. "GoServer/model"
	. "GoServer/model/device"
	mqtt "GoServer/mqttPacket"
	. "GoServer/utils/log"
	. "GoServer/utils/string"
	. "GoServer/utils/time"
	"encoding/json"
)

const (
	SYNC_UPDATE_TIME = 300000 //5*60*1000 5分钟同步一次
)

const (
	UPDATE_DEVICE_STATUS            = 0x02
	UPDATE_DEVICE_WORKER            = 0x04
	UPDATE_DEVICE_STATUS_AND_WORKER = 0x06
)

func changeDeviceStatus(device_sn string, deviceID uint64, updateFlags uint8, status int8, worker int8) {
	if deviceID != 0 {
		device := &DeviceInfo{}
		device.ID = deviceID

		var deviceUpdateMap = make(map[string]interface{})

		deviceUpdateMap["update_time"] = GetTimestampMs()

		if (updateFlags & UPDATE_DEVICE_STATUS) == UPDATE_DEVICE_STATUS {
			deviceUpdateMap["status"] = status
			Redis().SetDeviceStatusToRedis(device_sn, int(status))
		}

		if (updateFlags & UPDATE_DEVICE_WORKER) == UPDATE_DEVICE_WORKER {
			deviceUpdateMap["worker"] = worker
			Redis().SetDeviceWorkerToRedis(device_sn, int(worker))
		}

		if updateFlags != 0 {
			CreateAsyncSQLTaskWithUpdateMap(ASYNC_UP_DEVICE_INFO, device, deviceUpdateMap)
		}
	}
}

//状态刷新
func refreshDeviceStatus(deviceSN string, deviceID uint64, deviceStatus *DeviceStatus) {
	var timeout int64 = CHARGEING_TIME
	if deviceStatus.Behavior == mqtt.GISUNLINK_CHARGE_LEISURE {
		timeout = LEISURE_TIME
	}

	status := Redis().GetDeviceStatusFromRedis(deviceSN)
	worker := Redis().GetDeviceWorkerFormRedis(deviceSN)
	syncTime := Redis().GetDeviceSyncTimeFromRedis(deviceSN)

	var updateFlags uint8 = 0
	var workerStatus int8 = DEVICE_WORKING

	switch deviceStatus.Behavior {
	case mqtt.GISUNLINK_CHARGEING:
		{
			if status != int(DEVICE_WORKING) {
				updateFlags |= UPDATE_DEVICE_STATUS
				workerStatus = DEVICE_WORKING
			}
		}
	case mqtt.GISUNLINK_CHARGE_LEISURE:
		{
			if status != int(DEVICE_ONLINE) {
				updateFlags |= UPDATE_DEVICE_STATUS
				workerStatus = DEVICE_ONLINE
			}
		}
	}

	if uint8(worker) != deviceStatus.Worker {
		updateFlags |= UPDATE_DEVICE_WORKER
	}

	curTimeMS := GetTimestampMs()
	if syncTime <= 0 || (curTimeMS-syncTime) >= SYNC_UPDATE_TIME {
		Redis().SetDeviceSyncTimeToRedis(deviceSN, curTimeMS)
	}

	//更新状态
	changeDeviceStatus(deviceSN, deviceID, updateFlags, workerStatus, int8(deviceStatus.Worker))
	//更新令牌时间
	Redis().UpdateDeviceTokenExpiredTime(deviceSN, deviceStatus, timeout)
}

//保存设备上报状态
func createDeviceTransferLog(transfer *DeviceTransferLog) {
	if transfer == nil {
		return
	}
	transfer.CreateTime = GetTimestampMs()
	if err := ExecSQL().Create(&transfer).Error; err != nil {
		SystemLog("CreateDeviceTransferLog", err)
	}
}

//保存上报数据入库
func saveDeviceTransferDataOps(serverNode string, device_sn string, deviceID uint64, packet *mqtt.Packet) {
	var comNum uint8 = 0
	var payloadData string = ""
	switch packet.Json.Behavior {
	case mqtt.GISUNLINK_CHARGEING, mqtt.GISUNLINK_CHARGE_LEISURE: //运行中,空闲中
		comList := packet.Data.(*mqtt.ComList)
		comNum = comList.ComNum
		comPort := comList.ComPort
		if jsonString, err := json.Marshal(comPort); err == nil {
			payloadData = string(jsonString)
		}
		break
	case mqtt.GISUNLINK_START_CHARGE, mqtt.GISUNLINK_STOP_CHARGE, mqtt.GISUNLINK_CHARGE_FINISH, mqtt.GISUNLINK_CHARGE_NO_LOAD, mqtt.GISUNLINK_CHARGE_BREAKDOWN, mqtt.GISUNLINK_UPDATE_FIRMWARE, mqtt.GISUNLINK_COM_UPDATE, mqtt.GISUNLINK_COM_NO_UPDATE: //开始,停止,完成,空载,故障,升级,参数刷新,没有刷新参数
		comList := packet.Data.(*mqtt.ComList)
		comNum = comList.ComNum
		comPort := comList.ComPort
		if jsonString, err := json.Marshal(comPort); err == nil {
			payloadData = string(jsonString)
		}
		break
	case mqtt.GISUNLINK_CHARGE_TASK:
		comTaskStartTransfer := packet.Data.(*mqtt.ComTaskStartTransfer)
		comNum = comTaskStartTransfer.ComID
		if jsonString, err := json.Marshal(comTaskStartTransfer); err == nil {
			payloadData = string(jsonString)
		}
		break
	case mqtt.GISUNLINK_DEVIDE_STATUS:
		comTaskStatusQueryTransfer := packet.Data.(*mqtt.ComTaskStatusQueryTransfer)
		comNum = comTaskStatusQueryTransfer.ComID
		if jsonString, err := json.Marshal(comTaskStatusQueryTransfer); err == nil {
			payloadData = string(jsonString)
		}
		break
	case mqtt.GISUNLINK_EXIT_CHARGE_TASK:
		comTaskStopTransfer := packet.Data.(*mqtt.ComTaskStopTransfer)
		comNum = comTaskStopTransfer.ComID
		if jsonString, err := json.Marshal(comTaskStopTransfer); err == nil {
			payloadData = string(jsonString)
		}
		break
	case mqtt.GISUNLINK_SET_CONFIG:
		deviceSetConfigTransfer := packet.Data.(*mqtt.DeviceSetConfigTransfer)
		comNum = 0xFF
		if jsonString, err := json.Marshal(deviceSetConfigTransfer); err == nil {
			payloadData = string(jsonString)
		}
		break
	}

	log := &DeviceTransferLog{
		TransferID:   int64(packet.Json.ID),
		DeviceID:     deviceID,
		TransferAct:  packet.Json.Act,
		DeviceSn:     device_sn,
		ComNum:       comNum,
		TransferData: packet.Json.Data,
		PayloadData:  payloadData,
		Behavior:     packet.Json.Behavior,
		ServerNode:   serverNode,
		TransferTime: int64(packet.Json.Ctime),
	}

	createDeviceTransferLog(log)
}

func deviceExpiredMsgOps(pattern, channel, message string) {
	deviceSN := GetDeviceSN(message, ":")
	if deviceID := Redis().GetDeviceIDFromRedis(deviceSN); deviceID != 0 {
		switch message {
		case GetDeviceTokenKey(deviceSN): //这里处理过期key
			{
				SystemLog("deviceExpiredMsgOps: ", deviceSN, " DEVICE_OFFLINE")
				changeDeviceStatus(deviceSN, deviceID, UPDATE_DEVICE_STATUS, DEVICE_OFFLINE, 0)
			}
		case GetComdDataKey(deviceSN), GetDeviceInfoKey(deviceSN):
			{
				SystemLog("deviceExpiredMsgOps: ", message)
			}
		}
	}
}

