package Mqtt

import (
	. "GoServer/utils/time"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
)

type Packet struct {
	Json *JosnPacket
	Data interface{}
}

//转译base64数据
func base64ToBinaryData(base64String string) (binary []byte) {
	decodeBytes, err := base64.StdEncoding.DecodeString(base64String)
	if err == nil {
		binary = decodeBytes
	}
	return
}

func binaryDataiToBase64(binary []byte) (base64String string) {
	buf := bytes.Buffer{}
	b64 := base64.NewEncoder(base64.StdEncoding, &buf)
	if _, err := b64.Write(binary); err != nil {
		return ""
	}
	if err := b64.Close(); err != nil {
		return ""
	}
	return buf.String()
}

// 包解析
func (packet *Packet) analysisTransferBehavior() {
	binaryData := base64ToBinaryData(packet.Json.formatData())
	packet.Data = binaryConversionToInstance(binaryData, uint8(packet.Json.Behavior))
}

//按上传行为解析包结构
func (packet *Packet) analysisAction() {
	switch packet.Json.Act {
	case TRANSFER:
		packet.analysisTransferBehavior()
		break
	case TRANSFER_RESULT:
		packet.Data = &TransferResult{}
		break
	case DEVICE_INFO:
		packet.Data = &DeviceInfo{}
		break
	case FIRMWARE_UPDATE:
		packet.Data = &UpdateState{}
		break
	}
	if packet.Data != nil {
		err := json.Unmarshal([]byte(packet.Json.Data), &packet.Data)
		if err == nil {

		}
	}
}

func MessageUnpack(Payload []byte) (ok bool, packet *Packet) {
	ok = false
	Json := &JosnPacket{}
	err := json.Unmarshal(Payload, &Json)
	if err == nil {
		ok = true
		packet = &Packet{}
		packet.Json = Json
		packet.analysisAction()
		return
	}
	return
}

func MessagePack(data interface{}) (payload string, err error) {
	if data == nil {
		return "", errors.New("包结构为空")
	}

	var base64String string = ""
	var transfer string = TRANSFER
	var behavior uint8

	switch instance := data.(type) {
	case ComTaskStartTransfer:
		behavior = GISUNLINK_CHARGE_TASK
		base64String = binaryDataiToBase64(startTransferTaskConversionToBinary(&instance))
	case ComTaskStopTransfer:
		behavior = GISUNLINK_EXIT_CHARGE_TASK
		base64String = binaryDataiToBase64(stopTransferTaskConversionToBinary(&instance))
	case ComTaskStatusQueryTransfer:
		behavior = GISUNLINK_DEVIDE_STATUS
		base64String = binaryDataiToBase64(statusQueryTransferTaskConversionToBinary(&instance))
	case DeviceSetConfigTransfer:
		behavior = GISUNLINK_SET_CONFIG
		base64String = binaryDataiToBase64(setConfigTransferTaskConversionToBinary(&instance))
	case DeviceReStartTaskTransfer:
		behavior = GISUNLINK_RESTART
		reStartDeviceTransferTaskConversionToBinary(&instance)
	}

	packet := &JosnPacket{
		Act:      transfer,
		ID:       int(GetTimestamp()),
		Ctime:    int(GetTimestamp()),
		Behavior: behavior,
		Data:     base64String,
	}

	b, err := json.Marshal(packet)
	if err != nil {
		return "", err
	}
	payload = string(b)
	return
}
