package middleWare

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// ResponseHandler across domain
func ResponseHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.JSON(http.StatusOK, gin.H{
					"code":   http.StatusInternalServerError,
					"msg":    err,
					"result": "",
				})
				return
			}
		}()
		c.Next()
		err := c.Errors.ByType(gin.ErrorTypeAny).Last()
		if err != nil {
			if err.Meta != nil {
				c.JSON(http.StatusOK, err.Meta)
			} else {
				//if e, ok := err.Err.(errs.StandardError); ok {
				//	c.JSON(http.StatusOK, gin.H{
				//		"code": e.Code,
				//		"msg":  e.Msg,
				//	})
				//} else {
					c.JSON(http.StatusOK, gin.H{
						"code":   http.StatusInternalServerError,
						"msg":    err.Error(),
						"result": "",
					})
				//}
			}
		}
	}
}
