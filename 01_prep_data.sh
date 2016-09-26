#! /bin/bash

if [ $# -eq 0 ]; then
	echo "usage: "
	echo "	./01_prep_data port"
	echo "ex:	./01_prep_data 6379"
	exit
fi

for i in $(seq 1 100)
do
  i_with_prefix_zero="key"$(echo $i |awk '{printf("%09d\n",$0)}')
  redis-cli -p $1 -c hmset "04:"$i_with_prefix_zero h1 1 h2 2 h000001 1 h000002 2 h010001 1 h010002 2 h010003 3 > /dev/null
done
