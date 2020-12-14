package device

import (
	. "GoServer/model"
	. "GoServer/model/device"
	mqtt "GoServer/mqttPacket"
)

func taskHandle(task *AsyncSQLTask) {

}

func createComChargeTask(task *mqtt.ComTaskStartTransfer, deviceSN string, deviceID uint64) {

	CreateAsyncSQLTaskWithCallback(ASYNC_UPDATE_DEVICE_STATUS, request, lock, syncDeviceStatusTaskFunc)
	if task != nil && deviceID > 0 && len(deviceSN) > 1 {
		deviceCom := &DeviceCom{}
		deviceCom.Create(deviceID, uint64(task.Token), task.ComID, task.MaxEnergy, task.MaxTime, task.MaxElectricity)
		CreateAsyncSQLTask(asyncType, deviceCom)
	}
}

func deviceAckCreateComChargeTask(comList *mqtt.ComList, deviceSN string, deviceID uint64) {
	if len(comList.ComPort) >= 1 {
		comData := (comList.ComPort[0]).(mqtt.ComData)
		taskStart := &mqtt.ComTaskStartTransfer{
			ComID:          comData.Id,
			Token:          comData.Token,
			MaxEnergy:      comData.MaxEnergy,
			MaxElectricity: uint32(comData.MaxElectricity),
			MaxTime:        comData.MaxTime,
		}
		createComChargeTask(ASYNC_CREATE_COM_CHARGE_TASK_ACK, taskStart, deviceSN, deviceID)
	}
}

func exitComChargeTask(task *mqtt.ComTaskStopTransfer, deviceSN string, deviceID uint64) {
	if task != nil && deviceID > 0 && len(deviceSN) > 1 {
		deviceCom := &DeviceCom{
			DeviceID: deviceID,
			ChargeID: uint64(task.Token),
			ComID:    task.ComID,
		}
		CreateAsyncSQLTask(asyncType, deviceCom)
	}
}

func deviceAckExitComChargeTask(comList *mqtt.ComList, deviceSN string, deviceID uint64) {

}
