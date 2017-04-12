#!/bin/sh
pid=`pidof fdfs-agent`
if [ ! -z $pid ]; then
		`kill $pid`
fi

path="/usr/local/sandai/cluster/fdfs-agent"
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

echo "fdfs-agent start."
nohup ./fdfs-agent --conf="../conf/fdfs-agent-conf.json" >> $martini_log &
