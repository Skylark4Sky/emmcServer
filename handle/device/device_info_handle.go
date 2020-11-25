package device

import (
	. "GoServer/handle/user"
	. "GoServer/middleWare/dataBases/mysql"
	. "GoServer/middleWare/dataBases/redis"
	. "GoServer/model"
	. "GoServer/model/device"
	. "GoServer/utils/log"
	. "GoServer/utils/respond"
	. "GoServer/utils/string"
	"github.com/jinzhu/gorm"
	"math"
	"strconv"
	"strings"
)

const (
	SELECT_DEVICE_LIST              = 4
	SELECT_DEVICE_TRANSFER_LOG_LIST = 28
	SELECT_TMODULE_LIST             = 12
	SELECT_MODULE_CONNECT_LOG_LIST  = 20
	SYNC_DEVICE_STATUS              = 20
)

const (
	//批量修改设备状态分片
	SHARDING_SIZE = 20
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
	DEVICE_STATUS_KEY  = "status"
	SORT_FIELD_KEY     = "sortField"
	SORT_ORDER_KEY     = "sortOrder"
	STAR_TTIME_KEY     = "startTime"
	END_TIME_KEY       = "endTime"
	CREATE_TIME_KEY    = "create_time"
	UPDATE_TIME_KEY    = "update_time"
)

type RequestSyncData struct {
	UserID uint64 `fomr:"userID" json:"userID" binding:"required"`
}

type RequestListData struct {
	UserID      uint64      `fomr:"userID" json:"userID" binding:"required"`
	PageNum     int64       `form:"pageNum" json:"pageNum" binding:"required"`   //起始页
	PageSize    int64       `form:"pageSize" json:"pageSize" binding:"required"` //每页大小
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

func CheckUserRulesGroup(UserID uint64, roleValue int) (errMsg interface{}) {
	errMsg = nil
	userInfo := &UserInfo{}

	db := ExecSQL().Table("user_base")
	db = db.Select("user_role.id,user_role.rules")
	db = db.Joins("inner join user_role ON user_base.user_role = user_role.id")
	db = db.Where("uid = ?", UserID)

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

func addTimeCond(db *gorm.DB, timeField string, condMap map[string]interface{}) *gorm.DB {
	dbEntity := db

	var startTime int64 = 0
	var endTime int64 = 0

	if startTimeValue, ok := condMap[STAR_TTIME_KEY]; ok && startTimeValue != nil {
		startTime, _ = strconv.ParseInt(startTimeValue.(string), 10, 64)
	}
	if endTimeValue, ok := condMap[END_TIME_KEY]; ok && endTimeValue != nil {
		endTime, _ = strconv.ParseInt(endTimeValue.(string), 10, 64)
	}

	if startTime > 0 && endTime > 0 {
		cond := StringJoin([]interface{}{"(", timeField, " BETWEEN ? AND ?)"})
		dbEntity = dbEntity.Where(cond, startTime, endTime)
	} else {
		if startTime > 0 {
			cond := StringJoin([]interface{}{" ", timeField, " >= ?"})
			dbEntity = dbEntity.Where(cond, startTime)
		}
		if endTime > 0 {
			cond := StringJoin([]interface{}{" ", timeField, " <= ?"})
			dbEntity = dbEntity.Where(cond, endTime)
		}
	}
	return dbEntity
}

func generalSQLFormat(request *RequestListData, listSearch interface{}, condFilter func(db *gorm.DB, condMap map[string]interface{}) (*gorm.DB, string)) (errMsg interface{}, respond *RespondListData) {
	errMsg = nil
	respond = nil

	var total int64 = 0
	var orderCond string = ""

	db := ExecSQL()

	if request.RequestCond != nil {
		condMap := request.RequestCond.(map[string]interface{})
		if condFilter != nil {
			db, orderCond = condFilter(db, condMap)
		}
	}

	newScope := db.NewScope(listSearch)
	joinSQL := newScope.DB()
	tableName := newScope.GetModelStruct().TableName(db)

	totalRows := db.NewScope(listSearch).DB()

	if orderCond != "" {
		db = db.Order(orderCond).Limit(request.PageSize).Offset((request.PageNum - 1) * request.PageSize)
	} else {
		db = db.Order("id asc").Limit(request.PageSize).Offset((request.PageNum - 1) * request.PageSize)
	}

	db = db.Select("id")

	joinSQL = joinSQL.Joins(StringJoin([]interface{}{"inner join ? b ON ", tableName, ".id = b.id"}), db.Table(tableName).SubQuery())

	if err := joinSQL.Find(listSearch).Error; err != nil {
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

	if sortField, ok := condMap[SORT_FIELD_KEY]; ok {
		if sortOrder, ok := condMap[SORT_ORDER_KEY]; ok {
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

func addWhereCond(db *gorm.DB, condMap map[string]interface{}, key string) *gorm.DB {
	dbEntity := db
	switch key {
	case CREATE_TIME_KEY:
		{
			if keyValue, ok := condMap[TIMETYPE_KEY]; ok {
				timeType := keyValue.(string)
				if timeType == CREATE_TIME_KEY {
					dbEntity = addTimeCond(dbEntity, CREATE_TIME_KEY, condMap)
				}
			}
		}
	case UPDATE_TIME_KEY:
		{
			if keyValue, ok := condMap[TIMETYPE_KEY]; ok {
				timeType := keyValue.(string)
				if timeType == UPDATE_TIME_KEY {
					dbEntity = addTimeCond(dbEntity, UPDATE_TIME_KEY, condMap)
				}
			}
		}
	case DEVICE_STATUS_KEY:
		{
			if keyValue, ok := condMap[DEVICE_STATUS_KEY]; ok {
				if keyValue == "255" && keyValue != "" {
					break
				}
				if keyValue == "1" && keyValue != "" {
					cond := StringJoin([]interface{}{"(", key, " BETWEEN ? AND ?)"})
					dbEntity = dbEntity.Where(cond, DEVICE_ONLINE, DEVICE_WORKING)
				} else if keyValue != "" {
					cond := StringJoin([]interface{}{" ", key, " = ?"})
					dbEntity = dbEntity.Where(cond, keyValue)
				}
			}
		}
	default:
		if keyValue, ok := condMap[key]; ok {
			switch key {
			case ACCESS_WAY_KEY:
				if keyValue != "0" && keyValue != "" {
					cond := StringJoin([]interface{}{" ", key, " = ?"})
					dbEntity = dbEntity.Where(cond, keyValue)
				}
			case STAR_TTIME_KEY, END_TIME_KEY:
				break
			default:
				if keyValue != "" {
					cond := StringJoin([]interface{}{" ", key, " = ?"})
					dbEntity = dbEntity.Where(cond, keyValue)
				}
			}
		}
	}
	return dbEntity
}

func (request *RequestListData) GetDeviceList() (*RespondListData, interface{}) {

	var deviceList []DeviceInfo
	var respond *RespondListData = nil
	var errMsg interface{} = nil

	if errMsg, respond = generalSQLFormat(request, &deviceList, func(db *gorm.DB, condMap map[string]interface{}) (*gorm.DB, string) {
		db = addWhereCond(db, condMap, DEVICE_SN_KEY)
		db = addWhereCond(db, condMap, DEVICE_VERSION_KEY)
		db = addWhereCond(db, condMap, TYPE_KEY)
		db = addWhereCond(db, condMap, DEVICE_STATUS_KEY)
		db = addWhereCond(db, condMap, ACCESS_WAY_KEY)
		db = addWhereCond(db, condMap, CREATE_TIME_KEY)
		db = addWhereCond(db, condMap, UPDATE_TIME_KEY)
		return db, getOrderCond(condMap)
	}); errMsg != nil {
		return nil, errMsg
	}

	return respond, nil
}

func (request *RequestListData) GetDeviceTransferLogList() (interface{}, interface{}) {
	var transferList []DeviceTransferLog
	var respond *RespondListData = nil
	var errMsg interface{} = nil

	if errMsg, respond = generalSQLFormat(request, &transferList, func(db *gorm.DB, condMap map[string]interface{}) (*gorm.DB, string) {
		db = addWhereCond(db, condMap, TRANSFER_ID_KEY)
		db = addWhereCond(db, condMap, DEVICE_ID_KEY)
		db = addWhereCond(db, condMap, DEVICE_SN_KEY)
		db = addWhereCond(db, condMap, BEHAVIOR_KEY)
		db = addWhereCond(db, condMap, CREATE_TIME_KEY)

		defaultSortMap := map[string]interface{}{SORT_FIELD_KEY: "create_time", SORT_ORDER_KEY: DESCEND_ORDER}
		return db, getOrderCond(defaultSortMap)
	}); errMsg != nil {
		return nil, errMsg
	}

	return respond, nil
}

func (request *RequestListData) GetModuleList() (interface{}, interface{}) {
	var moduleList []ModuleInfo
	var respond *RespondListData = nil
	var errMsg interface{} = nil

	if errMsg, respond = generalSQLFormat(request, &moduleList, func(db *gorm.DB, condMap map[string]interface{}) (*gorm.DB, string) {
		db = addWhereCond(db, condMap, MODULE_SN_KEY)
		db = addWhereCond(db, condMap, DEVICE_ID_KEY)
		db = addWhereCond(db, condMap, ACCESS_WAY_KEY)
		db = addWhereCond(db, condMap, CREATE_TIME_KEY)
		db = addWhereCond(db, condMap, UPDATE_TIME_KEY)
		return db, getOrderCond(condMap)
	}); errMsg != nil {
		return nil, errMsg
	}

	return respond, nil
}

func (request *RequestListData) GetModuleConnectLogList() (interface{}, interface{}) {
	var connectList []ModuleConnectLog
	var respond *RespondListData = nil
	var errMsg interface{} = nil

	if errMsg, respond = generalSQLFormat(request, &connectList, func(db *gorm.DB, condMap map[string]interface{}) (*gorm.DB, string) {
		db = addWhereCond(db, condMap, MODULE_ID_KEY)
		db = addWhereCond(db, condMap, MODULE_SN_KEY)
		db = addWhereCond(db, condMap, ACCESS_WAY_KEY)
		db = addWhereCond(db, condMap, CREATE_TIME_KEY)
		return db, getOrderCond(condMap)
	}); errMsg != nil {
		return nil, errMsg
	}

	return respond, nil
}

func batchUpdatesDeviceStatus(entity interface{}, status int8) bool {
	deviceIDList := make([]uint64, 0)
	switch status {
	case DEVICE_OFFLINE:
		deviceMap := entity.(map[uint64]string)
		for deviceID, _ := range deviceMap {
			deviceIDList = append(deviceIDList, deviceID)
		}
		break
	case DEVICE_ONLINE, DEVICE_WORKING:
		deviceList := entity.([]interface{})
		for _, status := range deviceList {
			device := status.(DeviceStatus)
			deviceIDList = append(deviceIDList, device.ID)
		}
		break
	}

	var deviceIDListCount int = len(deviceIDList)

	if deviceIDListCount > 1 {
		//SystemLog("status ----> ", status, " deviceIDListCount ", deviceIDListCount, " deviceIDList ", deviceIDList)
		db := ExecSQL().Table("device_info")
		db = db.Where("id IN (?)", deviceIDList)
		db.Updates(map[string]interface{}{"status": status})
		return true
	}

	return false
}

func syncDeviceStatusTaskFunc(task *AsyncSQLTask) {
	if task.Entity == nil {
		task.Lock.Unlock()
		return
	}

	request := task.Entity.(*RequestSyncData)

	type ResultDeviceSNList struct {
		ID       uint64
		DeviceSN string
	}

	var resultList []ResultDeviceSNList
	db := ExecSQL().Table("device_info")
	db = db.Select("id,device_sn")
	db = db.Where("uid = ? AND status IN (?,?)", request.UserID, DEVICE_OFFLINE, DEVICE_ONLINE)
	if err := db.Scan(&resultList).Error; err != nil {
		SystemLog("syncDeviceStatusTaskFunc request: ", request, " Error: ", err)
		return
	}

	totalDevices := len(resultList)

	SystemLog("Beging update Device status total:", totalDevices)

	if totalDevices > 1 {
		var shardingNum int = 1
		if totalDevices > SHARDING_SIZE {
			shardingNum = int(math.Ceil(float64((totalDevices + (SHARDING_SIZE - 1)) / SHARDING_SIZE)))
		}

		var i int
		for i = 0; i < shardingNum; i++ {
			var start int = i * SHARDING_SIZE
			var offset int = SHARDING_SIZE
			if i == (shardingNum - 1) {
				offset = (totalDevices - (i * SHARDING_SIZE))
			}

			var k int
			deviceMap := make(map[uint64]string)
			for k = 0; k < offset; k++ {
				device := resultList[start+k]
				deviceMap[device.ID] = device.DeviceSN
			}

			onLine, workInLine := BatchReadDeviceTokenFromRedis(deviceMap)

			var needReletTime uint8 = 0

			if batchUpdatesDeviceStatus(deviceMap, DEVICE_OFFLINE) {
				needReletTime = (1 << 0)
			}

			if batchUpdatesDeviceStatus(onLine, DEVICE_ONLINE) {
				needReletTime = (1 << 1)
			}

			if batchUpdatesDeviceStatus(workInLine, DEVICE_WORKING) {
				needReletTime = (1 << 2)
			}

			if needReletTime > 1 {
				//再次续租锁时间
				err := task.Lock.Relet(int64(REDIS_LOCK_MINITIMEOUT))
				if err != nil {
					SystemLog("redis 锁续期 err: ", err)
				}
				SystemLog("redis 锁续期: ", REDIS_LOCK_MINITIMEOUT)
			}

		}
	}
}

func (request *RequestSyncData) SyncDeviceStatus() (interface{}, interface{}) {

	lock, ok, err := TryLock(StringJoin([]interface{}{"USERID_", request.UserID}), "syncDeviceStatus", int(REDIS_LOCK_DEFAULTIMEOUT))
	if err != nil {
		return nil, CreateErrorMessage(SYSTEM_ERROR, err)
	}
	if !ok {
		return nil, CreateErrorMessage(RESPOND_RESUBMIT, nil)
	}

	CreateAsyncSQLTaskWithCallback(ASYNC_UPDATE_DEVICE_STATUS, request, lock, syncDeviceStatusTaskFunc)
	return nil, CreateMessage(SUCCESS, "提交成功")
}
