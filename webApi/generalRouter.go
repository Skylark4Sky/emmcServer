package router

import (
	"GoServer/controller"
	"github.com/gin-gonic/gin"
)

func registerGeneralRouter(generalRouter *gin.RouterGroup) {
	//手机验证码
	generalRouter.POST("sms", action.SMSLimiter, action.GetSMS)
	//设备
	generalRouter.POST("device", action.DeviceRegister)
	generalRouter.GET("device", action.DeviceRegister)
}
