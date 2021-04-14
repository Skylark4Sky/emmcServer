package device

import (
	. "GoServer/middleWare/dataBases/redis"
	mqtt "GoServer/mqttPacket"
	. "GoServer/utils/float64"
	. "GoServer/utils/log"
	. "GoServer/utils/time"
	"time"
)

const (
	REDIS_LOCK_STATETIMEOUT = 5
)

//构建缓存数据实体
func calculateComData(comData *mqtt.ComData) *CacheComData {
	cacheComData := &CacheComData{
		ComData:              *comData,
		MaxChargeElectricity: comData.CurElectricity,
		CurPower:             0,
		AveragePower:         0,
		MaxPower:             0,
		SyncTime:             0,
		WriteFlags:           0,
	}
	return cacheComData
}

func deviceComStateDataCHK(newData, redisData *CacheComData, userID, deviceID uint64) bool {
	//如果token一致需检测数据变更
	if newData.Token == redisData.Token && (newData.Enable == COM_DISENABLE && redisData.Enable == COM_ENABLE) {
		comChargeTaskNeedCheck(redisData, userID, deviceID)
		SystemLog("实时: ", newData.Enable, " 缓存: ", redisData.Enable)
		SystemLog("deviceID", deviceID)
		SystemLog("实时: ---> ", newData)
		SystemLog("缓存: ---> ", redisData)
		return true
	}
	return false
}

func deviceComCharge(newComData, redisCacheData *CacheComData, userID, deviceID uint64) bool {
	//端口启用状态
	if newComData.Enable == COM_ENABLE {
		// 拷贝数据
		newComData.MaxPower = redisCacheData.MaxPower
		newComData.MaxChargeElectricity = redisCacheData.MaxChargeElectricity
		newComData.SyncTime = redisCacheData.SyncTime
		// 数值计算
		newComData.CurPower = CalculateCurComPower(CUR_VOLTAGE, newComData.CurElectricity, 5)              //当前P值
		newComData.AveragePower = CalculateCurAverageComPower(newComData.UseEnergy, newComData.UseTime, 5) //平均P值
		// 最大P值
		if CmpPower(newComData.CurPower, redisCacheData.MaxPower) > 0 {
			newComData.MaxPower = newComData.CurPower
		}
		// 最大E记录
		if newComData.CurElectricity > redisCacheData.MaxChargeElectricity {
			newComData.MaxChargeElectricity = newComData.CurElectricity
		}
		// 检测同步时间是否可以同步当前端口数据入库
		curTime := GetTimestampMs()
		if (redisCacheData.SyncTime == 0) || (curTime-redisCacheData.SyncTime) >= COM_DATA_SYNC_TIME {
			newComData.SyncTime = curTime
			comChargeTaskDataUpdate(newComData, userID, deviceID)
		}
		return true
	}

	//未开启时，检测值是否有变化 大于50ma保存异常数值
	if CmpElectricityException(float64(newComData.CurElectricity), float64(redisCacheData.CurElectricity), 50) {
		SystemLog(" Time:", TimeFormat(time.Now()), " ", deviceID, " 端口:", newComData.Id, " 异常---当前值:", newComData.CurElectricity, "上一次值为:", redisCacheData.CurElectricity)
	}

	return false
}

// 运行状态 & 空闲状态 处理
func deviceStateHandle(comList *mqtt.ComList, deviceSN string, userID, deviceID uint64) {
	redisComList := BatchReadDeviceComDataiFromRedis(deviceSN)

	//批量写入数据到缓存
	BatchWriteDeviceComDataToRedis(deviceSN, comList, func(comData *mqtt.ComData) *CacheComData {
		copyNewData := true
		newComData := calculateComData(comData)
		if redisCacheData, ok := redisComList[comData.Id]; ok {
			if deviceComStateDataCHK(newComData, &redisCacheData, userID, deviceID) {
				redisCacheData = *newComData
			}
			if !deviceComCharge(newComData, &redisCacheData, userID, deviceID) {
				copyNewData = false
				delete(redisComList, comData.Id)
			}
		}

		if copyNewData {
			redisComList[comData.Id] = *newComData
		}
		return newComData
	})

	//刷新设备状态 并 统计 数据值
	refreshDeviceStatus(deviceSN, deviceID, &DeviceStatus{
		Behavior:     comList.ComBehavior,
		Signal:       int8(comList.Signal),
		Worker:       Redis().TatolWorkerByDevice(deviceSN, redisComList),
		ID:           deviceID,
		ProtoVersion: comList.ComProtoVer,
	})
}

//上报数据处理
func deviceActBehaviorDataOps(packet *mqtt.Packet, deviceSN string, userID, deviceID uint64) {
	switch packet.Json.Behavior {
	//运行状态值运算
	case mqtt.GISUNLINK_CHARGEING, mqtt.GISUNLINK_CHARGE_LEISURE:
		deviceStateHandle(packet.Data.(*mqtt.ComList), deviceSN, userID, deviceID)
	case mqtt.GISUNLINK_CHARGE_TASK:
		comChargeTaskStart(packet.Data.(*mqtt.ComTaskStartTransfer), deviceSN, userID, deviceID, false)
	case mqtt.GISUNLINK_START_CHARGE:
		comChargeTaskStart(packet.Data.(*mqtt.ComList), deviceSN, userID, deviceID, true)
	case mqtt.GISUNLINK_EXIT_CHARGE_TASK:
		comChargeTaskStop(packet.Data.(*mqtt.ComTaskStopTransfer), deviceSN, userID, deviceID, false)
	case mqtt.GISUNLINK_STOP_CHARGE:
		comChargeTaskStop(packet.Data.(*mqtt.ComList), deviceSN, userID, deviceID, true)
	case mqtt.GISUNLINK_CHARGE_NO_LOAD, mqtt.GISUNLINK_CHARGE_FINISH, mqtt.GISUNLINK_CHARGE_BREAKDOWN:
		deviceInitiativeExitComChargeTask(packet.Data.(*mqtt.ComList), deviceSN, userID, deviceID, packet.Json.Behavior)
	}
}
