package device

import (
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

func createComChargeTask(task *mqtt.ComTaskStartTransfer, deviceSN string, deviceID uint64) {
	if task != nil && deviceID > 0 && len(deviceSN) > 1 {
		deviceCom := &DeviceCharge{}
		deviceCom.Create(deviceID, uint64(task.Token), task.ComID)
		deviceCom.Init(task.MaxEnergy, task.MaxTime, task.MaxElectricity)
		task := NewTask()
		task.Func = asyncDeviceChargeTaskFunc
		task.RunTaskWithTypeAndEntity(ASYNC_CREATE_COM_CHARGE_TASK, deviceCom)
	}
}

func deviceAckCreateComChargeTask(comList *mqtt.ComList, deviceSN string, deviceID uint64) {
	if len(comList.ComPort) >= 1 && len(deviceSN) > 1 && deviceID > 0 {
		comData := (comList.ComPort[0]).(mqtt.ComData)
		deviceCom := &DeviceCharge{}
		deviceCom.Create(deviceID, uint64(comData.Token), comData.Id)
		deviceCom.Init(comData.MaxEnergy, comData.MaxTime, uint32(comData.MaxElectricity))
		task := NewTask()
		task.Func = asyncDeviceChargeTaskFunc
		task.RunTaskWithTypeAndEntity(ASYNC_CREATE_COM_CHARGE_TASK_ACK, deviceCom)
	}
}

func exitComChargeTask(task *mqtt.ComTaskStopTransfer, deviceSN string, deviceID uint64) {
	if task != nil && deviceID > 0 && len(deviceSN) > 1 {
		deviceCom := &DeviceCharge{}
		deviceCom.Create(deviceID, uint64(task.Token), task.ComID)
		task := NewTask()
		task.Func = asyncDeviceChargeTaskFunc
		task.RunTaskWithTypeAndEntity(ASYNC_STOP_COM_CHARGE_TASK, deviceCom)
	}
}

func deviceAckExitComChargeTask(comList *mqtt.ComList, deviceSN string, deviceID uint64) {
	if len(comList.ComPort) >= 1 && len(deviceSN) > 1 && deviceID > 0 {
		comData := (comList.ComPort[0]).(mqtt.ComData)
		deviceCom := &DeviceCharge{}
		deviceCom.Create(deviceID, uint64(comData.Token), comData.Id)
		deviceCom.Init(comData.MaxEnergy, comData.MaxTime, uint32(comData.MaxElectricity))
		task := NewTask()
		task.Func = asyncDeviceChargeTaskFunc
		task.RunTaskWithTypeAndEntity(ASYNC_STOP_COM_CHARGE_TASK_ACK, deviceCom)
	}
}

func deviceInitiativeExitComChargeTask(comList *mqtt.ComList, deviceSN string, deviceID uint64, behavior uint8) {
	if len(comList.ComPort) >= 1 && len(deviceSN) > 1 && deviceID > 0 {
		comData := (comList.ComPort[0]).(mqtt.ComData)
		deviceCom := &DeviceCharge{}
		deviceCom.Create(deviceID, uint64(comData.Token), comData.Id)
		deviceCom.Init(comData.MaxEnergy, comData.MaxTime, uint32(comData.MaxElectricity))
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
