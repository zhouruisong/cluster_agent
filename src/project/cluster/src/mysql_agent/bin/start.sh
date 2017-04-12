#!/bin/sh
pid=`pidof mysql-agent`
if [ ! -z $pid ]; then
		`kill $pid`
fi

path="/usr/local/sandai/cluster/mysql-agent"
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

echo "mysql-agent start."
nohup ./mysql-agent --conf="../conf/mysql-agent-conf.json" >> $martini_log &
