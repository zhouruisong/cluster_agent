#!/bin/sh
pid=`pidof fdfs-agent`
if [ ! -z $pid ]; then
		echo "kill pid $pid"
		`kill $pid`
fi
echo "fdfs-agent stop."

