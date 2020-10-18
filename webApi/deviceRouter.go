package router

import (
	action "GoServer/controller"
	. "GoServer/middleWare/extension"
	"github.com/gin-gonic/gin"
)

func registerDeviceRouter(deviceRouter *gin.RouterGroup) {
	device := deviceRouter.Group("device")
	{
		{
			device.POST("connect", action.DeviceConnect)
			device.GET("connect", action.DeviceConnect)
			device.GET("start_charge", action.DeviceStartCharge)

			device.GET("stop_charge", action.DeviceStopCharge)
		}

		authDevice := device.Use(JwtIntercept)
		{
			//authDevice.POST("start_charge", action.DeviceStartCharge)
			//	authDevice.POST("stop_charge", action.DeviceStopCharge)
			authDevice.POST("status_query", action.DeviceStatusQuery)
			authDevice.POST("no_load_setting", action.DeviceNoLoadSetting)
			authDevice.POST("restart", action.DeviceRestart)
		}
	}
}
