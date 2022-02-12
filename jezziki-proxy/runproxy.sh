#! /usr/bin/bash
GO=/usr/local/go/bin/go
LFP=/home/admin/prod/jezziki-proxy/logs/jezziki-proxy.log
touch $LFP
$GO run . &>> $LFP