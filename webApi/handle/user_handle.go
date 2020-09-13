package handle

import (
	. "GoServer/dataBases/mysql"
	. "GoServer/middleWare"
	. "GoServer/utils"
	. "GoServer/webApi/model"
	"fmt"
	"github.com/gin-gonic/gin"
)

// 登录成功返回
type UserLoginRespond struct {
	UserID    int64  `json:"userID"`
	UserName  string `json:"username"`
	NickName  string `json:"nickname"`
	Gender    int8   `json:"gender"`
	Birthday  int64  `json:"birthday"`
	Signature string `json:"signature"`
	Mobile    string `json:"mobile"`
	Email     string `json:"email"`
	Face      string `json:"face"`
	Face200   string `json:"face200"`
	Srcface   string `json:"srcface"`
}

// 小程序登录结构
type WeAppLogin struct {
	Auth       UserAuth
	Code       string `form:"code"`
	OpenID     string
	SessionKey string
}

// 登录绑定
type UserLogin struct {
	UserBase UserBase
	Account  string `json:"account"`
	Pwsd     string `json:"pwsd"`
}

type asynSQLTask struct {
	entity interface{}
}

func (task *asynSQLTask) ExecTask() error {
	switch entity := task.entity.(type) {
	case UserLoginLog:
		if err := ExecSQL().Create(&entity).Error; err != nil {
			fmt.Println("add login log Error:", err.Error())
		}
	case UserAuth:
		if err := ExecSQL().Model(&entity).Update("update_time", entity.UpdateTime).Error; err != nil {
			fmt.Println("update auth time Error:", err.Error())
		}
	}
	return nil
}

func createLoginRespond(entity *UserBase) *UserLoginRespond {
	return &UserLoginRespond{
		UserID:    entity.UId,
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

func getLoginType(account string, entity *UserBase) UserType {
	loginType := UNKNOWN

	switch account {
	case entity.Email:
		loginType = EMAIL
	case entity.UserName:
		loginType = USERNAME
	case entity.Mobile:
		loginType = MOBILE
	}

	return loginType
}

func createLoginLog(ctx *gin.Context, Command int8, loginType UserType, userID int64) {
	var task asynSQLTask
	var log UserLoginLog

	log.Create(ctx.ClientIP(), Command, loginType, userID)
	task.entity = log
	var work Job = &task
	InsertAsynTask(work)
}

func updateAuthTime(entity *UserAuth) {
	var task asynSQLTask
	task.entity = *entity
	var work Job = &task
	InsertAsynTask(work)
}

func createNewWechat(ip string, user *UserBase, M *WeAppLogin) error {
	tx := ExecSQL().Begin()

	//建立新用户
	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		return err
	}

	auth := &UserAuth{
		UId:          user.UId,
		IDentityType: int8(WECHAT),
		IDentifier:   M.OpenID,
		Certificate:  M.SessionKey,
		CreateTime:   GetTimestamp(),
	}

	//登记授权
	if err := tx.Create(&auth).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 登记日志
	log := &UserRegisterLog{
		UId:            user.UId,
		RegisterMethod: uint8(WECHAT),
		RegisterTime:   GetTimestamp(),
		RegisterIP:     ip,
		//	RegisterClient string
	}

	if err := tx.Create(&log).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

// 普通登录
func (M *UserLogin) Login(ctx *gin.Context) (*JwtObj, *RetMsg) {
	entity := &M.UserBase
	err := ExecSQL().Where("email = ? or user_name = ? or mobile = ?", M.Account, M.Account, M.Account).First(&entity).Error
	if err != nil {
		if IsRecordNotFound(err) {
			return nil, CreateRetStatus(USER_NO_EXIST, nil)
		}
		return nil, CreateRetStatus(SYSTEM_ERROR, err)
	}

	var loginType UserType = getLoginType(M.Account, entity)

	if chkOk := PasswordVerify(M.Pwsd, entity.UserPwsd); chkOk != true {
		createLoginLog(ctx, LOGIN_FAILURED, loginType, entity.UId)
		return nil, CreateRetStatus(USER_PWSD_ERROR, nil)
	}

	JwtData, err := JwtGenerateToken(createLoginRespond(entity), entity.UId)
	if err != nil {
		createLoginLog(ctx, LOGIN_FAILURED, loginType, entity.UId)
		return nil, CreateRetStatus(SYSTEM_ERROR, err)
	}

	createLoginLog(ctx, LOGIN_SUCCEED, loginType, entity.UId)
	return JwtData, nil
}

//小程序登录
func (M *WeAppLogin) Login(ctx *gin.Context) (*JwtObj, *RetMsg) {
	entity := &M.Auth
	err := ExecSQL().Select("uid").Where("identifier = ?", M.OpenID).First(&entity).Error

	var hasRecord = true

	if err != nil {
		if IsRecordNotFound(err) {
			hasRecord = false
		} else {
			return nil, CreateRetStatus(SYSTEM_ERROR, err)
		}
	}

	var user UserBase

	if !hasRecord {
		// 建立新用户
		user.CreateByDefaultInfo(WECHAT)
		if err := createNewWechat(ctx.ClientIP(), &user, M); err != nil {
			return nil, CreateRetStatus(SYSTEM_ERROR, err)
		}

	} else {
		entity.UpdateTime = GetTimestamp()
		updateAuthTime(entity)
		if err := ExecSQL().Where("uid = ?", entity.UId).First(&user).Error; err != nil {
			return nil, CreateRetStatus(SYSTEM_ERROR, err)
		}
	}

	JwtData, err := JwtGenerateToken(createLoginRespond(&user), entity.UId)
	if err != nil {
		createLoginLog(ctx, LOGIN_FAILURED, WECHAT, entity.UId)
		return nil, CreateRetStatus(SYSTEM_ERROR, err)
	}

	createLoginLog(ctx, LOGIN_SUCCEED, WECHAT, entity.UId)
	return JwtData, nil
}
