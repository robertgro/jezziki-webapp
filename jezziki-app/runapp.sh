#! /usr/bin/bash
GO=/usr/local/go/bin/go
Host='192.168.0.8'
Node=$1
PortList=('0' '8080' '8081' '8082' '8083' '8084')

if [ -z "$Node" ]; then
    echo "Please provide a valid node argument between 1-5"
    exit 1
fi

LFP=/home/admin/prod/jezziki-app/logs/jezziki.node0${Node}.log
touch $LFP

for i in "${!PortList[@]}"; do

    if [ $Node -eq $i ] && [ $Node -ne 0 ]; then
        $GO run . -port=${PortList[i]} -host=${Host} &>> $LFP
        #echo "${Node} is equal to ${i}, Port is ${PortList[i]} on Host ${Host}"
    fi

done