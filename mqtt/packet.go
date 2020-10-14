package Mqtt

import (
	"encoding/base64"
	"encoding/json"
	//	"fmt"
)

type Packet struct {
	Json     *JosnPacket
	JsonData interface{}
}

//转译base64数据
func translateBinaryData(base64String string) (binary []byte) {
	decodeBytes, err := base64.StdEncoding.DecodeString(base64String)
	if err == nil {
		binary = decodeBytes
	}
	return
}

// 包解析
func (packet *Packet) analysisTransferBehavior() {
	binaryData := translateBinaryData(packet.Json.formatData())
	packet.JsonData = binaryConversionToInstance(binaryData, uint8(packet.Json.Behavior))
}

//按上传行为解析包结构
func (packet *Packet) analysisAction() {
	switch packet.Json.Act {
	case TRANSFER:
		packet.analysisTransferBehavior()
		break
	case TRANSFER_RESULT:
		packet.JsonData = &TransferResult{}
		break
	case DEVICE_INFO:
		packet.JsonData = &DeviceInfo{}
		break
	case FIRMWARE_UPDATE:
		packet.JsonData = &UpdateState{}
		break
	}

	if packet.JsonData != nil {
		err := json.Unmarshal([]byte(packet.Json.Data), &packet.JsonData)
		if err == nil {

		}
	}
}

func MessageHandler(Payload []byte) (ok bool, packet *Packet) {
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
