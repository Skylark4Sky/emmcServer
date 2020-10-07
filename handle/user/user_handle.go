package user

import (
	. "GoServer/middleWare/dataBases/mysql"
	. "GoServer/middleWare/extension"
	. "GoServer/model"
	. "GoServer/model/user"
	. "GoServer/utils/respond"
	. "GoServer/utils/security"
	. "GoServer/utils/time"
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

//小程序更新用户信息
type WeAppUptdae struct {
	NickName string `json:"nickName"`
	Gender   int8   `json:"gender"`
	Language string `json:"language"`
	Face200  string `json:"avatarUrl"`
	City     string `json:"city"`
	Province string `json:"province"`
	Country  string `json:"country"`
	UserID   int64  `json:"userID"`
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
	var log UserLoginLog
	log.Create(ctx.ClientIP(), Command, loginType, userID)
	CreateAsyncSQLTask(ASYNC_USER_LOGIN_LOG, log)
}

func updateAuthTime(entity UserAuth) {
	CreateAsyncSQLTask(ASYNC_UP_USER_AUTH_TIME, entity)
}

func createNewWechatUser(ip string, user *UserBase, M *WeAppLogin) {
	// 用户授权
	auth := UserAuth{
		UID:          user.UID,
		IdentityType: int8(WECHAT),
		Identifier:   M.OpenID,
		Certificate:  M.SessionKey,
		CreateTime:   user.CreateTime,
	}

	CreateAsyncSQLTask(ASYNC_CREATE_USER_AUTH, auth)

	// 登记日志
	log := UserRegisterLog{
		UID:            user.UID,
		RegisterMethod: int8(WECHAT),
		RegisterTime:   user.CreateTime,
		RegisterIP:     ip,
	}

	CreateAsyncSQLTask(ASYNC_CREATE_USER_REGISTER_LOG, log)

	// 拓展字段
	extra := UserExtra{
		UID:        user.UID,
		CreateTime: user.CreateTime,
	}

	CreateAsyncSQLTask(ASYNC_CREATE_USER_EXTRA, extra)

	// 地址
	location := UserLocation{
		UID: user.UID,
	}

	CreateAsyncSQLTask(ASYNC_CREATE_USER_LOCATION, location)
}

// 普通登录
func (M *UserLogin) Login(ctx *gin.Context) (*JwtObj, *MessageEntity) {
	entity := &M.UserBase
	err := ExecSQL().Where("email = ? or user_name = ? or mobile = ?", M.Account, M.Account, M.Account).First(&entity).Error
	if err != nil {
		if IsRecordNotFound(err) {
			return nil, CreateErrorMessage(USER_NO_EXIST, nil)
		}
		return nil, CreateErrorMessage(SYSTEM_ERROR, err)
	}

	var loginType UserType = getLoginType(M.Account, entity)

	if chkOk := PasswordVerify(M.Pwsd, entity.UserPwsd); chkOk != true {
		createLoginLog(ctx, LOGIN_FAILURED, loginType, entity.UID)
		return nil, CreateErrorMessage(USER_PWSD_ERROR, nil)
	}

	JwtData, err := JwtGenerateToken(createLoginRespond(entity), entity.UID)
	if err != nil {
		createLoginLog(ctx, LOGIN_FAILURED, loginType, entity.UID)
		return nil, CreateErrorMessage(SYSTEM_ERROR, err)
	}

	createLoginLog(ctx, LOGIN_SUCCEED, loginType, entity.UID)
	return JwtData, nil
}

//小程序登录
func (M *WeAppLogin) Login(ctx *gin.Context) (*JwtObj, *MessageEntity) {
	entity := &M.Auth
	err := ExecSQL().Select("uid").Where("identifier = ?", M.OpenID).First(&entity).Error

	var hasRecord = true

	if err != nil {
		if IsRecordNotFound(err) {
			hasRecord = false
		} else {
			return nil, CreateErrorMessage(SYSTEM_ERROR, err)
		}
	}

	var user UserBase

	if !hasRecord {
		// 建立新用户
		user.CreateByDefaultInfo(WECHAT)

		lastID, err := CreateSQLAndRetLastID(user)

		if err != nil {
			return nil, CreateErrorMessage(SYSTEM_ERROR, err)
		}

		user.UID = lastID

		//建立其它关联表
		createNewWechatUser(ctx.ClientIP(), &user, M)
	} else {
		entity.UpdateTime = GetTimestamp()
		updateAuthTime(*entity)
		if err := ExecSQL().Where("uid = ?", entity.UID).First(&user).Error; err != nil {
			return nil, CreateErrorMessage(SYSTEM_ERROR, err)
		}
	}

	JwtData, err := JwtGenerateToken(createLoginRespond(&user), user.UID)

	if err != nil {
		createLoginLog(ctx, LOGIN_FAILURED, WECHAT, user.UID)
		return nil, CreateErrorMessage(SYSTEM_ERROR, err)
	}

	createLoginLog(ctx, LOGIN_SUCCEED, WECHAT, user.UID)
	return JwtData, nil
}

func (weApp *WeAppUptdae) Save() {

	var curTimestam int = GetTimestamp()

	userBase := UserBase{
		NickName:   weApp.NickName,
		Gender:     weApp.Gender,
		Face200:    weApp.Face200,
		UpdateTime: curTimestam,
	}

	CreateAsyncSQLTaskWithRecordID(ASYNC_UPDATA_WEUSER_INFO, weApp.UserID, userBase)

	userLocation := UserLocation{
		CurrNation:   weApp.Country,
		CurrProvince: weApp.Province,
		CurrCity:     weApp.City,
		UpdateTime:   curTimestam,
	}

	CreateAsyncSQLTaskWithRecordID(ASYNC_UPDATA_WEUSER_LOCAL, weApp.UserID, userLocation)

	userExtra := UserExtra{
		Language:   weApp.Language,
		UpdateTime: curTimestam,
	}

	CreateAsyncSQLTaskWithRecordID(ASYNC_UPDATA_USER_EXTRA, weApp.UserID, userExtra)

	//	Language string `json:"language"`
}