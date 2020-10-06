package action

import (
	. "GoServer/handle/user"
	. "GoServer/middleWare/extension"
	. "GoServer/utils/config"
	. "GoServer/utils/respond"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

//普通用户登录
func Login(ctx *gin.Context) {
	var userLogin UserLogin

	if err := ctx.ShouldBind(&userLogin); err != nil {
		RespondMessage(ctx, CreateErrorMessage(PARAM_ERROR, err))
	}

	if userLogin.Account == "" {
		RespondMessage(ctx, CreateErrorMessage(PARAM_ERROR, "账号不能为空"))
		return
	}

	if userLogin.Pwsd == "" {
		RespondMessage(ctx, CreateErrorMessage(USER_PWSD_EMPTY, "密码不能为空"))
		return
	}

	data, err := userLogin.Login(ctx)

	if err != nil {
		RespondMessage(ctx, *err)
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
		RespondMessage(ctx, *loginErr)
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
	RespondMessage(ctx, nil)
}
