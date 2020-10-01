package router

import (
	. "GoServer/middleWare/extension"
	"github.com/gin-gonic/gin"
)

func registerPartnershipRouter(partnershipRouter *gin.RouterGroup) {
	partnership := partnershipRouter.Group("partnership").Use(JwtIntercept)
	{
		partnership.POST("add")
		partnership.POST("modify")
		partnership.POST("delete")
		partnership.POST("list")
	}
}
