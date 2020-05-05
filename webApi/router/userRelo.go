package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

const (
	WsInterface      = 1
	MachineInterface = 2
)

//行为拦截
func UserBehaviorIntercept() gin.HandlerFunc {
	return func(context *gin.Context) {

	}
}

//权限处理
func UserRole(role uint) gin.HandlerFunc {
	return func(context *gin.Context) {
		fmt.Println("UserRole 2222")
		return
	}
}
