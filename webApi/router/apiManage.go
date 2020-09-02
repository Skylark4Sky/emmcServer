package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func ApiRegisterManage(router *gin.Engine, prometheusHttp func(context *gin.Context)) {
	//性能监控
	func(general *gin.Engine) {
		general.GET("/metrics", gin.WrapH(promhttp.Handler()))
		general.GET("/check", func(context *gin.Context) { context.AbortWithStatusJSON(200, gin.H{"ok": true}) })
		general.Use(cors.Default())
	}(router)

	api := router.Group("api")
	{
		api.Use(prometheusHttp)
		// 注册全局通用接口
		registerGeneralRouter(api)
		// 注册用户接口路由
		registerUserRouter(api)
		// 注册合作关系接口路由 客户&供应商
		registerPartnershipRouter(api)
		// 注册商品接口路由
		registerProductRouter(api)
		// 注册设置接口路由
		registerSettingRouter(api)
	}

}
