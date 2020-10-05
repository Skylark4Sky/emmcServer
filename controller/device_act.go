package action

import (
	. "GoServer/handle/device"
	//	. "GoServer/utils/log"
	. "GoServer/utils/respond"
	. "GoServer/utils/security"
	//	. "GoServer/utils/time"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"strings"
	//	"time"
)

var (
	key []byte = []byte("78hrey23y28ogs89")
	iv  []byte = []byte("1234567890123456")
)

//设备登记
func DeviceConnect(ctx *gin.Context) {
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

	ciphertext, err := hex.DecodeString(postData.ModuleSN)

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
		respond := postData.Connect(ctx, urlParam.ClientID)
		RespondMessage(ctx, respond)
	} else {
		RespondMessage(ctx, CreateErrorMessage(AUTH_ERROR, "认证失败"))
	}
}

func DeviceList(ctx *gin.Context) {

}
