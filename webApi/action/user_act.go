package action

import (
	"GoServer/webApi/model"
	. "GoServer/webApi/utils"
	"github.com/gin-gonic/gin"
)

func Login(ctx *gin.Context) {
	var user model.UserLogin
	err := ctx.ShouldBind(&user)

	if ChkError(ctx, err) {
		return
	}

	data,err := user.Login(ctx.ClientIP())

	if ChkError(ctx, err) {
		return
	}

	if data == nil && err == nil {
		RetError(ctx,"用户不存在")
		return
	}

	RetData(ctx,data)
}