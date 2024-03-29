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

// Jwt验证
func checkRequestUserID(ctx *gin.Context, requestUserID uint64) (ok bool, userID uint64, errMsg interface{}) {
	userID = ctx.MustGet(JwtCtxUidKey).(uint64)
	ok = true
	if userID <= 0 {
		ok = false
		errMsg = CreateErrorMessage(PARAM_ERROR, nil)
	}
	return
}

// 检查请求参数
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

// 通用请求参数检查
func generalRequestChk(ctx *gin.Context, request *RequestListData) (bool, interface{}) {
	if _, err := checkRequestParam(ctx, request, MIN_PAGE_SIZE, MAX_PAGE_SIZE); err != nil {
		RespondMessage(ctx, err)
		return true, err
	}

	if err := CheckUserRulesGroup(request.UserID, SELECT_DEVICE_LIST); err != nil {
		return true, err
	}
	return false, nil
}

// 返回设备端口列表
func GetDeviceCom(ctx *gin.Context) {

}

// 返回设备充电列表
func GetDeviceChargeList(ctx *gin.Context) {
	var getListData RequestListData

	if _, errMsg := generalRequestChk(ctx, &getListData); errMsg != nil {
		RespondMessage(ctx, errMsg)
	}

	data, err := getListData.GetDeviceChargeList(getListData.UserID)

	if err != nil {
		RespondMessage(ctx, err)
		return
	}

	RespondMessage(ctx, CreateMessage(SUCCESS, data))
}

// 返回设备列表
func GetDeviceList(ctx *gin.Context) {
	var getListData RequestListData
	if _, errMsg := generalRequestChk(ctx, &getListData); errMsg != nil {
		RespondMessage(ctx, errMsg)
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

// 返回设备模块列表
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

// 返回设备连接日志
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

// 同步设备状态
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
