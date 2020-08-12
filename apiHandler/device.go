package apiHandler

import (
	. "GoServer/utils"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
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

func AES_CBCDecrypt(cipherText []byte, key []byte, iv []byte) (data string, err error) {
	//指定解密算法，返回一个AES算法的Block接口对象
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	blockMode := cipher.NewCBCDecrypter(block, iv)
	plainText := make([]byte, len(cipherText))
	blockMode.CryptBlocks(plainText, cipherText)
	return string(plainText), nil
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
		PrintInfo("DeviceRegister ", urlParam.ClientID, " ", postData.SN, " ", urlParam.Version)
		context.AbortWithStatusJSON(200, gin.H{"status": true, "clientID": urlParam.ClientID, "version": urlParam.Version, "deviceSN": postData.SN})
	} else {
		context.JSON(http.StatusBadRequest, gin.H{"status": false, "error": errors.New("Error clientID")})
	}
}
