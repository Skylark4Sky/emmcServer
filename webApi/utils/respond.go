package utils

import (
	. "GoServer/utils"
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

type RetMsg struct {
	Code              int64       `json:"code"`
	CurrentTimeMillis int64       `json:"currentTimeMillis"`
	Msg               string      `json:"msg"`
	Data              interface{} `json:"data"`
}

func CreateRetStatus(retCode int64, msg interface{}) *RetMsg {
	return CreateRetMsg(retCode, msg, nil)
}

func CreateRetMsg(retCode int64, msg interface{}, data interface{}) *RetMsg {
	retMsg, ok := retType[retCode]
	if ok {
		var msString string
		switch v := msg.(type) {
		case string:
			msString = v
		case error:
			msString = v.Error()
		default:
			msString = retMsg
		}
		return &RetMsg{
			Code:              retCode,
			Msg:               msString,
			CurrentTimeMillis: GetTimestampMs(),
			Data:              data,
		}
	}
	return nil
}

func RetError(ctx *gin.Context, msg interface{}) {
	data := msg.(RetMsg)
	ctx.AbortWithStatusJSON(200, gin.H{"code": data.Code, "currentTimeMillis": data.CurrentTimeMillis, "msg": data.Msg})
}

func RetData(ctx *gin.Context, data interface{}) {
	ctx.AbortWithStatusJSON(200, data)
}
