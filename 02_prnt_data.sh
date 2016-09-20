#! /bin/bash

if [ $# -eq 0 ]; then
	echo "usage: "
	echo "	./02_prnt_data port"
	echo "ex:	./02_prnt_data 6379"
	exit
fi

IFS=$'\n'
arr=($(redis-cli -p $1 -c scan 0 count 1001))

x=1
while :
do
	# 处理key
	for i in $(seq 1 `expr ${#arr[@]} - 1` )
	do
	  fields=($(redis-cli -p $1 -c hkeys ${arr[$i]} |grep ^h|sort))
	  
	  echo $x ${arr[$i]} ${fields[@]}
		
	  x=`expr $x + 1`
	  
	done	

	# 遍历完了
	if [ ${arr[0]} -eq 0 ]; then
	  	break
	else
		arr=($(redis-cli -p $1 -c scan ${arr[0]} count 1001))
	fi

done