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

func TranslateBinaryData(base64String string) (binary []byte) {
	decodeBytes, err := base64.StdEncoding.DecodeString(base64String)
	if err == nil {
		binary = decodeBytes
	}
	return
}

func (packet *Packet) AnalysisTransferBehavior() {
	binaryData := TranslateBinaryData(packet.Json.FormatData())
	packet.JsonData = BinaryConversionToInstance(binaryData, uint8(packet.Json.Behavior))
}

func (packet *Packet) AnalysisAction() {
	switch packet.Json.Act {
	case TRANSFER:
		packet.AnalysisTransferBehavior()
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
		packet.AnalysisAction()
		return
	}
	return
}
