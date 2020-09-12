package Service

import (
	. "GoServer/middleWare"
	. "GoServer/utils"
	. "GoServer/webApi/router"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
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
	setRunMode(GetWeb().Mode)

	if router != nil {
		router.NoMethod(ExceptionInterceptor)
		router.NoRoute(ExceptionInterceptor)
		ApiRegisterManage(router, PrometheusHttp)
	}

	server = &http.Server{
		Addr:           GetWeb().Port,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	Listen := make(chan error)
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			Listen <- err
		}
	}()

	var err error = nil
	select {
	case err = <-Listen:
	case <-time.After(1 * time.Second):
		break
	}
	close(Listen)
	return err
}

func StopWebService() {
	defer DBClose()
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
