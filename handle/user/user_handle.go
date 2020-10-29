package user

import (
	. "GoServer/middleWare/dataBases/mysql"
	. "GoServer/model"
	. "GoServer/model/user"
	. "GoServer/utils/string"
	"github.com/gin-gonic/gin"
)

// 登录成功返回
type UserLoginRespond struct {
	UserID    uint64 `json:"userID"`
	UserName  string `json:"username"`
	NickName  string `json:"nickname"`
	Gender    uint8  `json:"gender"`
	Birthday  int64  `json:"birthday"`
	Signature string `json:"signature"`
	Mobile    string `json:"mobile"`
	Email     string `json:"email"`
	Face      string `json:"face"`
	Face200   string `json:"face200"`
	Srcface   string `json:"srcface"`
}

func createLoginRespond(entity *UserBase) *UserLoginRespond {
	return &UserLoginRespond{
		UserID:    entity.UID,
		UserName:  entity.UserName,
		NickName:  entity.NickName,
		Gender:    entity.Gender,
		Birthday:  entity.Birthday,
		Signature: entity.Signature,
		Mobile:    entity.Mobile,
		Email:     entity.Email,
		Face:      entity.Face,
		Face200:   entity.Face200,
		Srcface:   entity.Srcface,
	}
}

func createLoginLog(ctx *gin.Context, Command uint8, loginType uint8, userID uint64) {
	log := &UserLoginLog{}
	log.Create(ctx.ClientIP(), Command, loginType, userID)
	CreateAsyncSQLTask(ASYNC_USER_LOGIN_LOG, log)
}

func updateAuthTime(entity *UserAuth) {
	CreateAsyncSQLTaskWithUpdateMap(ASYNC_UP_USER_AUTH_TIME, entity, map[string]interface{}{"update_time": entity.UpdateTime})
}

func CheckUserIsExist(user *AdminRegister) (bool, error) {
	entity := UserBase{}

	dict := make(map[string]string)

	if user.Name != "" {
		dict["user_name"] = user.Name
	}

	if user.Email != "" {
		dict["email"] = user.Email
	}

	if user.Mobile != "" {
		dict["mobile"] = user.Mobile
	}

	index := 0

	var itemValue []interface{}
	var condString string = ""

	for key, value := range dict {
		if index == 0 {
			condString = StringJoin([]interface{}{key, " = ?"})
		} else {
			condString = StringJoin([]interface{}{condString, " OR ", key, " = ?"})
		}
		itemValue = append(itemValue, value)
		index += 1
	}

	err := ExecSQL().Where(condString, itemValue...).First(&entity).Error
	var hasRecord = true
	if err != nil {
		if IsRecordNotFound(err) {
			hasRecord = false
		} else {
			return hasRecord, err
		}
	}

	return hasRecord, nil
}
