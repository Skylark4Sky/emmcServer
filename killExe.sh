#!/bin/bash
#
# File Name: killExe.sh
# System Environment: Linux VM_0_13_centos 3.10.0-1127.13.1.el7.x86_64 #1 SMP Tue Jun 23 15:46:38 UTC 2020 x86_64 x86_64 x86_64 GNU/Linux
# Created Time: 2020-08-12
# Author: root
# E-mail: root@qq.com
# Description: 
#
#########################################################################

#!/bin/sh

exeName=$1

es_pid=`ps aux | grep .$exeName. | grep -v "grep" | tr -s ' '| cut -d ' ' -f 2`

if [ $es_pid ]; then
	kill -9 $es_pid
	echo -e "\033[31mkill PID "$es_pid"\033[0m"
fi

