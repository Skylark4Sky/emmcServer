package device

import (
	. "GoServer/handle/device"
	. "GoServer/middleWare/extension"
	. "GoServer/utils/respond"
	"github.com/gin-gonic/gin"
)

func checkRequestParam(ctx *gin.Context, requestParam *RequestListData) (bool, interface{}) {
	userID := ctx.MustGet(JwtCtxUidKey)
	if userID.(uint64) <= 0 {
		return false, CreateErrorMessage(PARAM_ERROR, nil)
	}
	if err := ctx.ShouldBind(&requestParam); err != nil {
		return false, CreateErrorMessage(PARAM_ERROR, err)
	}
	if requestParam.UserID != userID {
		return false, CreateErrorMessage(PARAM_ERROR, nil)
	}

	return true, nil
}

// 返回设备列表
func GetDeviceList(ctx *gin.Context) {
	var getListData RequestListData
	if _, err := checkRequestParam(ctx, &getListData); err != nil {
		RespondMessage(ctx, err)
		return
	}

	data, err := getListData.GetDeviceList()

	if err != nil {
		RespondMessage(ctx, err)
		return
	}

	RespondMessage(ctx, CreateMessage(SUCCESS, data))
}

// 返回设备上报日志
func GetDeviceTransferLogList(ctx *gin.Context) {
	var getListData RequestListData
	if _, err := checkRequestParam(ctx, &getListData); err != nil {
		RespondMessage(ctx, err)
		return
	}

	data, err := getListData.GetDeviceTransferLogList()

	if err != nil {
		RespondMessage(ctx, err)
		return
	}

	RespondMessage(ctx, CreateMessage(SUCCESS, data))
}

func GetModuleList(ctx *gin.Context) {
	var getListData RequestListData
	if _, err := checkRequestParam(ctx, &getListData); err != nil {
		RespondMessage(ctx, err)
		return
	}

	data, err := getListData.GetModuleList()

	if err != nil {
		RespondMessage(ctx, err)
		return
	}

	RespondMessage(ctx, CreateMessage(SUCCESS, data))
}

func GetModuleConnectLogList(ctx *gin.Context) {
	var getListData RequestListData
	if _, err := checkRequestParam(ctx, &getListData); err != nil {
		RespondMessage(ctx, err)
		return
	}

	data, err := getListData.GetModuleConnectLogList()

	if err != nil {
		RespondMessage(ctx, err)
		return
	}

	RespondMessage(ctx, CreateMessage(SUCCESS, data))
}
