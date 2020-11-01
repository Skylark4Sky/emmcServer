package device

import (
	. "GoServer/handle/device"
	mqtt "GoServer/mqttPacket"
	. "GoServer/utils/respond"
	"github.com/gin-gonic/gin"
)

func DeviceStartCharge(ctx *gin.Context) {
	var postData mqtt.ComTaskStartTransfer
	postData.ComID = 1
	postData.Token = 123456789
	postData.MaxEnergy = 9000
	postData.MaxElectricity = 2272 //uint32(CalculateMaxComElectricity(500))
	postData.MaxTime = 6000

	payload, _ := mqtt.MessagePack(postData)

	mqttMsg := &MqMsg{}

	mqttMsg.Send("47.106.235.93:1883", "57ff69067878495148300967", payload)

	RespondMessage(ctx, CreateErrorMessage(PARAM_ERROR, payload))
}

func DeviceStopCharge(ctx *gin.Context) {
	var postData mqtt.ComTaskStopTransfer
	postData.ComID = 1
	postData.Token = 123456789
	postData.ForceStop = 1

	payload, _ := mqtt.MessagePack(postData)

	mqttMsg := &MqMsg{}
	mqttMsg.Send("47.106.235.93:1883", "57ff69067878495148300967", payload)
	RespondMessage(ctx, CreateErrorMessage(PARAM_ERROR, payload))

}

func DeviceStatusQuery(ctx *gin.Context) {

}

func DeviceNoLoadSetting(ctx *gin.Context) {

}

func DeviceRestart(ctx *gin.Context) {

}

func DeviceUpdateFirmware(ctx *gin.Context) {

}