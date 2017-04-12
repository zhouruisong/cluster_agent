#!/bin/sh
pid=`pidof mysql-agent`
if [ ! -z $pid ]; then
		echo "kill pid $pid"
		`kill $pid`
fi
echo "mysql-agent stop."

