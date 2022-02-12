#! /usr/bin/bash
# credits https://stackoverflow.com/questions/39186854/bash-script-cant-execute-go-command
GO=/usr/local/go/bin/go
$GO get -u
$GO mod tidy