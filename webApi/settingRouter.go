package router

import (
	. "GoServer/middleWare/extension"
	"github.com/gin-gonic/gin"
)

func registerSettingRouter(settingRouter *gin.RouterGroup) {
	setting := settingRouter.Group("setting").Use(JwtIntercept)
	{
		setting.POST("add")
		setting.POST("modify")
		setting.POST("delete")
	}
}