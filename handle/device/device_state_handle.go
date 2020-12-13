package device

import (
	. "GoServer/middleWare/dataBases/redis"
	mqtt "GoServer/mqttPacket"
	. "GoServer/utils/float64"
	. "GoServer/utils/log"
	. "GoServer/utils/time"
	"time"
)

//构建缓存数据实体
func calculateComData(comData *mqtt.ComData) *ComDataTotal {
	dataTotal := &ComDataTotal{
		ComData:              *comData,
		MaxChargeElectricity: comData.CurElectricity,
		CurPower:             0,
		AveragePower:         0,
		MaxPower:             0,
	}
	return dataTotal
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
		if comData.Enable == COM_DISENABLE {
			//大于50ma保存异常数值
			if (comData.CurElectricity >= 50) && (comData.CurElectricity != cacherData.CurElectricity) {
				SystemLog(" Time:", TimeFormat(time.Now()), " ", tokenKey, " 端口:", comData.Id, " 异常---当前值:", comData.CurElectricity, "上一次值为:", cacherData.CurElectricity)
			}
		} else {
			//comData.MaxPower
		}
		if comData.Token != cacherData.Token {
			SystemLog(" NewToKen:", comData.Token, " OldToKen:", cacherData.Token)
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

// 运行状态 & 空闲状态 处理
func deviceStateHandle(comList *mqtt.ComList, cacheKey string, deviceID uint64) {
	//批量读当前设备所有接口
	cacherComData := BatchReadDeviceComDataiFromRedis(cacheKey)

	//对比新旧数据
	analyseComData(cacheKey, comList, cacherComData)

	//批量写入数据
	if len(cacherComData) > 1 {
		BatchWriteDeviceComDataToRedis(cacheKey, comList, func(comData *mqtt.ComData) *ComDataTotal {
			dataTotal := calculateComData(comData)
			cacherData := cacherComData[comData.Id]
			//端口启用状态
			if comData.Enable == COM_ENABLE {
				// 数值计算
				dataTotal.CurPower = CalculateCurComPower(CUR_VOLTAGE, dataTotal.CurElectricity, 2)             //当前功率
				dataTotal.AveragePower = CalculateCurAverageComPower(dataTotal.UseEnergy, dataTotal.UseTime, 2) //平均功率
				//最大功率
				if CmpPower(dataTotal.CurPower, cacherData.MaxPower) >= 1 {
					dataTotal.MaxPower = dataTotal.CurPower
				}
				// 最大电流记录
				if comData.CurElectricity > cacherData.MaxChargeElectricity {
					dataTotal.MaxChargeElectricity = comData.CurElectricity
				}
			}
			return dataTotal
		})
	} else {
		BatchWriteDeviceComDataToRedis(cacheKey, comList, func(comData *mqtt.ComData) *ComDataTotal {
			return calculateComData(comData)
		})
	}

	//刷新设备状态
	refreshDeviceStatus(cacheKey, deviceID, &DeviceStatus{
		Behavior:     comList.ComBehavior,
		Signal:       int8(comList.Signal),
		Worker:       Redis().TatolWorkerByDevice(cacheKey, BatchReadDeviceComDataiFromRedis(cacheKey)),
		ID:           deviceID,
		ProtoVersion: comList.ComProtoVer,
	})
}

//上报数据处理
func deviceActBehaviorDataOps(packet *mqtt.Packet, deviceSN string, deviceID uint64) {
	switch packet.Json.Behavior {
	//运行状态值运算
	case mqtt.GISUNLINK_CHARGEING,mqtt.GISUNLINK_CHARGE_LEISURE:
		{
			deviceStateHandle(packet.Data.(*mqtt.ComList), deviceSN, deviceID)
		}
	//开始记录该项订单记录
	case mqtt.GISUNLINK_START_CHARGE:
		{
			SystemLog("CMD: GISUNLINK_START_CHARGE ", " deviceSN: ", deviceSN, " deviceID: ", deviceID)
		}
	//以下值均会触发停止订单行为
	case mqtt.GISUNLINK_CHARGE_NO_LOAD,mqtt.GISUNLINK_CHARGE_FINISH,mqtt.GISUNLINK_STOP_CHARGE,mqtt.GISUNLINK_CHARGE_BREAKDOWN:
		{
			SystemLog("CMD: GISUNLINK_CHARGE_FINISH ", " deviceSN: ", deviceSN, " deviceID: ", deviceID)
		}
	}
}