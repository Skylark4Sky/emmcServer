package action

import (
	. "GoServer/utils"
	"encoding/hex"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

var (
	key []byte = []byte("78hrey23y28ogs89")
	iv  []byte = []byte("1234567890123456")
)

type RequestParam struct {
	ClientID string `form:"clientID"`
	Version  string `form:"version"`
}

type RequestData struct {
	DeviceNo string `form:"deviceNo" json:"deviceNo" binding:"required"`
	Token    string `form:"token" json:"token" binding:"required"`
}

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
		RespondMessage(ctx, CreateErrorMessage(PARAM_ERROR, err))
		return
	}

	clientIDStringLen1 := len(clientID)
	clientIDStringLen2 := len(urlParam.ClientID)

	if clientIDStringLen1 < clientIDStringLen2 {
		RespondMessage(ctx, CreateErrorMessage(PARAM_ERROR, err))
		RespondMessage(ctx, CreateErrorMessage(AUTH_ERROR, "Error Token"))
		return
	}

	clientID = string([]byte(clientID)[:clientIDStringLen2])

	if strings.Compare(clientID, urlParam.ClientID) == 0 {
		requestTime := TimeFormat(time.Now())
		requestIP := ctx.ClientIP()
		MqttLog("[", requestIP, "] =========>> ", requestTime, " DeviceConnect ", urlParam.ClientID)
		MqttLog("[", requestIP, "] =========>> ", requestTime, " DeviceInfo ", postData.Token, " ", urlParam.Version)

		data := map[string]string{
			"url":  "http://www.gisunlink.com/GiSunLink.v2_to_v3.ota.bin",
			"size": "527300",
		}

		RespondMessage(ctx, CreateMessage(SUCCESS, data))
	} else {
		RespondMessage(ctx, CreateErrorMessage(PARAM_ERROR, "认证失败"))
	}
}
