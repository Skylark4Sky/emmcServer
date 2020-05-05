package router

import (
	"GoServer/webApi/action"
	. "GoServer/webApi/middleWare"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func ApiRegisterManage(router *gin.Engine, prometheusHttp func(context *gin.Context)) {
	//通用接口
	func(general *gin.Engine) {
		general.GET("/metrics", gin.WrapH(promhttp.Handler()))
		general.GET("/any-term", func(context *gin.Context) { context.AbortWithStatusJSON(200, gin.H{"ok": true}) })
		general.Use(cors.Default())
	}(router)

	api := router.Group("api").Use(prometheusHttp)
	{
		//不需要验证的接口
		func(api gin.IRoutes) {
			api.GET("health-check", func(context *gin.Context) { context.AbortWithStatusJSON(200, gin.H{"ok": true}) })
			//登陆
			api.POST("login", action.Login)
			//注册
			api.POST("register", func(context *gin.Context) { context.AbortWithStatusJSON(200, gin.H{"ok": true}) })
			//验证
			api.POST("checkCode", func(context *gin.Context) { context.AbortWithStatusJSON(200, gin.H{"ok": true}) })

		}(api)

		auth := api.Use(JwtIntercept)

		//需要验证的接口
		func(auth gin.IRoutes) {
			auth.GET("ws/:id", UserRole(WsInterface), UserBehaviorIntercept(), func(context *gin.Context) {
				context.AbortWithStatusJSON(200, gin.H{"ok1": true})
			})
			auth.GET("machine/:id", UserRole(MachineInterface), UserBehaviorIntercept(), func(context *gin.Context) {
				context.AbortWithStatusJSON(200, gin.H{"ok2": true})
			})
		}(auth)
	}

}
