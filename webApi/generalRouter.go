package router

import (
	globalApi "GoServer/controller/global"
	deviceApi "GoServer/controller/device"
	"github.com/gin-gonic/gin"
)

func registerGeneralRouter(generalRouter *gin.RouterGroup) {
	//手机验证码
	generalRouter.POST("sms", globalApi.SMSLimiter, globalApi.GetSMS)
	//设备
	generalRouter.POST("device", deviceApi.DeviceConnect)
	generalRouter.GET("device", deviceApi.DeviceConnect)
	//generalRouter.GET("deviceList", action.DeviceList)
}
