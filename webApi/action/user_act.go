package action

import (
	"GoServer/webApi/handle"
	. "GoServer/webApi/utils"
	"github.com/gin-gonic/gin"
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

}
