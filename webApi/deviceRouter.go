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
			authDevice.POST("startCharge", deviceApi.StartCharge)
			authDevice.POST("stopCharge", deviceApi.StopCharge)
			authDevice.POST("statusQuery", deviceApi.StatusQuery)
			authDevice.POST("noLoadSetting", deviceApi.NoLoadSetting)
			authDevice.POST("reStart", deviceApi.Restart)
			authDevice.POST("updateFirmware", deviceApi.UpdateFirmware)
			//常规控制器
			authDevice.POST("getDeviceList", deviceApi.GetDeviceList)
			authDevice.POST("getDeviceTransferLog", deviceApi.GetDeviceTransferLogList)
			authDevice.POST("getModuleList", deviceApi.GetModuleList)
			authDevice.POST("getModuleConnectLog", deviceApi.GetModuleConnectLogList)

			//authDevice.GET("getDeviceList", deviceApi.GetDeviceList)
			//authDevice.GET("getDeviceList", deviceApi.GetDeviceTransferLogList)
			//authDevice.GET("getDeviceList", deviceApi.GetModuleList)
			//authDevice.GET("getDeviceList", deviceApi.GetModuleConnectLogList)
		}
	}
}
