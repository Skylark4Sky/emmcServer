package action

import (
	. "GoServer/handle/user"
	. "GoServer/middleWare/extension"
	. "GoServer/utils/config"
	. "GoServer/utils/respond"
	. "GoServer/utils/security"
	. "GoServer/utils/string"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

//用户
func AddUser(ctx *gin.Context) {
	var register UserRegister
	if err := ctx.ShouldBind(&register); err != nil {
		RespondMessage(ctx, CreateErrorMessage(PARAM_ERROR, err))
		return
	}

	if register.Name == "" && register.Mobile == "" && register.Email == "" {
		RespondMessage(ctx, CreateErrorMessage(PARAM_ERROR, "参数错误"))
		return
	}

	userNameLen := len([]rune(register.Name))
	if userNameLen > 0 && userNameLen < 6 {
		RespondMessage(ctx, CreateErrorMessage(PARAM_ERROR, "用户名小于6位"))
		return
	}

	PwsdLen := len([]rune(register.Pwsd))

	if PwsdLen == 0 {
		RespondMessage(ctx, CreateErrorMessage(PARAM_ERROR, "密码不能空"))
		return
	} else if PwsdLen > 0 && PwsdLen < 6 {
		RespondMessage(ctx, CreateErrorMessage(PARAM_ERROR, "密码小于6位"))
		return
	}

	Pwsd, err := PasswordHash(register.Pwsd)
	if err != nil {
		RespondMessage(ctx, CreateErrorMessage(PARAM_ERROR, "参数错误"))
		return
	}

	register.Pwsd = Pwsd

	MobileLen := len([]rune(register.Mobile))
	if MobileLen > 0 && VerifyMobileFormat(register.Mobile) == false {
		RespondMessage(ctx, CreateErrorMessage(PARAM_ERROR, "手机格式错误"))
		return
	}

	EmailLen := len([]rune(register.Email))
	if EmailLen > 0 && VerifyEmailFormat(register.Email) == false {
		RespondMessage(ctx, CreateErrorMessage(PARAM_ERROR, "邮箱格式错误"))
		return
	}

	hasRecord, err := CheckUserIsExist(&register)

	if err != nil {
		RespondMessage(ctx, CreateErrorMessage(PARAM_ERROR, err))
		return
	}

	if hasRecord == true {
		RespondMessage(ctx, CreateErrorMessage(PARAM_ERROR, "用户已存在"))
		return
	}

	RespondMessage(ctx, register.Build(ctx))
}

//用户登录
func Login(ctx *gin.Context) {
	var login UserLogin

	if err := ctx.ShouldBind(&login); err != nil {
		RespondMessage(ctx, CreateErrorMessage(PARAM_ERROR, err))
		return
	}

	if login.Account == "" {
		RespondMessage(ctx, CreateErrorMessage(PARAM_ERROR, "账号不能为空"))
		return
	}

	if login.Pwsd == "" {
		RespondMessage(ctx, CreateErrorMessage(USER_PWSD_EMPTY, "密码不能为空"))
		return
	}

	data, err := login.Run(ctx)

	if err != nil {
		RespondMessage(ctx, err)
		return
	}

	RespondMessage(ctx, CreateMessage(SUCCESS, data))
}

func GetUserInfo(ctx *gin.Context) {
	userID := ctx.MustGet(JwtCtxUidKey)
	data, err := FetchUserInfo(userID.(uint64),ctx)

	if err != nil {
		RespondMessage(ctx, err)
		return
	}

	RespondMessage(ctx, CreateMessage(SUCCESS, data))
}

// 微信小程序登录
func WechatLogin(ctx *gin.Context) {
	var weApp WeAppLogin
	if err := ctx.ShouldBind(&weApp); err != nil {
		RespondMessage(ctx, CreateErrorMessage(PARAM_ERROR, err))
		return
	}

	if weApp.Code == "" {
		RespondMessage(ctx, CreateErrorMessage(PARAM_ERROR, "code不能为空"))
		return
	}

	weAppConfig, _ := GetWeApp()
	appID := weAppConfig.AppID
	secret := weAppConfig.AppSecret
	CodeToSessURL := weAppConfig.CodeToSessURL
	CodeToSessURL = strings.Replace(CodeToSessURL, "{appid}", appID, -1)
	CodeToSessURL = strings.Replace(CodeToSessURL, "{secret}", secret, -1)
	CodeToSessURL = strings.Replace(CodeToSessURL, "{code}", weApp.Code, -1)

	resp, err := http.Get(CodeToSessURL)
	defer resp.Body.Close()

	if err != nil || resp.StatusCode != 200 {
		RespondMessage(ctx, CreateErrorMessage(SYSTEM_ERROR, "获取微信用户授权失败"))
		return
	}

	var respData map[string]interface{}

	err = json.NewDecoder(resp.Body).Decode(&respData)
	if err != nil {
		RespondMessage(ctx, CreateErrorMessage(SYSTEM_ERROR, err))
		return
	}

	weApp.OpenID = respData["openid"].(string)
	weApp.SessionKey = respData["session_key"].(string)

	if weApp.OpenID == "" || weApp.SessionKey == "" {
		RespondMessage(ctx, CreateErrorMessage(SYSTEM_ERROR, "微信认证失败"))
		return
	}

	data, loginErr := weApp.Login(ctx)

	if loginErr != nil {
		RespondMessage(ctx, loginErr)
		return
	}

	RespondMessage(ctx, CreateMessage(SUCCESS, data))
}

func WeChatUpdateUserInfo(ctx *gin.Context) {
	userID := ctx.MustGet(JwtCtxUidKey)

	var weApp WeAppUptdae
	if err := ctx.ShouldBind(&weApp); err != nil {
		RespondMessage(ctx, CreateErrorMessage(PARAM_ERROR, err))
		return
	}

	if userID != weApp.UserID {
		RespondMessage(ctx, CreateErrorMessage(SYSTEM_ERROR, "参数错误"))
		return
	}

	//更新数据
	weApp.Save()

	//返回成功
	RespondMessage(ctx, CreateMessage(SUCCESS, nil))
}
