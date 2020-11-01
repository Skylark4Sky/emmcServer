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

const StartPage = 1

const (
	SELECT_DEVICE_LIST              = 4
	SELECT_DEVICE_TRANSFER_LOG_LIST = 28
	SELECT_TMODULE_LIST             = 12
	SELECT_MODULE_CONNECT_LOG_LIST  = 20
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

func checkUserRulesGroup(request *RequestListData, roleValue int) (isFind bool, errMsg interface{}) {

	isFind = false
	errMsg = nil

	userInfo := &UserInfo{}

	db := ExecSQL().Table("user_base")
	db = db.Select("user_role.id,user_role.rules")
	db = db.Joins("inner join user_role ON user_base.user_role = user_role.id")
	db = db.Where("uid = ?", request.UserID)

	//err := ExecSQL().Table("user_base").Select("user_role.id,user_role.rules").Joins("inner join user_role ON user_base.user_role = user_role.id").Where("uid = ?", request.UserID).Scan(&userInfo.User).Error
	if err := db.Scan(&userInfo.User).Error; err != nil {
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
	hasRole, errMsg := checkUserRulesGroup(request, SELECT_DEVICE_LIST)

	if errMsg != nil {
		return nil, errMsg
	}

	if !hasRole {
		return nil, CreateErrorMessage(SYSTEM_ERROR, "没有操作权限")
	}

	var deviceList []DeviceInfo
	var total int64 = 0

	db := ExecSQL().Debug()

	db = db.Limit(request.PageSize).Offset((request.PageNum - 1) * request.PageSize).Order("id desc")

	if request.RequestCond != nil {
		SystemLog("---request.RequestCond:", request.RequestCond)

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
	hasRole, errMsg := checkUserRulesGroup(request, SELECT_DEVICE_TRANSFER_LOG_LIST)

	if errMsg != nil {
		return nil, CreateErrorMessage(SYSTEM_ERROR, "没有操作权限")
	}

	if hasRole {

	}

	return nil, nil
}

func (request *RequestListData) GetModuleList() (interface{}, interface{}) {
	hasRole, errMsg := checkUserRulesGroup(request, SELECT_TMODULE_LIST)

	if errMsg != nil {
		return nil, CreateErrorMessage(SYSTEM_ERROR, "没有操作权限")
	}

	if hasRole {

	}

	return nil, nil
}

func (request *RequestListData) GetModuleConnectLogList() (interface{}, interface{}) {
	hasRole, errMsg := checkUserRulesGroup(request, SELECT_MODULE_CONNECT_LOG_LIST)

	if errMsg != nil {
		return nil, CreateErrorMessage(SYSTEM_ERROR, "没有操作权限")
	}

	if hasRole {

	}

	return nil, nil
}
