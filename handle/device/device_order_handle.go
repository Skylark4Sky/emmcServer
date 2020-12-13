package device

import (
	. "GoServer/model/device"
	. "GoServer/model"
	mqtt "GoServer/mqttPacket"
)

func createComChargeTask(task *mqtt.ComTaskStartTransfer, deviceSN string, deviceID uint64) {
	if (task != nil && deviceID > 0 && len(deviceSN) > 1) {
		deviceCom := &DeviceCom{}
		deviceCom.Create(deviceID,uint64(task.Token), task.ComID,task.MaxEnergy,task.MaxTime,task.MaxElectricity)
		CreateAsyncSQLTask(ASYNC_CREATE_COM_CHARGE_TASK, deviceCom)
	}
}

func exitComChargeTask(task *mqtt.ComTaskStopTransfer,deviceSN string, deviceID uint64) {
	if (task != nil && deviceID > 0 && len(deviceSN) > 1) {

	}
}