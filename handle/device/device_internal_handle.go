package device

import (
	. "GoServer/middleWare/dataBases/mysql"
	. "GoServer/middleWare/dataBases/redis"
	. "GoServer/model"
	. "GoServer/model/device"
	mqtt "GoServer/mqttPacket"
	. "GoServer/utils/float64"
	. "GoServer/utils/log"
	. "GoServer/utils/string"
	. "GoServer/utils/time"
	"encoding/json"
	"time"
)

const (
	UPDATE_DEVICE_STATUS            = 0x02
	UPDATE_DEVICE_WORKER            = 0x04
	UPDATE_DEVICE_STATUS_AND_WORKER = 0x06
)

func changeDeviceStatus(device_sn string, deviceID uint64, updateFlags uint8, status int8, worker int8) {
	if deviceID != 0 {
		device := &DeviceInfo{}
		device.ID = Redis().GetDeviceIDFromRedis(device_sn)

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
			SystemLog("deviceID: ", deviceID, " deviceUpdateMap: ", deviceUpdateMap)
			CreateAsyncSQLTaskWithUpdateMap(ASYNC_UP_DEVICE_INFO, device, deviceUpdateMap)
		}
	}
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
	case mqtt.GISUNLINK_START_CHARGE, mqtt.GISUNLINK_STOP_CHARGE,mqtt.GISUNLINK_CHARGE_FINISH, mqtt.GISUNLINK_CHARGE_NO_LOAD, mqtt.GISUNLINK_CHARGE_BREAKDOWN ,mqtt.GISUNLINK_UPDATE_FIRMWARE,mqtt.GISUNLINK_COM_UPDATE,mqtt.GISUNLINK_COM_NO_UPDATE: //开始,停止,完成,空载,故障,升级,参数刷新,没有刷新参数
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
				deviceID := Redis().GetDeviceIDFromRedis(deviceSN)
				changeDeviceStatus(deviceSN, deviceID, UPDATE_DEVICE_STATUS, DEVICE_OFFLINE, 0)
			}
		case GetComdDataKey(deviceSN), GetDeviceInfoKey(deviceSN):
			{
				SystemLog("deviceExpiredMsgOps: ", message)
			}
		}
	}
}

//比较数据
func analyseComData(tokenKey string, newData *mqtt.ComList, cacheData map[uint8]ComDataTotal) {
	//不存在缓存数据直接返回
	if len(cacheData) <= 0 {
		return
	}

	for _, comID := range newData.ComID {
		var index uint8 = comID
		if newData.ComNum <= 5 {
			index = (comID % 5)
		}
		comData := (newData.ComPort[int(index)]).(mqtt.ComData)
		comData.Id = comID
		cacherData := cacheData[comData.Id]
		//未开启时，检测值是否有变化
		if comData.Enable == 0 {
			//大于50ma保存异常数值
			if (comData.CurElectricity >= 50) && (comData.CurElectricity != cacherData.CurElectricity) {
				SystemLog(" Time:", TimeFormat(time.Now()), " ", tokenKey, " 端口:", comData.Id, " 异常---当前值:", comData.CurElectricity, "上一次值为:", cacherData.CurElectricity)
			}
		} else {

			//comData.MaxPower
		}
	}
}

func calculateComData(comData *mqtt.ComData) *ComDataTotal {
	dataTotal := &ComDataTotal{
		ComData:      *comData,
		CurPower:     0,
		AveragePower: 0,
		MaxPower:     0,
	}
	dataTotal.CurPower = CalculateCurComPower(CUR_VOLTAGE, dataTotal.CurElectricity, 2)
	dataTotal.AveragePower = CalculateCurAverageComPower(dataTotal.UseEnergy, dataTotal.UseTime, 2)
	return dataTotal

}

func deviceActBehaviorDataOps(packet *mqtt.Packet, cacheKey string, deviceID uint64, playload string) {
	switch packet.Json.Behavior {
	case mqtt.GISUNLINK_CHARGEING, mqtt.GISUNLINK_CHARGE_LEISURE:
		{
			comList := packet.Data.(*mqtt.ComList)
			//批量读当前所有接口
			cacherComData := BatchReadDeviceComDataiFromRedis(cacheKey)
			//对比新旧数据
			analyseComData(cacheKey, comList, cacherComData)

			//批量写入数据
			if len(cacherComData) > 1 {
				BatchWriteDeviceComDataToRedis(cacheKey, comList, func(comData *mqtt.ComData) *ComDataTotal {
					dataTotal := calculateComData(comData)
					cacherData := cacherComData[comData.Id]
					dataTotal.MaxPower = cacherData.MaxPower
					if CmpPower(dataTotal.CurPower, cacherData.MaxPower) == 1 {
						dataTotal.MaxPower = dataTotal.CurPower
					}
					return dataTotal

				})
			} else {
				BatchWriteDeviceComDataToRedis(cacheKey, comList, func(comData *mqtt.ComData) *ComDataTotal {
					return calculateComData(comData)
				})
			}

			deviceStatus := &DeviceStatus{
				Behavior:     comList.ComBehavior,
				Signal:       int8(comList.Signal),
				Worker:       Redis().TatolWorkerByDevice(cacheKey, BatchReadDeviceComDataiFromRedis(cacheKey)),
				ID:           deviceID,
				ProtoVersion: comList.ComProtoVer,
			}

			var timeout int64 = CHARGEING_TIME
			if packet.Json.Behavior == mqtt.GISUNLINK_CHARGE_LEISURE {
				timeout = LEISURE_TIME
			}

			status := Redis().GetDeviceStatusFromRedis(cacheKey)
			worker := Redis().GetDeviceWorkerFormRedis(cacheKey)

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

			//更新状态
			changeDeviceStatus(cacheKey, deviceID, updateFlags, workerStatus, int8(deviceStatus.Worker))

			//更新令牌时间
			Redis().UpdateDeviceTokenExpiredTime(cacheKey, deviceStatus, timeout)
		}
	}
}
