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

const (
	SELECT_DEVICE_LIST              = 4
	SELECT_DEVICE_TRANSFER_LOG_LIST = 28
	SELECT_TMODULE_LIST             = 12
	SELECT_MODULE_CONNECT_LOG_LIST  = 20
)

const (
	SELECT_CREATE_TIME = "create_time"
	SELECT_UPDATE_TIME = "update_time"
)

const (
	SELECT_DEVICE_SN         = "device_sn"
	SELECT_DEVICE_VERSION    = "device_version"
	SELECT_DEVICE_TYPE       = "type"
	SELECT_DEVICE_ACCESS_WAY = "access_way"
	SELECT_DEVICE_TIMETYPE   = "time"
)

const (
	SELECT_TRANSFER_ID       = "transfer_id"
	SELECT_TRANSFER_DEVICEID = "device_id"
	SELECT_TRANSFER_DEVICESN = "device_sn"
	SELECT_TRANSFER_BEHAVIOR = "behavior"
	SELECT_TRANSFER_TIMETYPE = "time"
)

const (
	SELECT_MODULE_SN         = "module_sn"
	SELECT_MODULE_VERSION    = "module_version"
	SELECT_MODULE_ACCESS_WAY = "access_way"
	SELECT_MODULE_TIMETYPE   = "time"
)

const (
	SELECT_CONNECT_MODULE_ID  = "module_id"
	SELECT_CONNECT_MODULE_SN  = "module_sn"
	SELECT_CONNECT_ACCESS_WAY = "access_way"
	SELECT_CONNECT_TIMETYPE   = "time"
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
	Total int64 `json:"total"`
	Size  int   `json:"size"`
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

	if startTime > 0 && endTime > 0 {
		cond := StringJoin([]interface{}{" (", timeField, " BETWEEN ? AND ?)"})
		dbEntity = dbEntity.Where(cond, startTime*1000, endTime*1000)
	} else {
		if startTime > 0 {
			cond := StringJoin([]interface{}{" ", timeField, " >= ?"})
			dbEntity = dbEntity.Where(cond, startTime*1000)
		}

		if endTime > 0 {
			cond := StringJoin([]interface{}{" ", timeField, " <= ?"})
			dbEntity = dbEntity.Where(cond, endTime*1000)
		}
	}

	return dbEntity
}

func generalSQLFormat(request *RequestListData, listSearch interface{}, condFilter func(db *gorm.DB, condMap map[string]interface{}) *gorm.DB) (errMsg interface{}, respond *RespondListData) {
	errMsg = nil
	respond = nil

	var total int64 = 0

	db := ExecSQL().Debug()

	if request.RequestCond != nil {
		condMap := request.RequestCond.(map[string]interface{})
		if condFilter != nil {
			db = condFilter(db, condMap)
		}
	} else {
		db = addTimeCond(db, "create_time", request.StartTime, request.EndTime)
	}

	totalRows := db.NewScope(listSearch).DB()

	db = db.Order("id desc").Limit(request.PageSize).Offset((request.PageNum - 1) * request.PageSize)

	if err := db.Find(listSearch).Error; err != nil {
		errMsg = CreateErrorMessage(SYSTEM_ERROR, err)
		return
	}

	if err := totalRows.Count(&total).Error; err != nil {
		errMsg = CreateErrorMessage(SYSTEM_ERROR, err)
		return
	}

	var listSize int = 0

	switch listSearch.(type) {
	case *[]DeviceInfo:
		{
			listSize = len(*listSearch.(*[]DeviceInfo))
		}
	case *[]DeviceTransferLog:
		{
			listSize = len(*listSearch.(*[]DeviceTransferLog))
		}
	case *[]ModuleInfo:
		{
			listSize = len(*listSearch.(*[]ModuleInfo))
		}
	case *[]ModuleConnectLog:
		{
			listSize = len(*listSearch.(*[]ModuleConnectLog))
		}
	}

	respond = &RespondListData{
		List: listSearch,
		Page: PageInfo{
			Total: total,
			Size:  listSize,
		},
	}

	return
}

func (request *RequestListData) GetDeviceList() (*RespondListData, interface{}) {
	errMsg := checkUserRulesGroup(request, SELECT_DEVICE_LIST)

	if errMsg != nil {
		return nil, errMsg
	}

	var deviceList []DeviceInfo
	var respond *RespondListData = nil

	if errMsg, respond = generalSQLFormat(request, &deviceList, func(db *gorm.DB, condMap map[string]interface{}) *gorm.DB {
		dbEntity := db
		for keyName, condValue := range condMap {
			cond := StringJoin([]interface{}{" ", keyName, " = ?"})
			switch keyName {
			case SELECT_DEVICE_TYPE, SELECT_DEVICE_SN, SELECT_DEVICE_VERSION, SELECT_DEVICE_ACCESS_WAY:
				{
					dbEntity = dbEntity.Where(cond, condValue)
				}
			case SELECT_DEVICE_TIMETYPE:
				{
					dbEntity = addTimeCond(dbEntity, condValue.(string), request.StartTime, request.EndTime)
				}
			}
		}
		return dbEntity
	}); errMsg != nil {
		return nil, errMsg
	}

	return respond, nil
}

func (request *RequestListData) GetDeviceTransferLogList() (interface{}, interface{}) {
	errMsg := checkUserRulesGroup(request, SELECT_DEVICE_TRANSFER_LOG_LIST)

	if errMsg != nil {
		return nil, CreateErrorMessage(SYSTEM_ERROR, "没有操作权限")
	}

	var transferList []DeviceTransferLog
	var respond *RespondListData = nil

	if errMsg, respond = generalSQLFormat(request, &transferList, func(db *gorm.DB, condMap map[string]interface{}) *gorm.DB {
		dbEntity := db
		
		if transferID, ok := condMap[SELECT_TRANSFER_ID]; ok {
			cond := StringJoin([]interface{}{" ", SELECT_TRANSFER_ID, " = ?"})
			dbEntity = dbEntity.Where(cond, transferID)
		}

		if deviceID, ok := condMap[SELECT_TRANSFER_DEVICEID]; ok {
			cond := StringJoin([]interface{}{" ", SELECT_TRANSFER_DEVICEID, " = ?"})
			dbEntity = dbEntity.Where(cond, deviceID)
		}

		if deviceSN, ok := condMap[SELECT_TRANSFER_DEVICESN]; ok {
			cond := StringJoin([]interface{}{" ", SELECT_TRANSFER_DEVICESN, " = ?"})
			dbEntity = dbEntity.Where(cond, deviceSN)
		}

		if behavior, ok := condMap[SELECT_TRANSFER_BEHAVIOR]; ok {
			cond := StringJoin([]interface{}{" ", SELECT_TRANSFER_BEHAVIOR, " = ?"})
			dbEntity = dbEntity.Where(cond, behavior)
		}

		if timeType, ok := condMap[SELECT_TRANSFER_TIMETYPE]; ok {
			dbEntity = addTimeCond(dbEntity, timeType.(string), request.StartTime, request.EndTime)
		}

		return dbEntity
	}); errMsg != nil {
		return nil, errMsg
	}

	return respond, nil
}

func (request *RequestListData) GetModuleList() (interface{}, interface{}) {
	errMsg := checkUserRulesGroup(request, SELECT_TMODULE_LIST)

	if errMsg != nil {
		return nil, CreateErrorMessage(SYSTEM_ERROR, "没有操作权限")
	}

	var moduleList []ModuleInfo
	var respond *RespondListData = nil

	if errMsg, respond = generalSQLFormat(request, &moduleList, func(db *gorm.DB, condMap map[string]interface{}) *gorm.DB {
		dbEntity := db
		for keyName, condValue := range condMap {
			cond := StringJoin([]interface{}{" ", keyName, " = ?"})
			switch keyName {
			case SELECT_MODULE_ACCESS_WAY, SELECT_MODULE_SN, SELECT_MODULE_VERSION:
				{
					dbEntity = dbEntity.Where(cond, condValue)
				}
			case SELECT_MODULE_TIMETYPE:
				{
					dbEntity = addTimeCond(dbEntity, condValue.(string), request.StartTime, request.EndTime)
				}
			}
		}
		return dbEntity
	}); errMsg != nil {
		return nil, errMsg
	}

	return respond, nil
}

func (request *RequestListData) GetModuleConnectLogList() (interface{}, interface{}) {
	errMsg := checkUserRulesGroup(request, SELECT_MODULE_CONNECT_LOG_LIST)

	if errMsg != nil {
		return nil, CreateErrorMessage(SYSTEM_ERROR, "没有操作权限")
	}

	var connectList []ModuleConnectLog
	var respond *RespondListData = nil

	if errMsg, respond = generalSQLFormat(request, &connectList, func(db *gorm.DB, condMap map[string]interface{}) *gorm.DB {
		dbEntity := db
		for keyName, condValue := range condMap {
			cond := StringJoin([]interface{}{" ", keyName, " = ?"})
			switch keyName {
			case SELECT_CONNECT_ACCESS_WAY, SELECT_CONNECT_MODULE_ID, SELECT_CONNECT_MODULE_SN:
				{
					dbEntity = dbEntity.Where(cond, condValue)
				}
			case SELECT_CONNECT_TIMETYPE:
				{
					dbEntity = addTimeCond(dbEntity, condValue.(string), request.StartTime, request.EndTime)
				}
			}
		}
		return dbEntity
	}); errMsg != nil {
		return nil, errMsg
	}

	return respond, nil
}
