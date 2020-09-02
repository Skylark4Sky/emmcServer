package router

import (
	"GoServer/webApi/action"
	"github.com/gin-gonic/gin"
)

func registerGeneralRouter(generalRouter *gin.RouterGroup) {
	//	func(api gin.IRoutes) {
	//手机验证码
	generalRouter.POST("sms", action.SMSLimiter, action.GetSMS)
	//设备
	generalRouter.POST("device", action.DeviceRegister)
	generalRouter.GET("device", action.DeviceRegister)
	//	}(api)

}
