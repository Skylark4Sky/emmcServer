package device

import (
	. "GoServer/handle/user"
	. "GoServer/middleWare/dataBases/mysql"
	. "GoServer/utils/respond"
	"strconv"
	"strings"
)

const (
	SELECT_DEVICE_LIST = 4
	SELECT_DEVICE_TRANSFER_LOG_LIST = 28
	SELECT_TMODULE_LIST = 12
	SELECT_MODULE_CONNECT_LOG_LIST = 20
)

type RequestListData struct {
	UserID    uint64 `fomr:"userID" json:"userID" binding:"required"`
	PageNum   int64  `form:"pageNum" json:"pageNum" binding:"required"`   //起始页
	PageSize  int64  `form:"pageSize" json:"pageSize" binding:"required"` //每页大小
	StartTime int64  `form:"startTime" json:"startTime"`
	EndTime   int64  `form:"endTime" json:"startTime"`
}

func checkUserRulesGroup(request *RequestListData, roleValue int) (isFind bool, errMsg interface{}) {

	isFind = false
	errMsg = nil

	userInfo := &UserInfo{}
	err := ExecSQL().Table("user_base").Select("user_role.id,user_role.rules").Joins("inner join user_role ON user_base.user_role = user_role.id").Where("uid = ?", request.UserID).Scan(&userInfo.User).Error

	if err != nil {
		if IsRecordNotFound(err) {
			errMsg = CreateErrorMessage(USER_NO_EXIST, nil)
			return
		}
		errMsg = CreateErrorMessage(SYSTEM_ERROR, err)
		return
	}

	if userInfo.User.Rules != "" && len(userInfo.User.Rules) >= 1 {
		countSplit := strings.Split(userInfo.User.Rules, ",")
		for _, ids := range countSplit {
			if role, err := strconv.Atoi(ids); err == nil {
				if role == roleValue {
					isFind = true
					break
				}
			}
		}
	}

	return
}

func (request *RequestListData) GetDeviceList() (interface{}, interface{}) {
	hasRole,errMsg  := checkUserRulesGroup(request,SELECT_DEVICE_LIST)

	if errMsg != nil {
		return nil, CreateErrorMessage(SYSTEM_ERROR, "没有操作权限")
	}

	if hasRole {

	}

	return nil, nil
}

func (request *RequestListData) GetDeviceTransferLogList() (interface{}, interface{}) {
	hasRole,errMsg  := checkUserRulesGroup(request,SELECT_DEVICE_TRANSFER_LOG_LIST)

	if errMsg != nil {
		return nil, CreateErrorMessage(SYSTEM_ERROR, "没有操作权限")
	}

	if hasRole {

	}

	return nil, nil
}

func (request *RequestListData) GetModuleList() (interface{}, interface{}) {
	hasRole,errMsg  := checkUserRulesGroup(request,SELECT_TMODULE_LIST)

	if errMsg != nil {
		return nil, CreateErrorMessage(SYSTEM_ERROR, "没有操作权限")
	}

	if hasRole {

	}

	return nil, nil
}

func (request *RequestListData) GetModuleConnectLogList() (interface{}, interface{}) {
	hasRole,errMsg  := checkUserRulesGroup(request,SELECT_MODULE_CONNECT_LOG_LIST)

	if errMsg != nil {
		return nil, CreateErrorMessage(SYSTEM_ERROR, "没有操作权限")
	}

	if hasRole {

	}

	return nil, nil
}
