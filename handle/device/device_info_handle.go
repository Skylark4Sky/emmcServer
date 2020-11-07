package device

import (
	. "GoServer/handle/user"
	. "GoServer/middleWare/dataBases/mysql"
	. "GoServer/model/device"
	. "GoServer/utils/log"
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
	ASCEND_ORDER  = "ascend"
	DESCEND_ORDER = "descend"
)

const (
	DEVICE_ID_KEY      = "device_id"
	DEVICE_SN_KEY      = "device_sn"
	DEVICE_VERSION_KEY = "device_version"
	TYPE_KEY           = "type"
	ACCESS_WAY_KEY     = "access_way"
	TRANSFER_ID_KEY    = "transfer_id"
	BEHAVIOR_KEY       = "behavior"
	TIMETYPE_KEY       = "time"
	MODULE_ID_KEY      = "module_id"
	MODULE_SN_KEY      = "module_sn"
	MODULE_VERSION_KEY = "module_version"
	SORTFIELD          = "sortField"
	SORTORDER          = "sortOrder"
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
	Total      int64 `json:"total"`
	CurPageNum int64 `json:"pageNum"`
}

type RespondListData struct {
	List interface{} `json:"datalist"`
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
		cond := StringJoin([]interface{}{"(", timeField, " BETWEEN ? AND ?)"})
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

func generalSQLFormat(request *RequestListData, listSearch interface{}, condFilter func(db *gorm.DB, condMap map[string]interface{}) (*gorm.DB, string)) (errMsg interface{}, respond *RespondListData) {
	errMsg = nil
	respond = nil

	var total int64 = 0
	var orderCond string = ""

	db := ExecSQL().Debug()

	if request.RequestCond != nil {
		condMap := request.RequestCond.(map[string]interface{})
		if condFilter != nil {
			db, orderCond = condFilter(db, condMap)
		}
	} else {
		db = addTimeCond(db, "create_time", request.StartTime, request.EndTime)
	}

	totalRows := db.NewScope(listSearch).DB()

	if orderCond != "" {
		db = db.Order(orderCond).Limit(request.PageSize).Offset((request.PageNum - 1) * request.PageSize)
	} else {
		db = db.Limit(request.PageSize).Offset((request.PageNum - 1) * request.PageSize)
	}

	if err := db.Find(listSearch).Error; err != nil {
		errMsg = CreateErrorMessage(SYSTEM_ERROR, err)
		return
	}

	if err := totalRows.Count(&total).Error; err != nil {
		errMsg = CreateErrorMessage(SYSTEM_ERROR, err)
		return
	}

	respond = &RespondListData{
		List: listSearch,
		Page: PageInfo{
			Total:      total,
			CurPageNum: request.PageNum,
		},
	}

	return
}

func getOrderCond(condMap map[string]interface{}) (orderCond string) {
	orderCond = ""
	if sortField, ok := condMap[SORTFIELD]; ok {
		if sortOrder, ok := condMap[SORTORDER]; ok {
			var order string = "desc"
			if sortOrder == ASCEND_ORDER {
				order = "asc"
			} else {
				order = "desc"
			}
			orderCond = StringJoin([]interface{}{sortField, " ", order})
		}
	}
	return
}

func addWhereCond(db *gorm.DB, request *RequestListData, condMap map[string]interface{}, key string) *gorm.DB {
	dbEntity := db
	if keyValue, ok := condMap[key]; ok {
		switch key {
		case TIMETYPE_KEY:
			dbEntity = addTimeCond(dbEntity, keyValue.(string), request.StartTime, request.EndTime)
		case ACCESS_WAY_KEY:
			if keyValue != "0" {
				cond := StringJoin([]interface{}{" ", key, " = ?"})
				dbEntity = dbEntity.Where(cond, keyValue)
			}
		default:
			cond := StringJoin([]interface{}{" ", key, " = ?"})
			dbEntity = dbEntity.Where(cond, keyValue)
		}
	}
	return dbEntity
}

func (request *RequestListData) GetDeviceList() (*RespondListData, interface{}) {
	errMsg := checkUserRulesGroup(request, SELECT_DEVICE_LIST)

	if errMsg != nil {
		return nil, errMsg
	}

	var deviceList []DeviceInfo
	var respond *RespondListData = nil

	if errMsg, respond = generalSQLFormat(request, &deviceList, func(db *gorm.DB, condMap map[string]interface{}) (*gorm.DB, string) {
		db = addWhereCond(db, request, condMap, DEVICE_SN_KEY)
		db = addWhereCond(db, request, condMap, DEVICE_VERSION_KEY)
		db = addWhereCond(db, request, condMap, TYPE_KEY)
		db = addWhereCond(db, request, condMap, ACCESS_WAY_KEY)
		db = addWhereCond(db, request, condMap, TIMETYPE_KEY)
		return db, getOrderCond(condMap)
	}); errMsg != nil {
		return nil, errMsg
	}

	SystemLog("respond:-->", respond, "request:", request)
	return respond, nil
}

func (request *RequestListData) GetDeviceTransferLogList() (interface{}, interface{}) {
	errMsg := checkUserRulesGroup(request, SELECT_DEVICE_TRANSFER_LOG_LIST)

	if errMsg != nil {
		return nil, CreateErrorMessage(SYSTEM_ERROR, "没有操作权限")
	}

	var transferList []DeviceTransferLog
	var respond *RespondListData = nil

	if errMsg, respond = generalSQLFormat(request, &transferList, func(db *gorm.DB, condMap map[string]interface{}) (*gorm.DB, string) {
		db = addWhereCond(db, request, condMap, TRANSFER_ID_KEY)
		db = addWhereCond(db, request, condMap, DEVICE_ID_KEY)
		db = addWhereCond(db, request, condMap, DEVICE_SN_KEY)
		db = addWhereCond(db, request, condMap, BEHAVIOR_KEY)
		db = addWhereCond(db, request, condMap, TIMETYPE_KEY)
		return db, getOrderCond(condMap)
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

	if errMsg, respond = generalSQLFormat(request, &moduleList, func(db *gorm.DB, condMap map[string]interface{}) (*gorm.DB, string) {
		db = addWhereCond(db, request, condMap, MODULE_SN_KEY)
		db = addWhereCond(db, request, condMap, MODULE_VERSION_KEY)
		db = addWhereCond(db, request, condMap, ACCESS_WAY_KEY)
		db = addWhereCond(db, request, condMap, TIMETYPE_KEY)
		return db, getOrderCond(condMap)
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

	if errMsg, respond = generalSQLFormat(request, &connectList, func(db *gorm.DB, condMap map[string]interface{}) (*gorm.DB, string) {
		db = addWhereCond(db, request, condMap, MODULE_ID_KEY)
		db = addWhereCond(db, request, condMap, MODULE_SN_KEY)
		db = addWhereCond(db, request, condMap, ACCESS_WAY_KEY)
		db = addWhereCond(db, request, condMap, TIMETYPE_KEY)
		return db, getOrderCond(condMap)
	}); errMsg != nil {
		return nil, errMsg
	}

	return respond, nil
}
