package router

import (
	. "GoServer/webApi/middleWare"
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
