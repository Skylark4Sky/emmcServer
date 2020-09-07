package utils

import (
	. "GoServer/webApi/middleWare"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RespondError(ctx *gin.Context, msg interface{}) {
	GaugeVecApiError.WithLabelValues("API").Inc()
	var ms string
	switch v := msg.(type) {
	case string:
		ms = v
	case error:
		ms = v.Error()
	default:
		ms = ""
	}
	ctx.AbortWithStatusJSON(200, gin.H{"error": false, "msg": ms})
}
func RetAuthError(ctx *gin.Context, msg interface{}) {
	ctx.AbortWithStatusJSON(http.StatusPreconditionFailed, gin.H{"ok": false, "msg": msg})
}

func RespondData(ctx *gin.Context, data interface{}) {
	ctx.AbortWithStatusJSON(200, data)
}

func RetSuccess(ctx *gin.Context) {
	ctx.AbortWithStatusJSON(200, gin.H{"ok": true, "msg": "success"})
}

func ChkError(ctx *gin.Context, err error) bool {
	if err != nil {
		RespondError(ctx, err.Error())
		return true
	}
	return false
}
