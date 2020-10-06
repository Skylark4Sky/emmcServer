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
			user.POST("register", func(context *gin.Context) { context.AbortWithStatusJSON(200, gin.H{"ok": true}) })
			user.POST("login", action.Login)
			user.POST("weAppLogin", action.WechatLogin)
			user.GET("weAppLogin", action.WechatLogin)
			user.POST("findpassword")
		}

		authUser := user.Use(JwtIntercept)
		{
			authUser.POST("setWeAppUser", action.WeChatUpdateUserInfo)
			authUser.POST("modify", func(context *gin.Context) { context.AbortWithStatusJSON(200, gin.H{"ok": true}) })
			authUser.POST("info", func(context *gin.Context) { context.AbortWithStatusJSON(200, gin.H{"ok": true}) })
		}
	}
}
