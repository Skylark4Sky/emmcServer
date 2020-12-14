package device

import (
	. "GoServer/middleWare/dataBases/redis"
	. "GoServer/model"
	mqtt "GoServer/mqttPacket"
	. "GoServer/utils/float64"
	. "GoServer/utils/log"
	. "GoServer/utils/time"
	"time"
)

//构建缓存数据实体
func calculateComData(comData *mqtt.ComData) *CacheComData {
	cacheComData := &CacheComData{
		ComData:              *comData,
		MaxChargeElectricity: comData.CurElectricity,
		CurPower:             0,
		AveragePower:         0,
		MaxPower:             0,
	}
	return cacheComData
}

//比较数据
func analyseComListData(deviceSN string, comList *mqtt.ComList, cacheComList map[uint8]CacheComData) {
	//不存在缓存数据直接返回
	if len(cacheComList) <= 0 {
		return
	}

	for _, comID := range comList.ComID {
		var index uint8 = comID
		if comList.ComNum <= 5 {
			index = (comID % 5)
		}
		comData := (comList.ComPort[int(index)]).(mqtt.ComData)
		comData.Id = comID
		cacheComData := cacheComList[comData.Id]
		//未开启时，检测值是否有变化
		if comData.Enable == COM_DISENABLE {
			//大于50ma保存异常数值
			if (comData.CurElectricity >= 50) && (comData.CurElectricity != cacheComData.CurElectricity) {
				SystemLog(" Time:", TimeFormat(time.Now()), " ", deviceSN, " 端口:", comData.Id, " 异常---当前值:", comData.CurElectricity, "上一次值为:", cacheComData.CurElectricity)
			}
		} else {
			//comData.MaxPower
		}
		if comData.Token != cacheComData.Token {
			SystemLog(" NewToKen:", comData.Token, " OldToKen:", cacheComData.Token)
		}
	}
}

// 运行状态 & 空闲状态 处理
func deviceStateHandle(comList *mqtt.ComList, deviceSN string, deviceID uint64) {
	//批量读当前设备所有接口
	cacheComList := BatchReadDeviceComDataiFromRedis(deviceSN)
	//数据分析
	analyseComListData(deviceSN, comList, cacheComList)
	//批量写入数据到缓存 并 统计 功率
	if len(cacheComList) > 1 {
		BatchWriteDeviceComDataToRedis(deviceSN, comList, func(comData *mqtt.ComData) *CacheComData {
			cacheComData := calculateComData(comData)
			cacheCom := cacheComList[comData.Id]
			//端口启用状态
			if comData.Enable == COM_ENABLE {
				// 数值计算
				cacheComData.CurPower = CalculateCurComPower(CUR_VOLTAGE, cacheComData.CurElectricity, 2)                //当前功率
				cacheComData.AveragePower = CalculateCurAverageComPower(cacheComData.UseEnergy, cacheComData.UseTime, 2) //平均功率
				//最大功率
				if CmpPower(cacheComData.CurPower, cacheCom.MaxPower) >= 1 {
					cacheComData.MaxPower = cacheComData.CurPower
				}
				// 最大电流记录
				if comData.CurElectricity > cacheCom.MaxChargeElectricity {
					cacheComData.MaxChargeElectricity = comData.CurElectricity
				}
			}
			return cacheComData
		})
	} else {
		BatchWriteDeviceComDataToRedis(deviceSN, comList, func(comData *mqtt.ComData) *CacheComData {
			return calculateComData(comData)
		})
	}

	//刷新设备状态
	refreshDeviceStatus(deviceSN, deviceID, &DeviceStatus{
		Behavior:     comList.ComBehavior,
		Signal:       int8(comList.Signal),
		Worker:       Redis().TatolWorkerByDevice(deviceSN, BatchReadDeviceComDataiFromRedis(deviceSN)),
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
		createComChargeTask(packet.Data.(*mqtt.ComTaskStartTransfer), deviceSN, deviceID)
	case mqtt.GISUNLINK_START_CHARGE:
		deviceAckCreateComChargeTask(packet.Data.(*mqtt.ComList), deviceSN, deviceID)
	case mqtt.GISUNLINK_EXIT_CHARGE_TASK:
		exitComChargeTask(packet.Data.(*mqtt.ComTaskStopTransfer), deviceSN, deviceID)
	case mqtt.GISUNLINK_CHARGE_NO_LOAD, mqtt.GISUNLINK_CHARGE_FINISH, mqtt.GISUNLINK_STOP_CHARGE, mqtt.GISUNLINK_CHARGE_BREAKDOWN:
		deviceAckExitComChargeTask(packet.Data.(*mqtt.ComList), deviceSN, deviceID)
	}
}
