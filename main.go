package main

import (
	. "GoServer/mqtt"
	. "GoServer/utils"
	. "GoServer/webApi"
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

func registerExitSignal() {
	signals := make(chan os.Signal, 1)
	exitRun := make(chan bool, 1)
	signal.Notify(signals, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)
	go func() {
		for exit := false; exit != true; {
			signal := <-signals
			fmt.Println("signal handle:->",signal)
			switch signal {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				StopMqttService()
				StopWebService()
				exitRun <- true
				exit = true;
				break
			case syscall.SIGUSR1:
				fmt.Println("usr1", signal)
			case syscall.SIGUSR2:
				fmt.Println("usr2", signal)
			default:
				fmt.Println("other", signal)
			}
		}
	}()
	<-exitRun
	fmt.Println("Server exit")
}

func main() {
	var err error

	if SystemConf() == nil {
		panic(ErrConfString)
	}

	if SystemConf().Service.Mqtt {
		if err = StartMqttService(); err != nil {
			fmt.Println("StartMqttService err:", err)
			panic(err)
		}
	}

	if SystemConf().Service.Web {
		if err = StatrWebService(); err != nil {
			fmt.Println("StatrWebService err:", err)
			panic(err)
		}
		string := fmt.Sprintf("curl -H \"Content-Type: application/json\" -X POST -d '{\"msg\":\"value\"}' \"http://localhost:%s/api/*\"",WebConf().Port)
		fmt.Println(string)
	}
	registerExitSignal()
}
