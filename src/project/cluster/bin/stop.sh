#!/bin/sh
if [ $# != 1 ];then
	echo "please input appname"
	exit 1
fi
pid=`pidof $1`
if [ ! -z $pid ]; then
		echo "kill pid $pid"
		`kill $pid`
fi
echo "$1 stop."

