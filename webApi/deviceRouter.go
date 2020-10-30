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
		}

		authDevice := device.Use(JwtIntercept)
		{
			//操作控制器 权重最高,需慎重处理
			authDevice.POST("startCharge", action.DeviceStartCharge)
			authDevice.POST("stopCharge", action.DeviceStopCharge)
			authDevice.POST("statusQuery", action.DeviceStatusQuery)
			authDevice.POST("noLoadSetting", action.DeviceNoLoadSetting)
			authDevice.POST("reStart", action.DeviceRestart)
			authDevice.POST("updateFirmware", action.DeviceUpdateFirmware)
			//常规控制器
			authDevice.POST("getDeviceList", action.GetDeviceList)
			authDevice.POST("getDeviceTransferLog", action.GetDeviceTransferLogList)
			authDevice.POST("getModuleList", action.GetModuleList)
			authDevice.POST("getModuleConnectLog", action.GetModuleConnectLogList)
		}
	}
}
