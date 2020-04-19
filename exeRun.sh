#!/bin/bash
#
# File Name: exeRun.sh
# System Environment: Darwin Lee-MacBook 19.3.0 Darwin Kernel Version 19.3.0: Thu Jan  9 20:58:23 PST 2020; root:xnu-6153.81.5~1/RELEASE_X86_64 x86_64
# Created Time: 2020-04-18
# Author: johan
# E-mail: johan@qq.com
# Description: 
#
#########################################################################

exeName=$1
logName=$2
ttyName=$(tty)
build="build"

if [ $exeName ]; then
	ulimit -c unlimited
	export GOTRACEBACK="system ./crash"
	env | grep -E "GO111MODULE|GOPROXY|PWD|GOTRACEBACK"
	runCmd="./$build/$exeName -v"

	if [ $logName ]; then
    #runCmd="${runCmd} > ./build/$logName &"
    runCmd="${runCmd} 1>/dev/null 2>&1 &"
	fi

	eval $runCmd
else	
	echo -e "\033[31m no exe run \033[0m"
fi