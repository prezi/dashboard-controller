#!/bin/bash
while read line
do
	arr=($line)
	# echo ${arr[0]}
	# echo ${arr[1]}
	# echo ${arr[2]}
	# echo '#######'
	# TODO 
	# ssh user@example.com 'bash -s' < local_script.sh
	# sshpass -p ${arr[2]} ssh user@example.com 'bash -s' < "mkdir -p dashboard-control"
	# sshpass -p ${arr[2]} scp -r bin ${arr[1]}@${arr[0]}:~/dashboard-control/
	# sshpass -p ${arr[2]} scp -r src ${arr[1]}@${arr[0]}:~/dashboard-control/
    # name=$line
    # echo "Text read from file - $name"
done < "device_authentication.conf"