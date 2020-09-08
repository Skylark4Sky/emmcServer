package action

import (
	. "GoServer/utils"
	"GoServer/webApi/handle"
	. "GoServer/webApi/utils"
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
	}

	if weApp.Code == "" {
		RetError(ctx, CreateRetStatus(PARAM_ERROR, "code不能为空"))
		return
	}

	weAppConfig := GetWeApp()
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

	var data map[string]interface{}

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		RetError(ctx, CreateRetStatus(SYSTEM_ERROR, err))
		return
	}

	weApp.OpenID = data["openid"].(string)
	weApp.SessionKey = data["session_key"].(string)

	data, err := weApp.Login(ctx)

	if err != nil {
		RetError(ctx, err)
		return
	}

}
