package Service

import (
	"GoServer/conf"
	. "GoServer/model"
	"bytes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

func prometheusHttp(context *gin.Context) {
	start := time.Now()
	method := context.Request.Method
	uri := context.Request.RequestURI
	GaugeVecApiMethod.WithLabelValues(method, uri).Inc()
	context.Next()
	end := time.Now()
	d := end.Sub(start) / time.Millisecond
	GaugeVecApiDuration.WithLabelValues(method).Set(float64(d))
}

func exceptionInterceptor(context *gin.Context) {
	var buffer bytes.Buffer
	buffer.WriteString("IP:" + context.ClientIP())
	buffer.WriteString(" -t ")
	buffer.WriteString(time.Now().Format(conf.GetConfig().GetSystem().Timeformat))
	buffer.WriteString(" -a ")
	buffer.WriteString(context.Request.Method)
	buffer.WriteString(" -p ")
	buffer.WriteString(context.Request.URL.Path)
	context.AbortWithStatusJSON(http.StatusPreconditionFailed, gin.H{"error": buffer.String()})
}

func setWebRouter(router *gin.Engine) {
	if router != nil {

		router.NoMethod(exceptionInterceptor)
		router.NoRoute(exceptionInterceptor)

		func(general *gin.Engine) {
			general.GET("/metrics", gin.WrapH(promhttp.Handler()))
			general.GET("/any-term", func(context *gin.Context) { context.AbortWithStatusJSON(200, gin.H{"ok": true}) })
			general.Use(cors.Default())
		}(router)
		ApiRegisterManage(router, prometheusHttp)
	}
}

func StatrWebService() error {
	gin.SetMode(gin.DebugMode)
	g := gin.New()

	setWebRouter(g)

	if err := g.Run(conf.GetConfig().GetWeb().Port); err != nil {
		return err
	}
	return nil

}
