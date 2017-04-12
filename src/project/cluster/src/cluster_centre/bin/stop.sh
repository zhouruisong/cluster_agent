#!/bin/sh
pid=`pidof cluster-centre`
if [ ! -z $pid ]; then
		echo "kill pid $pid"
		`kill $pid`
fi
echo "cluster-centre stop."

