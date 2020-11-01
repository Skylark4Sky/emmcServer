package router

import (
	deviceApi "GoServer/controller/device"
	. "GoServer/middleWare/extension"
	"github.com/gin-gonic/gin"
)

func registerDeviceRouter(deviceRouter *gin.RouterGroup) {
	device := deviceRouter.Group("device")
	{
		{
			device.POST("connect", deviceApi.DeviceConnect)
			device.GET("connect", deviceApi.DeviceConnect)
		}

		authDevice := device.Use(JwtIntercept)
		{
			//操作控制器 权重最高,需慎重处理
			authDevice.POST("startCharge", deviceApi.DeviceStartCharge)
			authDevice.POST("stopCharge", deviceApi.DeviceStopCharge)
			authDevice.POST("statusQuery", deviceApi.DeviceStatusQuery)
			authDevice.POST("noLoadSetting", deviceApi.DeviceNoLoadSetting)
			authDevice.POST("reStart", deviceApi.DeviceRestart)
			authDevice.POST("updateFirmware", deviceApi.DeviceUpdateFirmware)
			//常规控制器
			authDevice.POST("getDeviceList", deviceApi.GetDeviceList)
			authDevice.POST("getDeviceTransferLog", deviceApi.GetDeviceTransferLogList)
			authDevice.POST("getModuleList", deviceApi.GetModuleList)
			authDevice.POST("getModuleConnectLog", deviceApi.GetModuleConnectLogList)
		}
	}
}
