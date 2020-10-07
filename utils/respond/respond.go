package respond

import (
	. "GoServer/utils/time"
	"github.com/gin-gonic/gin"
)

const (
	SUCCESS       = 0
	USER_NO_EXIST = 100001 + iota
	USER_PWSD_ERROR
	USER_ACCOUNT_EMPTY
	USER_PWSD_EMPTY
	PARAM_ERROR
	AUTH_ERROR
	SYSTEM_ERROR
)

var retType = map[int64]string{
	SUCCESS:            "opt Success",
	USER_NO_EXIST:      "User is no exist",
	USER_PWSD_ERROR:    "User passWord error",
	USER_ACCOUNT_EMPTY: "User account error",
	USER_PWSD_EMPTY:    "User passWord empty",
	PARAM_ERROR:        "param error",
	AUTH_ERROR:         "auth error",
	SYSTEM_ERROR:       "",
}

type MessageSucceedEntityWithData struct {
	Code              int64       `json:"code"`
	CurrentTimeMillis uint64      `json:"currentTimeMillis"`
	Data              interface{} `json:"data"`
}

type MessageSucceedEntity struct {
	Code              int64  `json:"code"`
	CurrentTimeMillis uint64 `json:"currentTimeMillis"`
}

type MessageFailedEntity struct {
	Code              int64  `json:"code"`
	CurrentTimeMillis uint64 `json:"currentTimeMillis"`
	Msg               string `json:"msg"`
}

func messageBuilder(retCode int64, msg interface{}, data interface{}) interface{} {
	retMsg, ok := retType[retCode]
	if ok {
		var retString string
		switch v := msg.(type) {
		case string:
			retString = v
		case error:
			retString = v.Error()
		default:
			retString = retMsg
		}

		//code != Success copy error msg
		if retCode != SUCCESS {
			entity := &MessageFailedEntity{}
			entity.CurrentTimeMillis = uint64(GetTimestampMs())
			entity.Code = retCode
			entity.Msg = retString
			return entity
		} else {
			if data != nil {
				entity := &MessageSucceedEntityWithData{}
				entity.CurrentTimeMillis = uint64(GetTimestampMs())
				entity.Code = retCode
				entity.Data = data
				return entity
			} else {
				entity := &MessageSucceedEntity{}
				entity.CurrentTimeMillis = uint64(GetTimestampMs())
				entity.Code = retCode
				return entity
			}
		}
	}
	return nil
}

func CreateErrorMessage(retCode int64, msg interface{}) interface{} {
	return messageBuilder(retCode, msg, nil)
}

func CreateMessage(retCode int64, data interface{}) interface{} {
	return messageBuilder(retCode, nil, data)
}

func RespondMessage(ctx *gin.Context, data interface{}) {
	ctx.AbortWithStatusJSON(200, data)
}
