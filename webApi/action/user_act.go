package action

import (
	. "GoServer/utils"
	"GoServer/webApi/handle"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

//普通用户登录
func Login(ctx *gin.Context) {
	var userLogin handle.UserLogin

	if err := ctx.ShouldBind(&userLogin); err != nil {
		RetError(ctx, CreateRetStatus(PARAM_ERROR, err))
	}

	if userLogin.Account == "" {
		RetError(ctx, CreateRetStatus(PARAM_ERROR, "账号不能为空"))
		return
	}

	if userLogin.Pwsd == "" {
		RetError(ctx, CreateRetStatus(USER_PWSD_EMPTY, "密码不能为空"))
		return
	}

	data, err := userLogin.Login(ctx)

	if err != nil {
		RetError(ctx, *err)
		return
	}

	RetData(ctx, CreateRetMsg(SUCCESS, nil, data))
}

// 微信小程序登录
func WeAppLogin(ctx *gin.Context) {
	var weApp handle.WeAppLogin
	if err := ctx.ShouldBind(&weApp); err != nil {
		RetError(ctx, CreateRetStatus(PARAM_ERROR, err))
		return
	}

	if weApp.Code == "" {
		RetError(ctx, CreateRetStatus(PARAM_ERROR, "code不能为空"))
		return
	}

	weAppConfig , _:= GetWeApp()
	appID := weAppConfig.AppID
	secret := weAppConfig.AppSecret
	CodeToSessURL := weAppConfig.CodeToSessURL
	CodeToSessURL = strings.Replace(CodeToSessURL, "{appid}", appID, -1)
	CodeToSessURL = strings.Replace(CodeToSessURL, "{secret}", secret, -1)
	CodeToSessURL = strings.Replace(CodeToSessURL, "{code}", weApp.Code, -1)

	resp, err := http.Get(CodeToSessURL)
	defer resp.Body.Close()

	if err != nil || resp.StatusCode != 200 {
		RetError(ctx, CreateRetStatus(SYSTEM_ERROR, "获取微信用户授权失败"))
		return
	}

	var respData map[string]interface{}

	err = json.NewDecoder(resp.Body).Decode(&respData)
	if err != nil {
		RetError(ctx, CreateRetStatus(SYSTEM_ERROR, err))
		return
	}

	weApp.OpenID = respData["openid"].(string)
	weApp.SessionKey = respData["session_key"].(string)

	if weApp.OpenID == "" || weApp.SessionKey == "" {
		RetError(ctx, CreateRetStatus(SYSTEM_ERROR, "微信认证失败"))
		return
	}

	data, loginErr := weApp.Login(ctx)

	if loginErr != nil {
		RetError(ctx, *loginErr)
		return
	}

	RetData(ctx, CreateRetMsg(SUCCESS, nil, data))
}
