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
	//	var urlParam RequestParam
	//	if err := ctx.ShouldBindQuery(&urlParam); err != nil {
	//		RespondMessage(ctx, CreateErrorMessage(PARAM_ERROR, "参数错误1"))
	//		return
	//	}
	var postData RequestData
	if err := ctx.ShouldBind(&postData); err != nil {
		RespondMessage(ctx, CreateErrorMessage(PARAM_ERROR, "参数错误1"))
		return
	}

	ciphertext, err := hex.DecodeString(postData.Token)

	if err != nil {
		RespondMessage(ctx, CreateErrorMessage(PARAM_ERROR, "参数错误2"))
		return
	}

	ModuleSN, err := AES_CBCDecrypt(ciphertext, key, iv)

	if err != nil {
		RespondMessage(ctx, CreateErrorMessage(PARAM_ERROR, "参数错误3"))
		return
	}

	SNStringLen1 := len(ModuleSN)
	SNStringLen2 := len(postData.ModuleSN)

	if SNStringLen1 < SNStringLen2 {
		RespondMessage(ctx, CreateErrorMessage(AUTH_ERROR, "Error Token"))
		return
	}

	ModuleSN = string([]byte(ModuleSN)[:SNStringLen2])

	if strings.Compare(ModuleSN, postData.ModuleSN) == 0 {
		respond := postData.Connect(ctx)
		RespondMessage(ctx, respond)
	} else {
		RespondMessage(ctx, CreateErrorMessage(AUTH_ERROR, "认证失败"))
	}
}

func DeviceList(ctx *gin.Context) {

}
