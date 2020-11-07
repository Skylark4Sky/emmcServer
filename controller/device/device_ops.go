package device

import (
	. "GoServer/handle/device"
	. "GoServer/middleWare/extension"
	. "GoServer/utils/respond"
	. "GoServer/utils/string"
	"github.com/gin-gonic/gin"
)

const (
	MIN_PAGE_SIZE = 10
	MAX_PAGE_SIZE = 100
)

func checkRequestParam(ctx *gin.Context, requestParam *RequestListData, minSize int64, maxSize int64) (bool, interface{}) {
	userID := ctx.MustGet(JwtCtxUidKey)
	if userID.(uint64) <= 0 {
		return false, CreateErrorMessage(PARAM_ERROR, nil)
	}
	if err := ctx.ShouldBind(&requestParam); err != nil {
		return false, CreateErrorMessage(PARAM_ERROR, err)
	}

	if requestParam.PageNum <= 0 {
		return false, CreateErrorMessage(PARAM_ERROR, "起始页不能小于1")
	}

	if requestParam.PageSize < minSize || requestParam.PageSize > maxSize {
		errMsg := StringJoin([]interface{}{"页大小设置错误 ", MIN_PAGE_SIZE, " - ", MAX_PAGE_SIZE})
		return false, CreateErrorMessage(PARAM_ERROR, errMsg)
	}

	if requestParam.UserID != userID {
		return false, CreateErrorMessage(PARAM_ERROR, nil)
	}

	return true, nil
}

// 返回设备列表
func GetDeviceList(ctx *gin.Context) {
	var getListData RequestListData
	if _, err := checkRequestParam(ctx, &getListData, MIN_PAGE_SIZE, MAX_PAGE_SIZE); err != nil {
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
	if _, err := checkRequestParam(ctx, &getListData, MIN_PAGE_SIZE, MAX_PAGE_SIZE); err != nil {
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
	if _, err := checkRequestParam(ctx, &getListData, MIN_PAGE_SIZE, MAX_PAGE_SIZE); err != nil {
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
	if _, err := checkRequestParam(ctx, &getListData, MIN_PAGE_SIZE, MAX_PAGE_SIZE); err != nil {
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
