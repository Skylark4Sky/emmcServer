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
			//			user.POST("register", action.Register)      //admin register
			user.POST("login", action.Login)            //admin
			user.POST("weAppLogin", action.WechatLogin) //customer
			user.GET("weAppLogin", action.WechatLogin)  //customer
		}

		authUser := user.Use(JwtIntercept)
		{
			authUser.POST("setWeAppUser", action.WeChatUpdateUserInfo)
			authUser.POST("modify", func(context *gin.Context) { context.AbortWithStatusJSON(200, gin.H{"ok": true}) })
			authUser.POST("info", func(context *gin.Context) { context.AbortWithStatusJSON(200, gin.H{"ok": true}) })
		}
	}
}
