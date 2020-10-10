package Mqtt

import (
	. "GoServer/utils/float64"
	"bytes"
)

func BinaryConversionToTaskStartTransfer(binaryData []byte) (instance *ComTaskStartTransfer) {
	instance = &ComTaskStartTransfer{}
	bytesBuf := bytes.NewBuffer(binaryData)
	instance.comID = (GetUint8(bytesBuf))
	instance.token = (GetUint32(bytesBuf))
	instance.maxEnergy = (GetUint32(bytesBuf))
	instance.maxElectricity = (GetUint32(bytesBuf))
	instance.maxTime = (GetUint32(bytesBuf))
	return
}

func BinaryConversionToTaskStopTransfer(binaryData []byte) (instance *ComTaskStopTransfer) {
	instance = &ComTaskStopTransfer{}
	bytesBuf := bytes.NewBuffer(binaryData)
	instance.comID = (GetUint8(bytesBuf))
	instance.token = (GetUint32(bytesBuf))
	instance.forceStop = (GetUint8(bytesBuf))
	return
}

func BinaryConversionToTaskStatusQueryTransfer(binaryData []byte) (instance *ComTaskStatusQueryTransfer) {
	instance = &ComTaskStatusQueryTransfer{}
	bytesBuf := bytes.NewBuffer(binaryData)
	instance.comID = (GetUint8(bytesBuf))
	return
}

func BinaryConversionToTaskSetConfigTransfer(binaryData []byte) (instance *DeviceONLoadTimeSetConfigTransfer) {
	instance = &DeviceONLoadTimeSetConfigTransfer{}
	bytesBuf := bytes.NewBuffer(binaryData)
	instance.time = (GetUint8(bytesBuf))
	return
}

func BinaryConversionToReStartDeviceTaskTransfer(binaryData []byte) (instance *DeviceReStartTaskTransfer) {
	instance = &DeviceReStartTaskTransfer{}
	return
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

func BinaryConversionToComList(binaryData []byte, behavior uint8) (instance *ComList) {
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
		instance.Signal = (GetUint8(bytesBuf) ^ 0xFF) + 1
	}

	instance.ComBehavior = behavior
	instance.ComNum = GetUint8(bytesBuf)
	instance.ComID = GetBtyes(bytesBuf, uint32(instance.ComNum))
	instance.ComPort = make([]interface{}, instance.ComNum, instance.ComNum)
	instance.EnableCount = 0

	for index, _ := range instance.ComPort {
		com := ComData{}
		com.Token = GetUint32(bytesBuf)
		com.MaxEnergy = GetUint32(bytesBuf)
		com.UseEnergy = GetUint32(bytesBuf)
		com.MaxTime = GetUint32(bytesBuf)
		com.UseTime = GetUint32(bytesBuf)
		com.CurElectricity = GetUint32(bytesBuf)

		instance.UseEnergy += com.UseEnergy
		instance.UseElectricity += com.CurElectricity

		switch instance.ComProtoVer {
		case MAX_PROTO_VERSION0:
			com.ErrCode = GetUint8(bytesBuf)
			break
		case MAX_PROTO_VERSION1:
			com.ChipReset = GetUint16(bytesBuf)
			com.MaxElectricity = GetUint16(bytesBuf)
			com.ErrCode = GetUint8(bytesBuf)
			break
		case MAX_PROTO_VERSION2:
			com.ChipReset = GetUint16(bytesBuf)
			com.MaxElectricity = GetUint16(bytesBuf)
			com.ErrCode = GetUint8(bytesBuf)
			com.Enable = GetUint8(bytesBuf)
			instance.EnableCount += com.Enable
			break
		case MAX_PROTO_VERSION3:
			com.ChipReset = GetUint16(bytesBuf)
			com.MaxElectricity = GetUint16(bytesBuf)
			com.ErrCode = GetUint8(bytesBuf)
			com.Enable = GetUint8(bytesBuf)
			com.Behavior = GetUint8(bytesBuf)
			instance.EnableCount += com.Enable
			break
		}
		tempCurPower := com.CurPower
		com.CurPower = CalculateCurComPower(CUR_VOLTAGE, com.CurElectricity, 2)
		com.AveragePower = CalculateCurAverageComPower(com.UseEnergy, com.UseTime, 2)

		if CmpPower(com.CurPower,tempCurPower) == 1 {
			com.MaxPower = com.CurPower
		}
		instance.ComPort[index] = com
	}
	return
}

func BinaryConversionToInstance(binaryData []byte, behavior uint8) (instance interface{}) {
	instance = nil
	switch behavior {
	//下发
	case GISUNLINK_CHARGE_TASK: //任务
		instance = BinaryConversionToTaskStartTransfer(binaryData)
		break
	case GISUNLINK_DEVIDE_STATUS: //查询
		instance = BinaryConversionToTaskStatusQueryTransfer(binaryData)
		break
	case GISUNLINK_EXIT_CHARGE_TASK: //终止
		instance = BinaryConversionToTaskStopTransfer(binaryData)
		break
	case GISUNLINK_SET_CONFIG: //配置
		instance = BinaryConversionToTaskSetConfigTransfer(binaryData)
		break
	case GISUNLINK_RESTART: //重启
		instance = BinaryConversionToReStartDeviceTaskTransfer(binaryData)
		break
	//上报
	case GISUNLINK_START_CHARGE: //开始
		instance = BinaryConversionToComList(binaryData, GISUNLINK_START_CHARGE)
		break
	case GISUNLINK_CHARGEING: //执行中
		instance = BinaryConversionToComList(binaryData, GISUNLINK_CHARGEING)
		break
	case GISUNLINK_CHARGE_FINISH: //完成
		instance = BinaryConversionToComList(binaryData, GISUNLINK_CHARGE_FINISH)
		break
	case GISUNLINK_CHARGE_LEISURE: //空闲中
		instance = BinaryConversionToComList(binaryData, GISUNLINK_CHARGE_LEISURE)
		break
	case GISUNLINK_CHARGE_BREAKDOWN: //故障
		instance = BinaryConversionToComList(binaryData, GISUNLINK_CHARGE_BREAKDOWN)
		break
	case GISUNLINK_CHARGE_NO_LOAD: //空载
		instance = BinaryConversionToComList(binaryData, GISUNLINK_CHARGE_NO_LOAD)
		break
	case GISUNLINK_UPDATE_FIRMWARE: //升级
		instance = BinaryConversionToComList(binaryData, GISUNLINK_UPDATE_FIRMWARE)
		break
	case GISUNLINK_COM_UPDATE: //参数刷新
		instance = BinaryConversionToComList(binaryData, GISUNLINK_COM_UPDATE)
		break
	case GISUNLINK_STOP_CHARGE: //停止
		instance = BinaryConversionToComList(binaryData, GISUNLINK_STOP_CHARGE)
		break
	case GISUNLINK_COM_NO_UPDATE: //参数没有刷新
		instance = BinaryConversionToComList(binaryData, GISUNLINK_COM_NO_UPDATE)
		break
	}
	return
}
