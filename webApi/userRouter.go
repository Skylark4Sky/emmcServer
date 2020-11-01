package router

import (
	userApi "GoServer/controller/user"
	. "GoServer/middleWare/extension"
	"github.com/gin-gonic/gin"
)

func registerUserRouter(userRouter *gin.RouterGroup) {
	user := userRouter.Group("user")
	{
		{
			user.POST("login", userApi.Login)
			user.POST("logout", userApi.Logout)
			user.POST("weAppLogin", userApi.WechatLogin)
			user.GET("weAppLogin", userApi.WechatLogin)
		}

		authUser := user.Use(JwtIntercept)
		{
			authUser.POST("addUser", userApi.AddUser)
			authUser.POST("getUserInfo", userApi.GetUserInfo)
			// 微信小程序登录
			authUser.POST("updateWeAppUserInfo", userApi.WeChatUpdateUserInfo)
		}

	}
}
