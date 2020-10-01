package middleWare

import (
	. "GoServer/utils/time"
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"time"
)

var (
	GaugeVecApiDuration = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "apiDuration",
		Help: "api耗时单位ms",
	}, []string{"WSorAPI"})
	GaugeVecApiMethod = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "apiCount",
		Help: "各种网络请求次数",
	}, []string{"method", "path"})
	GaugeVecApiError = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "apiErrorCount",
		Help: "请求api错误的次数type: api/ws",
	}, []string{"type"})
)

func init() {
	prometheus.MustRegister(GaugeVecApiMethod, GaugeVecApiDuration, GaugeVecApiError)
}

func PrometheusHttp(ctx *gin.Context) {
	start := time.Now()
	method := ctx.Request.Method
	uri := ctx.Request.RequestURI
	GaugeVecApiMethod.WithLabelValues(method, uri).Inc()
	ctx.Next()
	end := time.Now()
	d := end.Sub(start) / time.Millisecond
	GaugeVecApiDuration.WithLabelValues(method).Set(float64(d))
}

func ExceptionInterceptor(ctx *gin.Context) {
	var buffer bytes.Buffer
	buffer.WriteString("IP:" + ctx.ClientIP())
	buffer.WriteString(" -t ")
	buffer.WriteString(TimeFormat(time.Now()))
	buffer.WriteString(" -a ")
	buffer.WriteString(ctx.Request.Method)
	buffer.WriteString(" -p ")
	buffer.WriteString(ctx.Request.URL.Path)
	ctx.AbortWithStatusJSON(http.StatusPreconditionFailed, gin.H{"error": buffer.String()})
}
