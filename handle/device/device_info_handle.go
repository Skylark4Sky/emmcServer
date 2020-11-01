package device

import (
	. "GoServer/handle/user"
	. "GoServer/middleWare/dataBases/mysql"
	. "GoServer/utils/respond"
	"strconv"
	"strings"
)

const (
	SELECT_DEVICE_LIST = 110
)

type RequestListData struct {
	UserID    uint64 `fomr:"userID" json:"userID" binding:"required"`
	PageNum   int64  `form:"pageNum" json:"pageNum" binding:"required"`   //起始页
	PageSize  int64  `form:"pageSize" json:"pageSize" binding:"required"` //每页大小
	StartTime int64  `form:"startTime" json:"startTime"`
	EndTime   int64  `form:"endTime" json:"startTime"`
}

func checkUserRules(entity *UserInfo, roleValue int) (isFind bool) {
	isFind = false
	if entity.User.Rules != "" && len(entity.User.Rules) >= 1 {
		countSplit := strings.Split(entity.User.Rules, ",")
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
	userInfo := &UserInfo{}
	err := ExecSQL().Debug().Table("user_base").Select("user_role.rules").Joins("inner join user_role ON user_base.user_role = user_role.id").Where("uid = ?", request.UserID).Scan(&userInfo.User).Error

	if err != nil {
		if IsRecordNotFound(err) {
			return nil, CreateErrorMessage(USER_NO_EXIST, nil)
		}
		return nil, CreateErrorMessage(SYSTEM_ERROR, err)
	}

	if checkUserRules(userInfo,SELECT_DEVICE_LIST) == false {
		return nil, CreateErrorMessage(SYSTEM_ERROR, "没有操作权限")
	}

	return userInfo, nil
}

func (request *RequestListData) GetDeviceTransferLogList() (interface{}, interface{}) {
	return nil, nil
}

func (request *RequestListData) GetModuleList() (interface{}, interface{}) {
	return nil, nil
}

func (request *RequestListData) GetModuleConnectLogList() (interface{}, interface{}) {
	return nil, nil
}
