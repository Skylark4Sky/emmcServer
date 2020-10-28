package user

import (
	. "GoServer/middleWare/dataBases/mysql"
	. "GoServer/middleWare/extension"
	. "GoServer/model"
	. "GoServer/model/user"
	. "GoServer/utils/respond"
	. "GoServer/utils/time"
	"github.com/gin-gonic/gin"
)

// 小程序登录结构
type WeAppLogin struct {
	Auth       UserAuth
	Code       string `form:"code" json:"code" binding:"required"`
	OpenID     string
	SessionKey string
}

//小程序更新用户信息
type WeAppUptdae struct {
	NickName string `json:"nickName"`
	Gender   uint8  `json:"gender"`
	Language string `json:"language"`
	Face200  string `json:"avatarUrl"`
	City     string `json:"city"`
	Province string `json:"province"`
	Country  string `json:"country"`
	UserID   int64  `json:"userID"`
}

func createUserExtraInfo(ip string, user *UserBase) {
	// 登记日志
	log := &UserRegisterLog{
		UID:            user.UID,
		RegisterMethod: WECHAT,
		RegisterTime:   user.CreateTime,
		RegisterIP:     ip,
	}

	CreateAsyncSQLTask(ASYNC_CREATE_USER_REGISTER_LOG, log)

	// 拓展字段
	extra := &UserExtra{
		UID:        user.UID,
		CreateTime: user.CreateTime,
	}

	CreateAsyncSQLTask(ASYNC_CREATE_USER_EXTRA, extra)

	// 地址
	location := &UserLocation{
		UID: user.UID,
	}

	CreateAsyncSQLTask(ASYNC_CREATE_USER_LOCATION, location)
}

func createNewWechatUser(ip string, user *UserBase, M *WeAppLogin) {
	// 用户授权
	auth := &UserAuth{
		UID:          user.UID,
		IdentityType: WECHAT,
		Identifier:   M.OpenID,
		Certificate:  M.SessionKey,
		CreateTime:   user.CreateTime,
	}

	CreateAsyncSQLTask(ASYNC_CREATE_USER_AUTH, auth)
	createUserExtraInfo(ip, user)
}

//小程序登录
func (M *WeAppLogin) Login(ctx *gin.Context) (*JwtObj, interface{}) {
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

	var user *UserBase = &UserBase{}

	if !hasRecord {
		// 建立新用户
		user.CreateByDefaultInfo(WECHAT)

		lastID, err := CreateSQLAndRetLastID(user)

		if err != nil {
			return nil, CreateErrorMessage(SYSTEM_ERROR, err)
		}

		user.UID = lastID

		//建立其它关联表
		createNewWechatUser(ctx.ClientIP(), user, M)
	} else {
		entity.UpdateTime = GetTimestamp()
		updateAuthTime(entity)
		if err := ExecSQL().Where("uid = ?", entity.UID).First(&user).Error; err != nil {
			return nil, CreateErrorMessage(SYSTEM_ERROR, err)
		}
	}

	JwtData, err := JwtGenerateToken(createLoginRespond(user), user.UID)

	if err != nil {
		createLoginLog(ctx, LOGIN_FAILURED, WECHAT, user.UID)
		return nil, CreateErrorMessage(SYSTEM_ERROR, err)
	}

	createLoginLog(ctx, LOGIN_SUCCEED, WECHAT, user.UID)
	return JwtData, nil
}

//更新用户信息
func (weApp *WeAppUptdae) Save() {
	var curTimestam int64 = GetTimestamp()

	userBase := &UserBase{
		NickName:   weApp.NickName,
		Gender:     weApp.Gender,
		Face200:    weApp.Face200,
		UpdateTime: curTimestam,
	}
	CreateAsyncSQLTaskWithRecordID(ASYNC_UPDATA_WEUSER_INFO, weApp.UserID, userBase)

	userLocation := &UserLocation{
		CurrNation:   weApp.Country,
		CurrProvince: weApp.Province,
		CurrCity:     weApp.City,
		UpdateTime:   curTimestam,
	}
	CreateAsyncSQLTaskWithRecordID(ASYNC_UPDATA_WEUSER_LOCAL, weApp.UserID, userLocation)

	userExtra := &UserExtra{
		Language:   weApp.Language,
		UpdateTime: curTimestam,
	}
	CreateAsyncSQLTaskWithRecordID(ASYNC_UPDATA_USER_EXTRA, weApp.UserID, userExtra)
}
