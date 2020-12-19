package device

import (
	. "GoServer/middleWare/dataBases/redis"
	mqtt "GoServer/mqttPacket"
	. "GoServer/utils/float64"
	. "GoServer/utils/log"
	. "GoServer/utils/time"
	"math"
	"time"
)

const (
	COM_DATA_SYNC_TIME = 3000000
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

//数据检查
func comListDataChk(deviceSN string, comList *mqtt.ComList, redisCacheComList map[uint8]CacheComData) {
	//不存在数据直接返回
	if comList != nil {
		return
	}

	for _, comID := range comList.ComID {
		var index uint8 = comID
		if comList.ComNum <= 5 {
			index = (comID % 5)
		}
		comData := (comList.ComPort[int(index)]).(mqtt.ComData)
		comData.Id = comID
		if redisComData, ok := redisCacheComList[comData.Id]; ok {
			//未开启时，检测值是否有变化
			if comData.Enable == COM_DISENABLE {
				//大于50ma保存异常数值
				if (comData.CurElectricity != redisComData.CurElectricity) && (math.Abs(float64(comData.CurElectricity)-float64(redisComData.CurElectricity)) >= 50) {
					SystemLog(" Time:", TimeFormat(time.Now()), " ", deviceSN, " 端口:", comData.Id, " - ", redisComData.Id, " 异常---当前值:", comData.CurElectricity, "上一次值为:", redisComData.CurElectricity)
				}
			}
		}
	}
}

// 运行状态 & 空闲状态 处理
func deviceStateHandle(comList *mqtt.ComList, deviceSN string, deviceID uint64) {
	redisComList := BatchReadDeviceComDataiFromRedis(deviceSN)
	//数据分析
	comListDataChk(deviceSN, comList, redisComList)
	//批量写入数据到缓存 并 统计 功率
	if len(redisComList) > 1 {
		BatchWriteDeviceComDataToRedis(deviceSN, comList, func(comData *mqtt.ComData) *CacheComData {
			newComData := calculateComData(comData)
			redisCacheData := redisComList[comData.Id]
			//端口启用状态
			if comData.Enable == COM_ENABLE {
				// 数值计算
				newComData.CurPower = CalculateCurComPower(CUR_VOLTAGE, newComData.CurElectricity, 5)              //当前P值
				newComData.AveragePower = CalculateCurAverageComPower(newComData.UseEnergy, newComData.UseTime, 5) //平均P值
				// 最大P值
				newComData.MaxPower = redisCacheData.MaxPower
				if CmpPower(newComData.CurPower, redisCacheData.MaxPower) > 0 {
					newComData.MaxPower = newComData.CurPower
				}
				// 最大E记录
				newComData.MaxChargeElectricity = redisCacheData.MaxChargeElectricity
				if comData.CurElectricity > redisCacheData.MaxChargeElectricity {
					newComData.MaxChargeElectricity = comData.CurElectricity
				}
				// 检测同步时间是否可以同步当前端口数据入库
				curTime := GetTimestampMs()
				newComData.SyncTime = redisCacheData.SyncTime
				if (redisCacheData.SyncTime == 0) || (curTime-redisCacheData.SyncTime) >= COM_DATA_SYNC_TIME {
					newComData.SyncTime = curTime
					comChargeTaskDataUpdate(newComData, deviceID)
				}

				redisComList[comData.Id] = *newComData
			}
			return newComData
		})
	} else {
		BatchWriteDeviceComDataToRedis(deviceSN, comList, func(comData *mqtt.ComData) *CacheComData {
			cacheComData := calculateComData(comData)
			redisComList[comData.Id] = *cacheComData
			return cacheComData
		})
	}

	//刷新设备状态
	refreshDeviceStatus(deviceSN, deviceID, &DeviceStatus{
		Behavior:     comList.ComBehavior,
		Signal:       int8(comList.Signal),
		Worker:       Redis().TatolWorkerByDevice(deviceSN, redisComList),
		ID:           deviceID,
		ProtoVersion: comList.ComProtoVer,
	})
}

//上报数据处理
func deviceActBehaviorDataOps(packet *mqtt.Packet, deviceSN string, deviceID uint64) {
	switch packet.Json.Behavior {
	//运行状态值运算
	case mqtt.GISUNLINK_CHARGEING, mqtt.GISUNLINK_CHARGE_LEISURE:
		deviceStateHandle(packet.Data.(*mqtt.ComList), deviceSN, deviceID)
	case mqtt.GISUNLINK_CHARGE_TASK:
		comChargeTaskStart(packet.Data.(*mqtt.ComTaskStartTransfer), deviceSN, deviceID, false)
	case mqtt.GISUNLINK_START_CHARGE:
		comChargeTaskStart(packet.Data.(*mqtt.ComList), deviceSN, deviceID, true)
	case mqtt.GISUNLINK_EXIT_CHARGE_TASK:
		comChargeTaskStop(packet.Data.(*mqtt.ComTaskStopTransfer), deviceSN, deviceID, false)
	case mqtt.GISUNLINK_STOP_CHARGE:
		comChargeTaskStop(packet.Data.(*mqtt.ComList), deviceSN, deviceID, true)
	case mqtt.GISUNLINK_CHARGE_NO_LOAD, mqtt.GISUNLINK_CHARGE_FINISH, mqtt.GISUNLINK_CHARGE_BREAKDOWN:
		deviceInitiativeExitComChargeTask(packet.Data.(*mqtt.ComList), deviceSN, deviceID, packet.Json.Behavior)
	}
}
