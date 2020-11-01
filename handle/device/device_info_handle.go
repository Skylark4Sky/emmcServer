package device

import (
	. "GoServer/handle/user"
	. "GoServer/middleWare/dataBases/mysql"
	. "GoServer/utils/respond"
)

type RequestListData struct {
	UserID    uint64 `fomr:"userID" json:"userID" binding:"required"`
	PageNum   int64  `form:"pageNum" json:"pageNum" binding:"required"`   //起始页
	PageSize  int64  `form:"pageSize" json:"pageSize" binding:"required"` //每页大小
	StartTime int64  `form:"startTime" json:"startTime"`
	EndTime   int64  `form:"endTime" json:"startTime"`
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
