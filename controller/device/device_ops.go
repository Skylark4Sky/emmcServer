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

func checkRequestUserID(ctx *gin.Context, requestUserID uint64) (ok bool, userID uint64, errMsg interface{}) {
	userID = ctx.MustGet(JwtCtxUidKey).(uint64)
	ok = true
	if userID <= 0 {
		ok = false
		errMsg = CreateErrorMessage(PARAM_ERROR, nil)
	}

	return
}

func checkRequestParam(ctx *gin.Context, requestParam *RequestListData, minSize int64, maxSize int64) (bool, interface{}) {

	if err := ctx.ShouldBind(&requestParam); err != nil {
		return false, CreateErrorMessage(PARAM_ERROR, err)
	}

	ok, userID, errMsg := checkRequestUserID(ctx, requestParam.UserID)
	if !ok {
		return ok, errMsg
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

	if errMsg := CheckUserRulesGroup(getListData.UserID, SELECT_DEVICE_LIST); errMsg != nil {
		RespondMessage(ctx, errMsg)
		return
	}

	data, err := getListData.GetDeviceList(getListData.UserID)

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

	if errMsg := CheckUserRulesGroup(getListData.UserID, SELECT_DEVICE_TRANSFER_LOG_LIST); errMsg != nil {
		RespondMessage(ctx, errMsg)
		return
	}

	data, err := getListData.GetDeviceTransferLogList(getListData.UserID)

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

	if errMsg := CheckUserRulesGroup(getListData.UserID, SELECT_TMODULE_LIST); errMsg != nil {
		RespondMessage(ctx, errMsg)
		return
	}

	data, err := getListData.GetModuleList(getListData.UserID)

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

	if errMsg := CheckUserRulesGroup(getListData.UserID, SELECT_MODULE_CONNECT_LOG_LIST); errMsg != nil {
		RespondMessage(ctx, errMsg)
		return
	}

	data, err := getListData.GetModuleConnectLogList(getListData.UserID)

	if err != nil {
		RespondMessage(ctx, err)
		return
	}

	RespondMessage(ctx, CreateMessage(SUCCESS, data))
}

func SyncDeviceStatus(ctx *gin.Context) {
	var sync RequestSyncData

	if err := ctx.ShouldBind(&sync); err != nil {
		CreateErrorMessage(PARAM_ERROR, err)
		return
	}

	ok, _, errMsg := checkRequestUserID(ctx, sync.UserID)
	if !ok {
		RespondMessage(ctx, errMsg)
		return
	}

	if errMsg := CheckUserRulesGroup(sync.UserID, SYNC_DEVICE_STATUS); errMsg != nil {
		RespondMessage(ctx, errMsg)
		return
	}

	data, err := sync.SyncDeviceStatus()

	if err != nil {
		RespondMessage(ctx, err)
		return
	}
	RespondMessage(ctx, CreateMessage(SUCCESS, data))
}
