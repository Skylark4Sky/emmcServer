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
	SN string `form:"deviceSN" json:"deviceSN" binding:"required"`
}

//设备登记
func DeviceRegister(ctx *gin.Context) {
	var urlParam RequestParam
	if err := ctx.ShouldBindQuery(&urlParam); err != nil {
		RetError(ctx, CreateRetStatus(PARAM_ERROR, err))
		return
	}

	var postData RequestData
	if err := ctx.ShouldBind(&postData); err != nil {
		RetError(ctx, CreateRetStatus(PARAM_ERROR, err))
		return
	}

	ciphertext, err := hex.DecodeString(postData.SN)
	if err != nil {
		RetError(ctx, CreateRetStatus(PARAM_ERROR, err))
		return
	}

	clientID, err := AES_CBCDecrypt(ciphertext, key, iv)

	if err != nil {
		RetError(ctx, CreateRetStatus(PARAM_ERROR, err))
		return
	}

	clientIDStringLen1 := len(clientID)
	clientIDStringLen2 := len(urlParam.ClientID)

	if clientIDStringLen1 < clientIDStringLen2 {
		RetError(ctx, CreateRetStatus(PARAM_ERROR, err))
		RetError(ctx, CreateRetStatus(AUTH_ERROR, "Error ClientID"))
		return
	}

	clientID = string([]byte(clientID)[:clientIDStringLen2])

	if strings.Compare(clientID, urlParam.ClientID) == 0 {
		requestTime := TimeFormat(time.Now())
		requestIP := ctx.ClientIP()
		MqttLog("[", requestIP, "] =========>> ", requestTime, " DeviceConnect ", urlParam.ClientID)
		MqttLog("[", requestIP, "] =========>> ", requestTime, " DeviceInfo ", postData.SN, " ", urlParam.Version)
		ctx.AbortWithStatusJSON(200, gin.H{"status": true, "clientID": urlParam.ClientID, "version": urlParam.Version, "deviceSN": postData.SN})
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": false, "error": errors.New("Error clientID")})
	}
}
