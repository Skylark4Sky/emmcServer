#*
#* _COPYRIGHT_
#*
#* File Name: Makefile
#* System Environment: Linux VM_0_15_centos 3.10.0-514.26.2.el7.x86_64 #1 SMP Tue Jul 4 15:04:05 UTC 2017 x86_64 x86_64 x86_64 GNU/Linux
#* Created Time: 2019-12-10
#* Author: root
#* E-mail: root@qq.com
#* Description: 
#*
#*
#date -u '+%Y-%m-%dT%I:%M:%S%p'

export GO111MODULE=on
export GOPROXY=https://goproxy.cn/

BINNAME := "GoServer"
LDFLAGS := "-s -w -X 'main.BuildTime=$(shell date -u '+%F-%Z/%T')' -X 'main.GoVersion=$(shell go version)'"
GO ?= go
GOFMT ?= gofmt "-s"
PACKAGES ?= $(shell $(GO) list ./... | grep -v /vendor/)
VETPACKAGES ?= $(shell $(GO) list ./... | grep -v /vendor/ | grep -v /examples/)
GOFILES := $(shell find . -name "*.go" -type f -not -path "./vendor/*")
TIME := "log"$(shell date +"%Y%m%d")

exe:
	@sh ./exeRun.sh $(BINNAME)
bg: release
	@sh ./killExe.sh
	@sh ./exeRun.sh $(BINNAME) $(TIME)
run: debug
	@sh ./killExe.sh
	@sh ./exeRun.sh $(BINNAME)
release:
	@$(GO) build -ldflags $(LDFLAGS) -o ./build/$(BINNAME) main.go
#	upx -9 ./build/$(BINNAME)
debug:
	@$(GO) build -ldflags $(LDFLAGS) -o ./build/$(BINNAME) main.go
clean:
	@rm -fr ./build/$(BINNAME)
	@rm -fr ./crash
.PHONY: release

