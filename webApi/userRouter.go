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
			user.POST("adminLogin", action.AdminUserLogin)
			user.POST("weAppLogin", action.WechatLogin)
			user.GET("weAppLogin", action.WechatLogin)
		}

		authUser := user.Use(JwtIntercept)
		{
			authUser.POST("adminAdd", action.AdminUserAdd)
			authUser.POST("updateWeAppUserInfo", action.WeChatUpdateUserInfo)
			authUser.POST("modifyUserInfo", func(context *gin.Context) { context.AbortWithStatusJSON(200, gin.H{"ok": true}) })
			authUser.POST("modifyUserRole", func(context *gin.Context) { context.AbortWithStatusJSON(200, gin.H{"ok": true}) })
			authUser.POST("getUserInfo", func(context *gin.Context) { context.AbortWithStatusJSON(200, gin.H{"ok": true}) })
		}
	}
}
