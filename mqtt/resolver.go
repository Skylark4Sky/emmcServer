package mqtt

import "bytes"

func BinaryConversionToTaskStartTransfer(binaryData []byte) (instance *ComTaskStartTransfer) {
	instance = &ComTaskStartTransfer{}
	bytesBuf := bytes.NewBuffer(binaryData)
	instance.comID = getUint8(bytesBuf)
	instance.token = getUint32(bytesBuf)
	instance.maxEnergy = getUint32(bytesBuf)
	instance.maxElectricity = getUint32(bytesBuf)
	instance.maxTime = getUint32(bytesBuf)
	return
}

func BinaryConversionToTaskStopTransfer(binaryData []byte) (instance *ComTaskStopTransfer) {
	instance = &ComTaskStopTransfer{}
	bytesBuf := bytes.NewBuffer(binaryData)
	instance.comID = getUint8(bytesBuf)
	instance.token = getUint32(bytesBuf)
	instance.forceStop = getUint8(bytesBuf)
	return
}

func BinaryConversionToComList(binaryData []byte, behavior uint8) (instance *ComList) {
	instance = &ComList{}
	bytesBuf := bytes.NewBuffer(binaryData)

	instance.comProtoVer = getProtocolVersion(len(binaryData))

	if instance.comProtoVer == 0 {
		instance = nil
		return
	}

	if instance.comProtoVer == MAX_PROTO_VERSION0 {
		instance.signal = 0xff
	} else {
		instance.signal = (getUint8(bytesBuf) ^ 0xFF) + 1
	}

	instance.comBehavior = behavior
	instance.comNum = getUint8(bytesBuf)
	instance.comID = getBtyes(bytesBuf, uint32(instance.comNum))
	instance.comPort = make([]interface{}, instance.comNum, instance.comNum)
	instance.enableCount = 0

	for index, _ := range instance.comPort {
		com := ComData{}
		com.token = getUint32(bytesBuf)
		com.maxEnergy = getUint32(bytesBuf)
		com.useEnergy = getUint32(bytesBuf)
		com.maxTime = getUint32(bytesBuf)
		com.useTime = getUint32(bytesBuf)
		com.curElectricity = getUint32(bytesBuf)

		instance.useEnergy += com.useEnergy
		instance.useElectricity += com.curElectricity

		switch instance.comProtoVer {
		case MAX_PROTO_VERSION0:
			com.errCode = getUint8(bytesBuf)
			break
		case MAX_PROTO_VERSION1:
			com.chipReset = getUint16(bytesBuf)
			com.maxElectricity = getUint16(bytesBuf)
			com.errCode = getUint8(bytesBuf)
			break
		case MAX_PROTO_VERSION2:
			com.chipReset = getUint16(bytesBuf)
			com.maxElectricity = getUint16(bytesBuf)
			com.errCode = getUint8(bytesBuf)
			com.enable = getUint8(bytesBuf)
			instance.enableCount += com.enable
			break
		case MAX_PROTO_VERSION3:
			com.chipReset = getUint16(bytesBuf)
			com.maxElectricity = getUint16(bytesBuf)
			com.errCode = getUint8(bytesBuf)
			com.enable = getUint8(bytesBuf)
			com.behavior = getUint8(bytesBuf)
			instance.enableCount += com.enable
			break
		}

		instance.comPort[index] = com
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
		break
	case GISUNLINK_EXIT_CHARGE_TASK: //终止
		instance = BinaryConversionToTaskStopTransfer(binaryData)
		break
	case GISUNLINK_SET_CONFIG: //配置
		break
	case GISUNLINK_RESTART: //重启
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
		break
	}
	return
}
