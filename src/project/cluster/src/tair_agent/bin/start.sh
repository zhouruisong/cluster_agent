#!/bin/sh
pid=`pidof tair-agent`
if [ ! -z $pid ]; then
		`kill $pid`
fi

path="/usr/local/sandai/cluster/tair-agent"
logs="$path/logs"
martini_log="$logs/access.log"
#echo "$logs"

if [ ! -d "$logs" ];then
		mkdir -p $logs
		if [ ! 0 -eq $? ];then
				echo "mkdir $logs failed."
				exit $?
		fi
fi

echo "tair-agent start."
nohup ./tair-agent --conf="../conf/tair-agent-conf.json" >> $martini_log &
