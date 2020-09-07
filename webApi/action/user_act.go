package action

import (
	"GoServer/webApi/handle"
	. "GoServer/webApi/utils"
	"github.com/gin-gonic/gin"
)

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
