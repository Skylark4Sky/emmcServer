package action

import (
	"GoServer/webApi/handle"
	. "GoServer/webApi/utils"
	"github.com/gin-gonic/gin"
)

//普通用户登录
func Login(ctx *gin.Context) {
	var userLogin handle.UserLogin
	err := ctx.ShouldBind(&userLogin)

	if ChkError(ctx, err) {
		return
	}

	data, err := userLogin.Login(ctx.ClientIP())

	if ChkError(ctx, err) {
		return
	}

	if data == nil && err == nil {
		RespondError(ctx, "用户不存在")
		return
	}

	RespondData(ctx, data)
}

// 微信小程序登录
func WeAppLogin(ctx *gin.Context) {

}
