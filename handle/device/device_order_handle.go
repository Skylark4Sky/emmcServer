package device

import (
	. "GoServer/model/device"
	mqtt "GoServer/mqttPacket"
	. "GoServer/model/asyncTask"
)

func asyncDeviceChargeTaskFunc(task *AsyncTaskEntity) {
	switch task.Type {
	case ASYNC_CREATE_COM_CHARGE_TASK:
		entity := task.Entity.(*DeviceCom)
		DeviceComChargeTaskOps(entity, false)
	case ASYNC_CREATE_COM_CHARGE_TASK_ACK:
		entity := task.Entity.(*DeviceCom)
		DeviceComChargeTaskOps(entity, true)
	case ASYNC_STOP_COM_CHARGE_TASK:
	case ASYNC_STOP_COM_CHARGE_TASK_ACK:
	}
}

func createComChargeTask(task *mqtt.ComTaskStartTransfer, deviceSN string, deviceID uint64) {
	if task != nil && deviceID > 0 && len(deviceSN) > 1 {
		deviceCom := &DeviceCom{}
		deviceCom.Create(deviceID, uint64(task.Token), task.ComID)
		deviceCom.Init(task.MaxEnergy, task.MaxTime, task.MaxElectricity)
		task := NewTask()
		//task.Func = asyncDeviceChargeTaskFunc
		task.RunTaskWithTypeAndEntity(ASYNC_CREATE_COM_CHARGE_TASK, deviceCom)
	}
}

func deviceAckCreateComChargeTask(comList *mqtt.ComList, deviceSN string, deviceID uint64) {
	if len(comList.ComPort) >= 1 && len(deviceSN) > 1 && deviceID > 0 {
		comData := (comList.ComPort[0]).(mqtt.ComData)
		deviceCom := &DeviceCom{}
		deviceCom.Create(deviceID, uint64(comData.Token), comData.Id)
		deviceCom.Init(comData.MaxEnergy, comData.MaxTime,uint32(comData.MaxElectricity))
		task := NewTask()
		task.Func = asyncDeviceChargeTaskFunc
		task.RunTaskWithTypeAndEntity(ASYNC_CREATE_COM_CHARGE_TASK_ACK, deviceCom)
	}
}

func exitComChargeTask(task *mqtt.ComTaskStopTransfer, deviceSN string, deviceID uint64) {
	if task != nil && deviceID > 0 && len(deviceSN) > 1 {
		deviceCom := &DeviceCom{}
		deviceCom.Create(deviceID, uint64(task.Token), task.ComID)
		task := NewTask()
		task.Func = asyncDeviceChargeTaskFunc
		task.RunTaskWithTypeAndEntity(ASYNC_STOP_COM_CHARGE_TASK, deviceCom)
	}
}

func deviceAckExitComChargeTask(comList *mqtt.ComList, deviceSN string, deviceID uint64) {
	if len(comList.ComPort) >= 1 && len(deviceSN) > 1 && deviceID > 0 {
		comData := (comList.ComPort[0]).(mqtt.ComData)
		deviceCom := &DeviceCom{}
		deviceCom.Create(deviceID, uint64(comData.Token), comData.Id)
		task := NewTask()
		task.Func = asyncDeviceChargeTaskFunc
		task.RunTaskWithTypeAndEntity(ASYNC_STOP_COM_CHARGE_TASK_ACK, deviceCom)
	}
}
