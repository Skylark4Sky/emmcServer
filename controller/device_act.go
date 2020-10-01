package action

import (
	. "GoServer/handle/device"
	. "GoServer/utils/log"
	. "GoServer/utils/respond"
	. "GoServer/utils/security"
	. "GoServer/utils/time"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

var (
	key []byte = []byte("78hrey23y28ogs89")
	iv  []byte = []byte("1234567890123456")
)

//设备登记
func DeviceRegister(ctx *gin.Context) {
	var urlParam RequestParam
	if err := ctx.ShouldBindQuery(&urlParam); err != nil {
		RespondMessage(ctx, CreateErrorMessage(PARAM_ERROR, "参数错误"))
		return
	}

	var postData RequestData
	if err := ctx.ShouldBind(&postData); err != nil {
		RespondMessage(ctx, CreateErrorMessage(PARAM_ERROR, "参数错误"))
		return
	}

	ciphertext, err := hex.DecodeString(postData.Token)

	if err != nil {
		RespondMessage(ctx, CreateErrorMessage(PARAM_ERROR, "参数错误"))
		return
	}

	clientID, err := AES_CBCDecrypt(ciphertext, key, iv)

	if err != nil {
		RespondMessage(ctx, CreateErrorMessage(PARAM_ERROR, "参数错误"))
		return
	}

	clientIDStringLen1 := len(clientID)
	clientIDStringLen2 := len(urlParam.ClientID)

	if clientIDStringLen1 < clientIDStringLen2 {
		RespondMessage(ctx, CreateErrorMessage(AUTH_ERROR, "Error Token"))
		return
	}

	clientID = string([]byte(clientID)[:clientIDStringLen2])

	if strings.Compare(clientID, urlParam.ClientID) == 0 {
		requestTime := TimeFormat(time.Now())
		requestIP := ctx.ClientIP()
		MqttLog("[", requestIP, "] =========>> ", requestTime, " DeviceConnect ", urlParam.ClientID)
		MqttLog("[", requestIP, "] =========>> ", requestTime, " DeviceInfo ", postData.Token, " ", urlParam.Version)
		//	data := &FirmwareInfo{
		//		URL:  "http://www.gisunlink.com/GiSunLink.ota.bin",
		//		Size: 476448,
		//	}

		RespondMessage(ctx, CreateMessage(SUCCESS, nil))
	} else {
		RespondMessage(ctx, CreateErrorMessage(AUTH_ERROR, "认证失败"))
	}
}
