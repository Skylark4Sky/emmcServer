package device

import (
	. "GoServer/middleWare/dataBases/mysql"
	. "GoServer/middleWare/dataBases/redis"
	. "GoServer/model/device"
	mqtt "GoServer/mqttPacket"
	. "GoServer/utils/float64"
	. "GoServer/utils/log"
	. "GoServer/utils/time"
	. "GoServer/model"
	"time"
	. "GoServer/utils/string"
)

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
func saveDeviceTransferDataOps(serverNode string, device_sn string, packet *mqtt.Packet) {
	var comNum uint8 = 0
	switch packet.Json.Behavior {
	case mqtt.GISUNLINK_CHARGEING, mqtt.GISUNLINK_CHARGE_LEISURE: //运行中,空闲中
		comList := packet.Data.(*mqtt.ComList)
		comNum = comList.ComNum
		break
	case mqtt.GISUNLINK_START_CHARGE, mqtt.GISUNLINK_CHARGE_FINISH, mqtt.GISUNLINK_CHARGE_NO_LOAD, mqtt.GISUNLINK_CHARGE_BREAKDOWN: //开始,完成,空载,故障
		comList := packet.Data.(*mqtt.ComList)
		comNum = comList.ComNum
		for _, comID := range comList.ComID {
			comNum = comID
		}
		break
	}

	log := &DeviceTransferLog{
		TransferID:   int64(packet.Json.ID),
		DeviceID:     Redis().GetDeviceIDFromRedis(device_sn),
		TransferAct:  packet.Json.Act,
		DeviceSn:     device_sn,
		ComNum:       comNum,
		TransferData: packet.Json.Data,
		Behavior:     packet.Json.Behavior,
		ServerNode:   serverNode,
		TransferTime: int64(packet.Json.Ctime),
	}
	createDeviceTransferLog(log)
}

func deviceExpiredMsgOps(pattern, channel, message string) {
	deviceSN := GetDeviceSN(message,":")
	if deviceID := Redis().GetDeviceIDFromRedis(deviceSN); deviceID != 0 {
		switch message {
		case GetDeviceTokenKey(deviceSN): //这里处理过期key
			{
				device := &DeviceInfo{}
				deviceUpdateMap := map[string]interface{}{"update_time": GetTimestampMs(), "status": DEVICE_OFFLINE}
				CreateAsyncSQLTaskWithUpdateMap(ASYNC_UP_DEVICE_INFO, device, deviceUpdateMap)
			}
		case GetDeviceIDKey(deviceSN):
			{
			}
		case GetRawDataKey(deviceSN):
			{
			}
		case GetComdDataKey(deviceSN):
			{
			}
		case GetDeviceInfoKey(deviceSN):
			{
			}
		}
	}
}

//比较数据
func analyseComData(tokenKey string, newData *mqtt.ComList, cacheData map[uint8]mqtt.ComData) {
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

func deviceActBehaviorDataOps(packet *mqtt.Packet, cacheKey string, playload string) {
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
				BatchWriteDeviceComDataToRedis(cacheKey, comList, func(comData *mqtt.ComData) {
					cacherData := cacherComData[comData.Id]
					comData.MaxPower = cacherData.MaxPower
					if CmpPower(comData.CurPower, cacherData.MaxPower) == 1 {
						comData.MaxPower = comData.CurPower
					}
				})
			} else {
				BatchWriteDeviceComDataToRedis(cacheKey, comList, func(comData *mqtt.ComData) {})
			}

			deviceStatus := &DeviceStatus{
				Behavior:     comList.ComBehavior,
				Signal:       int8(comList.Signal),
				Worker:       Redis().TatolWorkerByDevice(cacheKey, BatchReadDeviceComDataiFromRedis(cacheKey)),
				ProtoVersion: comList.ComProtoVer,
			}

			var timeout int64 = CHARGEING_TIME
			if packet.Json.Behavior == mqtt.GISUNLINK_CHARGE_LEISURE {
				timeout = LEISURE_TIME
			}

			//更新令牌时间
			Redis().UpdateDeviceTokenExpiredTime(cacheKey, deviceStatus, timeout)
			//写入原始数据
			Redis().UpdateDeviceRawDataToRedis(cacheKey, playload)
		}
	}
}
