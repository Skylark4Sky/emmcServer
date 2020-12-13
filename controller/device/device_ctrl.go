package device

import (
	. "GoServer/handle/device"
	. "GoServer/middleWare/extension"
	mqtt "GoServer/mqttPacket"
	. "GoServer/utils/respond"
	"github.com/gin-gonic/gin"
)

func StartCharge(ctx *gin.Context) {
	//	userID := ctx.MustGet(JwtCtxUidKey)

	//	if userID.(uint64) <= 0 {
	//		RespondMessage(ctx, CreateErrorMessage(PARAM_ERROR, nil))
	//	}

	var postData mqtt.ComTaskStartTransfer
	postData.ComID = 3
	postData.Token = 123456789
	postData.MaxEnergy = 19000
	postData.MaxElectricity = 2272 //uint32(CalculateMaxComElectricity(500))
	postData.MaxTime = 6000

	payload, _ := mqtt.MessagePack(postData)

	mqttMsg := &MqMsg{}

	mqttMsg.Send("47.106.235.93:1883", "57ff69067878495148300967", payload)

	RespondMessage(ctx, CreateErrorMessage(PARAM_ERROR, payload))
}

func StopCharge(ctx *gin.Context) {
	//	userID := ctx.MustGet(JwtCtxUidKey)
	//
	//	if userID.(uint64) <= 0 {
	//		RespondMessage(ctx, CreateErrorMessage(PARAM_ERROR, nil))
	//	}

	var postData mqtt.ComTaskStopTransfer
	postData.ComID = 1
	postData.Token = 123456789
	postData.ForceStop = 1

	payload, _ := mqtt.MessagePack(postData)

	mqttMsg := &MqMsg{}
	mqttMsg.Send("47.106.235.93:1883", "57ff69067878495148300967", payload)
	RespondMessage(ctx, CreateErrorMessage(PARAM_ERROR, payload))

}

func StatusQuery(ctx *gin.Context) {
	userID := ctx.MustGet(JwtCtxUidKey)

	if userID.(uint64) <= 0 {
		RespondMessage(ctx, CreateErrorMessage(PARAM_ERROR, nil))
	}

}

func NoLoadSetting(ctx *gin.Context) {
	userID := ctx.MustGet(JwtCtxUidKey)

	if userID.(uint64) <= 0 {
		RespondMessage(ctx, CreateErrorMessage(PARAM_ERROR, nil))
	}
}

func Restart(ctx *gin.Context) {
	userID := ctx.MustGet(JwtCtxUidKey)

	if userID.(uint64) <= 0 {
		RespondMessage(ctx, CreateErrorMessage(PARAM_ERROR, nil))
	}
}

func UpdateFirmware(ctx *gin.Context) {
	userID := ctx.MustGet(JwtCtxUidKey)

	if userID.(uint64) <= 0 {
		RespondMessage(ctx, CreateErrorMessage(PARAM_ERROR, nil))
	}
}
