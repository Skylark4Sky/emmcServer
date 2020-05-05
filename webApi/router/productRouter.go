package router

import (
	. "GoServer/webApi/middleWare"
	"github.com/gin-gonic/gin"
)

func registerProductRouter(productRouter *gin.RouterGroup) {
	product := productRouter.Group("product").Use(JwtIntercept)
	{
		product.POST("add")
		product.POST("modify")
		product.POST("delete")
		product.POST("list")
	}
}