package device

import (
	. "GoServer/middleWare/dataBases/redis"
	. "GoServer/model/asyncTask"
	. "GoServer/model/device"
	mqtt "GoServer/mqttPacket"
)

func asyncDeviceChargeTaskFunc(task *AsyncTaskEntity) {
	var state uint32 = 0
	switch task.Type {
	case ASYNC_CREATE_COM_CHARGE_TASK:
		state = COM_CHARGE_START_BIT
	case ASYNC_CREATE_COM_CHARGE_TASK_ACK:
		state = COM_CHARGE_START_ACK_BIT
	case ASYNC_STOP_COM_CHARGE_TASK:
		state = COM_CHARGE_STOP_BIT
	case ASYNC_STOP_COM_CHARGE_TASK_ACK:
		state = COM_CHARGE_STOP_ACK_BIT
	case ASYNC_INITIATIVE_EXIT_COM_CHARGE_TASK:
		if initiative_type, ok := task.Param["initiative_type"]; ok {
			state = initiative_type.(uint32)
		}
	}
	entity := task.Entity.(*DeviceCharge)
	if state >= COM_CHARGE_START_BIT {
		DeviceComChargeTaskOps(entity, state)
	}
}

//结算
func settlementChargeTaskData() {

}

func comChargeTaskStart(iface interface{}, deviceSN string, deviceID uint64, ack bool) {
	if iface != nil && deviceID > 0 && len(deviceSN) > 1 {
		deviceCom := &DeviceCharge{}
		task := NewTask()
		task.Func = asyncDeviceChargeTaskFunc
		if !ack {
			entity := iface.(*mqtt.ComTaskStartTransfer)
			deviceCom.Create(deviceID, uint64(entity.Token), entity.ComID)
			deviceCom.Init(entity.MaxEnergy, entity.MaxTime, entity.MaxElectricity)
			task.RunTaskWithTypeAndEntity(ASYNC_CREATE_COM_CHARGE_TASK, deviceCom)
		} else {
			comList := iface.(*mqtt.ComList)
			entity := (comList.ComPort[0]).(mqtt.ComData)
			deviceCom.Create(deviceID, uint64(entity.Token), entity.Id)
			deviceCom.Init(entity.MaxEnergy, entity.MaxTime, uint32(entity.MaxElectricity))
			task.RunTaskWithTypeAndEntity(ASYNC_CREATE_COM_CHARGE_TASK_ACK, deviceCom)
		}
	}
}

func comChargeTaskStop(iface interface{}, deviceSN string, deviceID uint64, ack bool) {
	if iface != nil && deviceID > 0 && len(deviceSN) > 1 {
		cacheData := &CacheComData{}
		deviceCom := &DeviceCharge{}

		task := NewTask()
		task.Func = asyncDeviceChargeTaskFunc

		if !ack {
			entity := iface.(*mqtt.ComTaskStopTransfer)
			Redis().GetDeviceComDataFormRedis(deviceSN, entity.ComID, cacheData)

			deviceCom.Create(deviceID, uint64(entity.Token), entity.ComID)
			deviceCom.ChangeValue(cacheData.UseEnergy, cacheData.UseTime, cacheData.MaxChargeElectricity, cacheData.AveragePower, cacheData.MaxPower)
			task.RunTaskWithTypeAndEntity(ASYNC_STOP_COM_CHARGE_TASK, deviceCom)
		} else {
			comList := iface.(*mqtt.ComList)
			entity := (comList.ComPort[0]).(mqtt.ComData)
			Redis().GetDeviceComDataFormRedis(deviceSN, entity.Id, cacheData)

			deviceCom.Create(deviceID, uint64(entity.Token), entity.Id)
			deviceCom.ChangeValue(cacheData.UseEnergy, cacheData.UseTime, cacheData.MaxChargeElectricity, cacheData.AveragePower, cacheData.MaxPower)
			task.RunTaskWithTypeAndEntity(ASYNC_STOP_COM_CHARGE_TASK_ACK, deviceCom)
		}
	}
}

func deviceInitiativeExitComChargeTask(comList *mqtt.ComList, deviceSN string, deviceID uint64, behavior uint8) {
	if len(comList.ComPort) >= 1 && len(deviceSN) > 1 && deviceID > 0 {
		comData := (comList.ComPort[0]).(mqtt.ComData)

		cacheData := &CacheComData{}
		Redis().GetDeviceComDataFormRedis(deviceSN, comData.Id, cacheData)

		deviceCom := &DeviceCharge{}
		deviceCom.Create(deviceID, uint64(comData.Token), comData.Id)
		deviceCom.Init(comData.MaxEnergy, comData.MaxTime, uint32(comData.MaxElectricity))
		deviceCom.ChangeValue(cacheData.UseEnergy, cacheData.UseTime, cacheData.MaxChargeElectricity, cacheData.AveragePower, cacheData.MaxPower)

		task := NewTask()
		switch behavior {
		case mqtt.GISUNLINK_CHARGE_FINISH:
			task.Param = map[string]interface{}{"initiative_type": COM_CHARGE_FINISH_BIT}
		case mqtt.GISUNLINK_CHARGE_BREAKDOWN:
			task.Param = map[string]interface{}{"initiative_type": COM_CHARGE_BREAKDOWN_BIT}
		case mqtt.GISUNLINK_CHARGE_NO_LOAD:
			task.Param = map[string]interface{}{"initiative_type": COM_CHARGE_NO_LOAD_BIT}
		}
		task.Func = asyncDeviceChargeTaskFunc
		task.RunTaskWithTypeAndEntity(ASYNC_INITIATIVE_EXIT_COM_CHARGE_TASK, deviceCom)
	}
}
