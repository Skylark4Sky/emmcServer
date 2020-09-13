package main

import (
	. "GoServer/service"
	. "GoServer/utils"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

var (
	BuildTime = ""
	GoVersion = ""
)

func init() {
	fmt.Println("\nServer runing")
	if len(os.Args) == 2 && (os.Args[1] == "-v" || os.Args[1] == "-version") {
		fmt.Println("go version: \t" + GoVersion)
		fmt.Println("Build Time: \t" + BuildTime)
		fmt.Println("")
	}
}

func waitSignalExit() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		fmt.Println(sig)
		done <- true
	}()
	<-done
}

func main() {

	var err error
	var waitExit bool

	if _,err := GetConfig(); err != nil {
		panic(err)
	}



	if system, _ := GetSystem(); system != nil {
		if system.Service.Mqtt {
			if err = StartMqttService(); err != nil {
				fmt.Println("StartMqttService err:", err)
				panic(err)
			}
		}

		if system.Service.Web {
			if err = StatrWebService(); err != nil {
				fmt.Println("StatrWebService err:", err)
				panic(err)
			}
		}

		if system.Service.Mqtt || system.Service.Web {
			waitExit = true
		}
	}

	if waitExit {
		waitSignalExit()
	}

}
