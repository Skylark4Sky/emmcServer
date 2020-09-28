package Service

import (
	. "GoServer/dataBases/mysql"
	"GoServer/middleWare"
	. "GoServer/utils"
	. "GoServer/webApi/router"
	"context"
	"fmt"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"time"
)

var server *http.Server = nil

const (
	DEBUG_MODE   uint32 = 0
	RELEASE_MODE uint32 = 1
	TEST_MODE    uint32 = 2
)

func setRunMode(mode uint32) {
	switch mode {
	case DEBUG_MODE:
		gin.SetMode(gin.DebugMode)
		fmt.Println("cur mode gin.DebugMode")
		break
	case RELEASE_MODE:
		gin.SetMode(gin.ReleaseMode)
		fmt.Println("cur mode gin.ReleaseMode")
		break
	case TEST_MODE:
		gin.SetMode(gin.TestMode)
		fmt.Println("cur mode gin.TestMode")
		break
	}
}

func StatrWebService() error {
	router := gin.New()

	webOption, _ := GetWeb()

	setRunMode(webOption.Mode)

	if router != nil {
		router.Use(middleWare.Cors())
		//router.Use(gin.Recovery())
		router.Use(middleWare.Recovery())
		router.Use(middleWare.RequestLogger())
		//router.Use(middleWare.ResponseHandler())
		router.NoMethod(middleWare.ExceptionInterceptor)
		router.NoRoute(middleWare.ExceptionInterceptor)
		ApiRegisterManage(router, middleWare.PrometheusHttp)
	}

	server = &http.Server{
		Addr:           webOption.Port,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// gracehttp可平滑重启
	if err := gracehttp.Serve(server); err != nil {
		WebLog("Start Server failed", zap.Error(err))
		return err
	}

	defer func(srv *http.Server) {
		err := recover()
		if err != nil {
			err := fmt.Errorf("panic %s", err)
			WebLog("Server Shutdown:", zap.Error(err))
			return
		}
		WebLog("Server Shutdown Success")
	}(server)

	WebLog("Start Server Success")
	//
	//Listen := make(chan error)
	//go func() {
	//	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
	//		Listen <- err
	//	}
	//}()
	//
	//var err error = nil
	//select {
	//case err = <-Listen:
	//case <-time.After(1 * time.Second):
	//	break
	//}
	//close(Listen)
	return nil
}

func StopWebService() {
	defer SQLClose()
	if server != nil {
		now := time.Now()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			fmt.Println("err", err)
		}
		select {
		case <-ctx.Done():
			fmt.Println("timeout of 5 seconds.")
		default:
			//fmt.Println("work")
		}
		fmt.Println("------exited--------", time.Since(now), ctx)
	}
	fmt.Println("StopWebService")
}
