#!/bin/bash
while read line
do
	arr=($line)
	# TODO 
	# ssh user@example.com 'bash -s' < local_script.sh
	sshpass -p ${arr[2]} ssh user@example.com 'bash -s' < "mkdir -p dashboard-control"
	sshpass -p ${arr[2]} scp -r bin ${arr[1]}@${arr[0]}:~/dashboard-control/
	sshpass -p ${arr[2]} scp -r src ${arr[1]}@${arr[0]}:~/dashboard-control/
    # name=$line
    # echo "Text read from file - $name"
done < "raspberry_authdata.conf"