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

// 后台用户查询
type AdminUser struct {
	UID            uint64
	UserRole       uint8
	UserName       string
	UserPwsd	   string
	NickName       string
	Gender         uint8
	Birthday       int64
	Signature      string
	Face200        string
	Mobile         string
	Email          string
	Rules		   string
}

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
	//
	//{
	//roleId:
	//	'admin',
	//		permissionId: 'deviceManage',
	//	permissionName: '设备管理',
	//	actions: null,
	//	actionEntitySet: null,
	//	actionList: null,
	//	dataAccess: null
	//}
		//ExecSQL().Debug().Table("user_base").Select("user_base.uid,user_base.user_name,user_base.nick_name,user_base.gender,user_base.birthday,user_base.signature,user_base.face200,u.mobile,user_role.rules").Joins("inner join user_role ON user_base.user_role = user_role.id ").Where("email = ? or user_name = ? or mobile = ?", M.Account, M.Account, M.Account).Scan(&results)


	//uid	user_name	nick_name	gender	birthday	signature	face200	mobile	rules
//	ExecSQL().Debug().Joins("inner join user_role as r ON u.user_role = r.id").Where("email = ? or user_name = ? or mobile = ?", M.Account, M.Account, M.Account).Find(&user)

	//开始查用户权限
	//SELECT  FROM user_base as u  inner join user_role as r ON u.user_role = r.id WHERE (email = '13725467898' or user_name = '13725467898' or mobile = '13725467898');
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
		UserRole:       ADMIN_USER,
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
