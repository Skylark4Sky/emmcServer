package device

import (
	. "GoServer/handle/user"
	. "GoServer/middleWare/dataBases/mysql"
	. "GoServer/model/device"
	. "GoServer/utils/respond"
	. "GoServer/utils/string"
	"github.com/jinzhu/gorm"
	"strconv"
	"strings"
)

const StartPage = 1

const (
	SELECT_DEVICE_LIST              = 4
	SELECT_DEVICE_TRANSFER_LOG_LIST = 28
	SELECT_TMODULE_LIST             = 12
	SELECT_MODULE_CONNECT_LOG_LIST  = 20
)

const (
	SELECT_TRANSFER_LOG_DEVICEID     = "device_id"
	SELECT_TRANSFER_LOG_BEHAVIOR     = "behavior"
	SELECT_TRANSFER_LOG_DEVICESN     = "device_sn"
	SELECT_TRANSFER_LOG_TRANSFERID   = "transfer_id"
	SELECT_TRANSFER_LOG_TRANSFERTIME = "transfer_time"
	SELECT_TRANSFER_LOG_CREATETIME   = "create_time"
)

const (
	SELECT_DEVICE_LIST_TYPE            = "type"
	SELECT_DEVICE_LIST_DEVICE_SN       = "device_sn"
	SELECT_DEVICE_LIST_DEVICE_VERSION  = "device_version"
	SELECT_DEVICE_LIST_TYPE_ACCESS_WAY = "access_way"
	SELECT_DEVICE_LIST_TIMETYPE        = "time"
)

type RequestListData struct {
	UserID      uint64      `fomr:"userID" json:"userID" binding:"required"`
	PageNum     int64       `form:"pageNum" json:"pageNum" binding:"required"`   //起始页
	PageSize    int64       `form:"pageSize" json:"pageSize" binding:"required"` //每页大小
	StartTime   int64       `form:"startTime" json:"startTime"`
	EndTime     int64       `form:"endTime" json:"endTime"`
	RequestCond interface{} `form:"requestCond" json:"requestCond"`
}

type PageInfo struct {
	Total int64 `json:"total,omitempty"`
	Size  int   `json:"size,omitempty"`
}

type RespondListData struct {
	List interface{} `json:"list"`
	Page PageInfo    `json:"page,omitempty"`
}

func checkUserRulesGroup(request *RequestListData, roleValue int) (errMsg interface{}) {
	errMsg = nil
	userInfo := &UserInfo{}

	db := ExecSQL().Table("user_base")
	db = db.Select("user_role.id,user_role.rules")
	db = db.Joins("inner join user_role ON user_base.user_role = user_role.id")
	db = db.Where("uid = ?", request.UserID)

	if err := db.Scan(&userInfo.User).Error; err != nil {
		if IsRecordNotFound(err) {
			errMsg = CreateErrorMessage(USER_NO_EXIST, nil)
			return
		}
		errMsg = CreateErrorMessage(SYSTEM_ERROR, err)
		return
	}

	var isFind = false

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

	if isFind == false {
		errMsg = CreateErrorMessage(SYSTEM_ERROR, "没有操作权限")
	}

	return nil
}

func addTimeCond(db *gorm.DB, timeField string, startTime, endTime int64) *gorm.DB {
	dbEntity := db
	if startTime > 0 {
		cond := StringJoin([]interface{}{" ", timeField, " >= ?"})
		dbEntity = dbEntity.Where(cond, startTime*1000)
	}

	if endTime > 0 {
		cond := StringJoin([]interface{}{" ", timeField, " <= ?"})
		dbEntity = dbEntity.Where(cond, endTime*1000)
	}
	return dbEntity
}

func (request *RequestListData) GetDeviceList() (*RespondListData, interface{}) {
	errMsg := checkUserRulesGroup(request, SELECT_DEVICE_LIST)

	if errMsg != nil {
		return nil, errMsg
	}

	var deviceList []DeviceInfo
	var total int64 = 0

	db := ExecSQL().Debug()

	db = db.Limit(request.PageSize).Offset((request.PageNum - 1) * request.PageSize).Order("id desc")

	if request.RequestCond != nil {
		condMap := request.RequestCond.(map[string]interface{})

		for keyName, condValue := range condMap {
			fmt.Println("keyName:", keyName)
			cond := StringJoin([]interface{}{" ", keyName, " = ?"})
			switch keyName {
			case SELECT_DEVICE_LIST_TYPE, SELECT_DEVICE_LIST_DEVICE_SN, SELECT_DEVICE_LIST_DEVICE_VERSION, SELECT_DEVICE_LIST_TYPE_ACCESS_WAY:
				{
					db = db.Where(cond, condValue)
				}
			case SELECT_DEVICE_LIST_TIMETYPE:
				{
					db = addTimeCond(db, condValue.(string), request.StartTime, request.EndTime)
				}
			}
		}
	} else {
		db = addTimeCond(db, "create_time", request.StartTime, request.EndTime)
	}

	if request.PageNum == StartPage {
		if err := db.Find(&deviceList).Count(&total).Error; err != nil {
			return nil, CreateErrorMessage(SYSTEM_ERROR, err)
		}
	} else {
		if err := db.Find(&deviceList).Error; err != nil {
			return nil, CreateErrorMessage(SYSTEM_ERROR, err)
		}
	}

	var respond = RespondListData{
		List: deviceList,
		Page: PageInfo{
			Total: total,
			Size:  len(deviceList),
		},
	}
	return &respond, nil
}

func (request *RequestListData) GetDeviceTransferLogList() (interface{}, interface{}) {
	errMsg := checkUserRulesGroup(request, SELECT_DEVICE_TRANSFER_LOG_LIST)

	if errMsg != nil {
		return nil, CreateErrorMessage(SYSTEM_ERROR, "没有操作权限")
	}

	var transferList []DeviceTransferLog
	var total int64 = 0

	db := ExecSQL().Debug()

	db = db.Limit(request.PageSize).Offset((request.PageNum - 1) * request.PageSize).Order("id desc")

	if request.RequestCond != nil {
		condMap := request.RequestCond.(map[string]interface{})

		for keyName, condValue := range condMap {
			cond := StringJoin([]interface{}{" ", keyName, " = ?"})
			switch keyName {
			case SELECT_DEVICE_LIST_TYPE, SELECT_DEVICE_LIST_DEVICE_SN, SELECT_DEVICE_LIST_DEVICE_VERSION, SELECT_DEVICE_LIST_TYPE_ACCESS_WAY:
				{
					db = db.Where(cond, condValue)
				}
			case SELECT_DEVICE_LIST_TIMETYPE:
				{
					db = addTimeCond(db, condValue.(string), request.StartTime, request.EndTime)
				}
			}
		}
	} else {
		db = addTimeCond(db, "create_time", request.StartTime, request.EndTime)
	}

	if request.PageNum == StartPage {
		if err := db.Find(&transferList).Count(&total).Error; err != nil {
			return nil, CreateErrorMessage(SYSTEM_ERROR, err)
		}
	} else {
		if err := db.Find(&deviceList).Error; err != nil {
			return nil, CreateErrorMessage(SYSTEM_ERROR, err)
		}
	}

	var respond = RespondListData{
		List: deviceList,
		Page: PageInfo{
			Total: total,
			Size:  len(deviceList),
		},
	}

	return &respond, nil
}

func (request *RequestListData) GetModuleList() (interface{}, interface{}) {
	errMsg := checkUserRulesGroup(request, SELECT_TMODULE_LIST)

	if errMsg != nil {
		return nil, CreateErrorMessage(SYSTEM_ERROR, "没有操作权限")
	}

	var deviceList []DeviceInfo
	var total int64 = 0

	db := ExecSQL().Debug()

	db = db.Limit(request.PageSize).Offset((request.PageNum - 1) * request.PageSize).Order("id desc")

	if request.RequestCond != nil {
		condMap := request.RequestCond.(map[string]interface{})

		for keyName, condValue := range condMap {
			fmt.Println("keyName:", keyName)
			cond := StringJoin([]interface{}{" ", keyName, " = ?"})
			switch keyName {
			case SELECT_DEVICE_LIST_TYPE, SELECT_DEVICE_LIST_DEVICE_SN, SELECT_DEVICE_LIST_DEVICE_VERSION, SELECT_DEVICE_LIST_TYPE_ACCESS_WAY:
				{
					db = db.Where(cond, condValue)
				}
			case SELECT_DEVICE_LIST_TIMETYPE:
				{
					db = addTimeCond(db, condValue.(string), request.StartTime, request.EndTime)
				}
			}
		}
	} else {
		db = addTimeCond(db, "create_time", request.StartTime, request.EndTime)
	}

	if request.PageNum == StartPage {
		if err := db.Find(&deviceList).Count(&total).Error; err != nil {
			return nil, CreateErrorMessage(SYSTEM_ERROR, err)
		}
	} else {
		if err := db.Find(&deviceList).Error; err != nil {
			return nil, CreateErrorMessage(SYSTEM_ERROR, err)
		}
	}

	var respond = RespondListData{
		List: deviceList,
		Page: PageInfo{
			Total: total,
			Size:  len(deviceList),
		},
	}

	return &respond, nil
}

func (request *RequestListData) GetModuleConnectLogList() (interface{}, interface{}) {
	errMsg := checkUserRulesGroup(request, SELECT_MODULE_CONNECT_LOG_LIST)

	if errMsg != nil {
		return nil, CreateErrorMessage(SYSTEM_ERROR, "没有操作权限")
	}

	var deviceList []DeviceInfo
	var total int64 = 0

	db := ExecSQL().Debug()

	db = db.Limit(request.PageSize).Offset((request.PageNum - 1) * request.PageSize).Order("id desc")

	if request.RequestCond != nil {
		condMap := request.RequestCond.(map[string]interface{})

		for keyName, condValue := range condMap {
			fmt.Println("keyName:", keyName)
			cond := StringJoin([]interface{}{" ", keyName, " = ?"})
			switch keyName {
			case SELECT_DEVICE_LIST_TYPE, SELECT_DEVICE_LIST_DEVICE_SN, SELECT_DEVICE_LIST_DEVICE_VERSION, SELECT_DEVICE_LIST_TYPE_ACCESS_WAY:
				{
					db = db.Where(cond, condValue)
				}
			case SELECT_DEVICE_LIST_TIMETYPE:
				{
					db = addTimeCond(db, condValue.(string), request.StartTime, request.EndTime)
				}
			}
		}
	} else {
		db = addTimeCond(db, "create_time", request.StartTime, request.EndTime)
	}

	if request.PageNum == StartPage {
		if err := db.Find(&deviceList).Count(&total).Error; err != nil {
			return nil, CreateErrorMessage(SYSTEM_ERROR, err)
		}
	} else {
		if err := db.Find(&deviceList).Error; err != nil {
			return nil, CreateErrorMessage(SYSTEM_ERROR, err)
		}
	}

	var respond = RespondListData{
		List: deviceList,
		Page: PageInfo{
			Total: total,
			Size:  len(deviceList),
		},
	}

	return &respond, nil
}
