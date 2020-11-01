package router

import (
	"GoServer/controller"
	. "GoServer/middleWare/extension"
	"github.com/gin-gonic/gin"
)

func registerUserRouter(userRouter *gin.RouterGroup) {
	user := userRouter.Group("user")
	{
		{
			user.POST("login", action.Login)
			user.POST("logout", action.Logout)
			user.POST("weAppLogin", action.WechatLogin)
			user.GET("weAppLogin", action.WechatLogin)
		}

		authUser := user.Use(JwtIntercept)
		{
			authUser.POST("addUser", action.AddUser)
			authUser.POST("getUserInfo", action.GetUserInfo)
			// 微信小程序登录
			authUser.POST("updateWeAppUserInfo", action.WeChatUpdateUserInfo)
		}
	}
}
