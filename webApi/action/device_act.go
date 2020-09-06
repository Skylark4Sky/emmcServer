package action

import (
	. "GoServer/utils"
	. "GoServer/webApi/utils"
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
func DeviceRegister(context *gin.Context) {
	var urlParam RequestParam
	if err := context.ShouldBindQuery(&urlParam); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var postData RequestData
	if err := context.ShouldBind(&postData); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ciphertext, err := hex.DecodeString(postData.SN)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	clientID, err := AES_CBCDecrypt(ciphertext, key, iv)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	clientIDStringLen1 := len(clientID)
	clientIDStringLen2 := len(urlParam.ClientID)

	if clientIDStringLen1 < clientIDStringLen2 {
		context.JSON(http.StatusBadRequest, gin.H{"error": errors.New("Error clientID")})
		return
	}

	clientID = string([]byte(clientID)[:clientIDStringLen2])

	if strings.Compare(clientID, urlParam.ClientID) == 0 {
		requestTime := time.Now().Format(GetSystem().Timeformat)
		requestIP := context.ClientIP()
		PrintInfo("[", requestIP, "] =========>> ", requestTime, " DeviceConnect ", urlParam.ClientID)
		PrintInfo("[", requestIP, "] =========>> ", requestTime, " DeviceInfo ", postData.SN, " ", urlParam.Version)
		context.AbortWithStatusJSON(200, gin.H{"status": true, "clientID": urlParam.ClientID, "version": urlParam.Version, "deviceSN": postData.SN})
	} else {
		context.JSON(http.StatusBadRequest, gin.H{"status": false, "error": errors.New("Error clientID")})
	}
}
