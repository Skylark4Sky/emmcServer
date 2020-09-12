package router

import (
	"GoServer/webApi/action"
	. "GoServer/webApi/middleWare"
	"github.com/gin-gonic/gin"
)

func registerUserRouter(userRouter *gin.RouterGroup) {
	user := userRouter.Group("user")
	{
		{
			user.POST("register", func(context *gin.Context) { context.AbortWithStatusJSON(200, gin.H{"ok": true}) })
			user.POST("login", action.Login)
			user.POST("weAppLogin", action.WeAppLogin)
			user.GET("weAppLogin", action.WeAppLogin)
			user.POST("findpassword")
		}

		authUser := user.Use(JwtIntercept)
		{
			authUser.POST("modify", func(context *gin.Context) { context.AbortWithStatusJSON(200, gin.H{"ok": true}) })
			authUser.POST("info", func(context *gin.Context) { context.AbortWithStatusJSON(200, gin.H{"ok": true}) })
		}
	}
}
