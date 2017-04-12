#!/bin/sh
pid=`pidof tair-agent`
if [ ! -z $pid ]; then
		echo "kill pid $pid"
		`kill $pid`
fi
echo "tair-agent stop."

