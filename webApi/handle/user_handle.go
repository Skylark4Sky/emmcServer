package handle

import (
	. "GoServer/utils"
	. "GoServer/webApi/middleWare"
	. "GoServer/webApi/model"
	. "GoServer/webApi/utils"
	"fmt"
	"github.com/gin-gonic/gin"
)

// 登录成功返回
type UserLoginRespond struct {
	UserID    int64
	UserName  string
	NickName  string
	Gender    int8
	Birthday  int64
	Signature string
	Mobile    string
	Email     string
	Face      string
	Face200   string
	Srcface   string
}

// 登录绑定
type UserLogin struct {
	UserBase UserBase
	Account  string `json:"account"`
	Pwsd     string `json:"pwsd"`
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

func addLoginLog(ctx *gin.Context, Command int8, loginType UserType, userID int64) {
	//userAgent := ctx.Request.Header.Get("User-Agent")
	//fmt.Println(userAgent)

	var entity = &UserLoginLog{
		UId:        userID,
		Type:       int8(loginType),
		CreateTime: GetTimestamp(),
		Command:    Command,
		Lastip:     ctx.ClientIP(),
	}
	err := DBInstance.Debug().Create(entity).Error
	if err != nil {
		fmt.Println(err)
	}
	return
}

// 普通登录
func (M *UserLogin) Login(ctx *gin.Context) (*JwtObj, *RetMsg) {
	if M.Pwsd == "" {
		return nil, CreateRetStatus(USER_PWSD_EMPTY, nil)
	}

	entity := &M.UserBase
	cond := fmt.Sprintf("email = '%s' or user_name = '%s' or mobile = '%s'", M.Account, M.Account, M.Account)
	err := DBInstance.Debug().Where(cond).First(&entity).Error
	if err != nil {
		if IsRecordNotFound(err) {
			return nil, CreateRetStatus(USER_NO_EXIST, nil)
		}
		return nil, CreateRetStatus(SYSTEM_ERROR, err)
	}

	var loginType UserType = getLoginType(M.Account, entity)

	if chkOk := PasswordVerify(M.Pwsd, entity.UserPwsd); chkOk != true {
		addLoginLog(ctx, LOGIN_FAILURED, loginType, entity.UId)
		return nil, CreateRetStatus(USER_PWSD_ERROR, nil)
	}

	JwtData, err := JwtGenerateToken(createLoginRespond(entity), entity.UId)
	if err != nil {
		addLoginLog(ctx, LOGIN_FAILURED, loginType, entity.UId)
		return nil, CreateRetStatus(SYSTEM_ERROR, err)
	}

	addLoginLog(ctx, LOGIN_SUCCEED, loginType, entity.UId)
	return JwtData, nil
}

//小程序登录
func WeAppLogin(ctx *gin.Context) (*JwtObj, *RetMsg) {
	return nil, nil
}
