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

// 后台用户登录
type AdminLogin struct {
	UserBase UserBase
	Account  string `form:"account" json:"account" binding:"required"`
	Pwsd     string `form:"pwsd" json:"pwsd" binding:"required"`
}

// 后台用户注册
type AdminRegister struct {
	Source    uint8  `form:"source" json:"source" binding:"required"`
	Name      string `form:"userName" json:"userName"`
	Pwsd      string `form:"userPwsd" json:"userPwsd" binding:"required"`
	NickName  string `form:"nickName" json:"nickName"`
	Gender    uint8  `form:"gender" json:"gender"`
	Birthday  int64  `form:"birthDay" json:"birthDay"`
	Signature string `form:"signature" json:"signature"`
	Mobile    string `form:"mobile" json:"mobile"`
	Email     string `form:"email" json:"email"`
}

// 普通登录
func (M *AdminLogin) Login(ctx *gin.Context) (*JwtObj, interface{}) {
	entity := &M.UserBase
	err := ExecSQL().Where("email = ? or user_name = ? or mobile = ?", M.Account, M.Account, M.Account).First(&entity).Error
	if err != nil {
		if IsRecordNotFound(err) {
			return nil, CreateErrorMessage(USER_NO_EXIST, nil)
		}
		return nil, CreateErrorMessage(SYSTEM_ERROR, err)
	}

	var loginType uint8 = getLoginType(M.Account, entity)

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

func (M *AdminRegister) Register(ctx *gin.Context) interface{} {

	var user CreateUserInfo
	user.Base = UserBase{
		RegisterSource: M.Source,
		UserRole:       2,
		UserName:       M.Name,
		UserPwsd:       M.Pwsd,
		NickName:       M.NickName,
		Gender:         M.Gender,
		Birthday:       M.Birthday,
		Signature:      M.Signature,
		Mobile:         M.Mobile,
		Email:          M.Email,
		CreateTime:     GetTimestamp(),
	}

	user.Log = UserRegisterLog{
		RegisterMethod: M.Source,
		RegisterTime:   user.Base.CreateTime,
		RegisterIP:     ctx.ClientIP(),
	}

	user.Extra = UserExtra{
		CreateTime: user.Base.CreateTime,
	}

	user.Location = UserLocation{}

	CreateAsyncSQLTask(ASYNC_CREATE_NORMAL_USER, user)

	return CreateMessage(SUCCESS, nil)
}
