package mqttPacket

import (
	"bytes"
)

func binaryConversionToStartTransferTask(binaryData []byte) (instance *ComTaskStartTransfer) {
	instance = &ComTaskStartTransfer{}
	bytesBuf := bytes.NewBuffer(binaryData)
	instance.ComID = (getUint8(bytesBuf))
	instance.Token = (getUint32(bytesBuf))
	instance.MaxEnergy = (getUint32(bytesBuf))
	instance.MaxElectricity = (getUint32(bytesBuf))
	instance.MaxTime = (getUint32(bytesBuf))
	return
}

func startTransferTaskConversionToBinary(instance *ComTaskStartTransfer) []byte {
	binaryData := new(bytes.Buffer)
	setUint8(binaryData, instance.ComID)
	setUint32(binaryData, instance.Token)
	setUint32(binaryData, instance.MaxEnergy)
	setUint32(binaryData, instance.MaxElectricity)
	setUint32(binaryData, instance.MaxTime)
	return binaryData.Bytes()
}

func binaryConversionToStopTransferTask(binaryData []byte) (instance *ComTaskStopTransfer) {
	instance = &ComTaskStopTransfer{}
	bytesBuf := bytes.NewBuffer(binaryData)
	instance.ComID = (getUint8(bytesBuf))
	instance.Token = (getUint32(bytesBuf))
	instance.ForceStop = (getUint8(bytesBuf))
	return
}

func stopTransferTaskConversionToBinary(instance *ComTaskStopTransfer) []byte {
	binaryData := new(bytes.Buffer)
	setUint8(binaryData, instance.ComID)
	setUint32(binaryData, instance.Token)
	setUint8(binaryData, instance.ForceStop)
	return binaryData.Bytes()
}

func binaryConversionToStatusQueryTransferTask(binaryData []byte) (instance *ComTaskStatusQueryTransfer) {
	instance = &ComTaskStatusQueryTransfer{}
	bytesBuf := bytes.NewBuffer(binaryData)
	instance.ComID = (getUint8(bytesBuf))
	return
}

func statusQueryTransferTaskConversionToBinary(instance *ComTaskStatusQueryTransfer) []byte {
	binaryData := new(bytes.Buffer)
	setUint8(binaryData, instance.ComID)
	return binaryData.Bytes()
}

func binaryConversionToSetConfigTransferTask(binaryData []byte) (instance *DeviceSetConfigTransfer) {
	instance = &DeviceSetConfigTransfer{}
	bytesBuf := bytes.NewBuffer(binaryData)
	instance.Time = (getUint8(bytesBuf))
	return
}

func setConfigTransferTaskConversionToBinary(instance *DeviceSetConfigTransfer) []byte {
	binaryData := new(bytes.Buffer)
	setUint8(binaryData, instance.Time)
	return binaryData.Bytes()
}

func binaryConversionToReStartDeviceTransferTask(binaryData []byte) (instance *DeviceReStartTaskTransfer) {
	instance = &DeviceReStartTaskTransfer{}
	return
}

func reStartDeviceTransferTaskConversionToBinary(instance *DeviceReStartTaskTransfer) []byte {
	return nil
}

func getTransferVersion(length int) (version uint16) {
	version = 0
	if length == MIN_PROTO_VERSION0 || length == MAX_PROTO_VERSION0 {
		version = MAX_PROTO_VERSION0
	} else if length == MIN_PROTO_VERSION1 || length == MAX_PROTO_VERSION1 {
		version = MAX_PROTO_VERSION1
	} else if length == MIN_PROTO_VERSION2 || length == MAX_PROTO_VERSION2 {
		version = MAX_PROTO_VERSION2
	} else if length == MIN_PROTO_VERSION3 || length == MAX_PROTO_VERSION3 {
		version = MAX_PROTO_VERSION3
	} else if length == MIN_PROTO_VERSION3 || length == MAX_PROTO_VERSION4 {
		version = MAX_PROTO_VERSION3
	}
	return
}

func binaryConversionToComList(binaryData []byte, behavior uint8) (instance *ComList) {
	instance = &ComList{}
	bytesBuf := bytes.NewBuffer(binaryData)
	dataLen := len(binaryData)
	instance.ComProtoVer = getTransferVersion(dataLen)

	if instance.ComProtoVer == 0 {
		instance = nil
		return
	}

	if instance.ComProtoVer == MAX_PROTO_VERSION0 {
		instance.Signal = 0xff
	} else {
		instance.Signal = (getUint8(bytesBuf) ^ 0xFF) + 1
	}

	instance.ComBehavior = behavior
	instance.ComNum = getUint8(bytesBuf)
	instance.ComID = getBtyes(bytesBuf, uint32(instance.ComNum))
	instance.ComPort = make([]interface{}, instance.ComNum, instance.ComNum)
	instance.EnableCount = 0

	for index, _ := range instance.ComPort {
		com := ComData{}
		com.Token = getUint32(bytesBuf)
		com.MaxEnergy = getUint32(bytesBuf)
		com.UseEnergy = getUint32(bytesBuf)
		com.MaxTime = getUint32(bytesBuf)
		com.UseTime = getUint32(bytesBuf)
		com.CurElectricity = getUint32(bytesBuf)

		instance.UseEnergy += com.UseEnergy
		instance.UseElectricity += com.CurElectricity

		switch instance.ComProtoVer {
		case MAX_PROTO_VERSION0:
			com.ErrCode = getUint8(bytesBuf)
			break
		case MAX_PROTO_VERSION1:
			com.ChipReset = getUint16(bytesBuf)
			com.MaxElectricity = getUint16(bytesBuf)
			com.ErrCode = getUint8(bytesBuf)
			break
		case MAX_PROTO_VERSION2:
			com.ChipReset = getUint16(bytesBuf)
			com.MaxElectricity = getUint16(bytesBuf)
			com.ErrCode = getUint8(bytesBuf)
			com.Enable = getUint8(bytesBuf)
			instance.EnableCount += com.Enable
			break
		case MAX_PROTO_VERSION3:
			com.ChipReset = getUint16(bytesBuf)
			com.MaxElectricity = getUint16(bytesBuf)
			com.ErrCode = getUint8(bytesBuf)
			com.Enable = getUint8(bytesBuf)
			com.Behavior = getUint8(bytesBuf)
			instance.EnableCount += com.Enable
			break
		}
		instance.ComPort[index] = com
	}
	return
}

func binaryConversionToInstance(binaryData []byte, behavior uint8) (instance interface{}) {
	instance = nil
	switch behavior {
	//下发
	case GISUNLINK_CHARGE_TASK: //任务
		instance = binaryConversionToStartTransferTask(binaryData)
		break
	case GISUNLINK_DEVIDE_STATUS: //查询
		instance = binaryConversionToStatusQueryTransferTask(binaryData)
		break
	case GISUNLINK_EXIT_CHARGE_TASK: //终止
		instance = binaryConversionToStopTransferTask(binaryData)
		break
	case GISUNLINK_SET_CONFIG: //配置
		instance = binaryConversionToSetConfigTransferTask(binaryData)
		break
	case GISUNLINK_RESTART: //重启
		instance = binaryConversionToReStartDeviceTransferTask(binaryData)
		break
	//上报
	case GISUNLINK_START_CHARGE: //开始
		instance = binaryConversionToComList(binaryData, GISUNLINK_START_CHARGE)
		break
	case GISUNLINK_CHARGEING: //执行中
		instance = binaryConversionToComList(binaryData, GISUNLINK_CHARGEING)
		break
	case GISUNLINK_CHARGE_FINISH: //完成
		instance = binaryConversionToComList(binaryData, GISUNLINK_CHARGE_FINISH)
		break
	case GISUNLINK_CHARGE_LEISURE: //空闲中
		instance = binaryConversionToComList(binaryData, GISUNLINK_CHARGE_LEISURE)
		break
	case GISUNLINK_CHARGE_BREAKDOWN: //故障
		instance = binaryConversionToComList(binaryData, GISUNLINK_CHARGE_BREAKDOWN)
		break
	case GISUNLINK_CHARGE_NO_LOAD: //空载
		instance = binaryConversionToComList(binaryData, GISUNLINK_CHARGE_NO_LOAD)
		break
	case GISUNLINK_UPDATE_FIRMWARE: //升级
		instance = binaryConversionToComList(binaryData, GISUNLINK_UPDATE_FIRMWARE)
		break
	case GISUNLINK_COM_UPDATE: //参数刷新
		instance = binaryConversionToComList(binaryData, GISUNLINK_COM_UPDATE)
		break
	case GISUNLINK_STOP_CHARGE: //停止
		instance = binaryConversionToComList(binaryData, GISUNLINK_STOP_CHARGE)
		break
	case GISUNLINK_COM_NO_UPDATE: //参数没有刷新
		instance = binaryConversionToComList(binaryData, GISUNLINK_COM_NO_UPDATE)
		break
	}
	return
}
