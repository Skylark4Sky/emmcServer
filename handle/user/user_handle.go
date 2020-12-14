package user

import (
	. "GoServer/middleWare/dataBases/mysql"
	. "GoServer/model/user"
	. "GoServer/model/asyncTask"
	. "GoServer/utils/string"
	"github.com/gin-gonic/gin"
)

type LoginRespond struct {
	UserInfo interface{} `json:"userInfo"`
	Token    interface{} `json:"tokenInfo"`
}

type UserData struct {
	UID       uint64 `json:"uid"`
	UserName  string `json:"username"`
	UserPwsd  string `json:"-"`
	NickName  string `json:"nickname"`
	Gender    uint8  `json:"gender"`
	Birthday  int64  `json:"birthday"`
	Signature string `json:"signature"`
	Face200   string `json:"face200"`
	Mobile    string `json:"mobile"`
	Email     string `json:"email"`
	RulesID	  int64  `json:"-"`
	Rules     string `json:"-"`
}

func createLoginRespond(entity *UserBase) *UserData {
	return &UserData{
		UID:       entity.UID,
		UserName:  entity.UserName,
		NickName:  entity.NickName,
		Gender:    entity.Gender,
		Birthday:  entity.Birthday,
		Signature: entity.Signature,
		Mobile:    entity.Mobile,
		Email:     entity.Email,
		Face200:   entity.Face200,
	}
}

func createLoginLog(ctx *gin.Context, Command uint8, loginType uint8, userID uint64) {
	log := &UserLoginLog{}
	log.Create(ctx.ClientIP(), Command, loginType, userID)
	NewAsyncTaskWithParam(ASYNC_USER_LOGIN_LOG,log)
}

func updateAuthTime(entity *UserAuth) {
	task := NewTask()
	task.Param = map[string]interface{}{"update_time": entity.UpdateTime}
	task.RunTaskWithTypeAndEntity(ASYNC_UP_USER_AUTH_TIME,entity)
}

func CheckUserIsExist(user *UserRegister) (bool, error) {
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
