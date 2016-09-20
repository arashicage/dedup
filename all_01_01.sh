#! /bin/bash
for i in $(cat url.01|grep -v px);
do
  host=$(echo $i |cut -f1 -d:);
  port=$(echo $i |cut -f2 -d:);
  echo "process" $host $port
  #echo "ok" $host $port > log/$host"_"$port".log"
  ./dedup_01.sh $host $port > log/"01_"$host"_"$port".log"
done;
