package router

import (
	. "GoServer/middleWare/extension"
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