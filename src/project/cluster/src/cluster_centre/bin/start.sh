#!/bin/sh
pid=`pidof cluster-centre`
if [ ! -z $pid ]; then
		`kill $pid`
fi

path="/usr/local/sandai/cluster/cluster-centre"
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

echo "cluster-centre start."
nohup ./cluster-centre --conf="../conf/cluster-centre-conf.json" >> $martini_log &
