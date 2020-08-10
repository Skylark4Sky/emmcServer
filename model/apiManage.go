package model

import (
	"GoServer/apiHandler"
	. "GoServer/middleWare"
	"github.com/gin-gonic/gin"
)

func ApiRegisterManage(router *gin.Engine, prometheusHttp func(context *gin.Context)) {
	api := router.Group("api").Use(prometheusHttp)
	{
		func(api gin.IRoutes) {
			api.GET("health-check", func(context *gin.Context) { context.AbortWithStatusJSON(200, gin.H{"ok": true}) })
			//登陆
			api.POST("login", func(context *gin.Context) { context.AbortWithStatusJSON(200, gin.H{"ok": true}) })
			//注册
			api.POST("register", func(context *gin.Context) { context.AbortWithStatusJSON(200, gin.H{"ok": true}) })
			//验证
			api.POST("checkCode", func(context *gin.Context) { context.AbortWithStatusJSON(200, gin.H{"ok": true}) })
			//设备
			api.POST("device", apiHandler.DeviceRegister)
			api.GET("device", apiHandler.DeviceRegister)
		}(api)

		auth := api.Use(JwtIntercept)

		func(auth gin.IRoutes) {
			auth.GET("ws/:id", UserRole(WsInterface), UserBehaviorIntercept(), func(context *gin.Context) { context.AbortWithStatusJSON(200, gin.H{"ok1": true}) })
			auth.GET("machine/:id", UserRole(MachineInterface), UserBehaviorIntercept(), func(context *gin.Context) {
				context.AbortWithStatusJSON(200, gin.H{"ok2": true})
			})
		}(auth)
	}

}
