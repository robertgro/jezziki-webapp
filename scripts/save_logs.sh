#!/bin/bash

# Define a timestamp function
timestamp() {
  date +"%Y-%m-%d_%H-%M-%S" # current time
  # format & credits https://stackoverflow.com/questions/17066250/create-timestamp-variable-in-bash-script
}

# echo $(timestamp)

journalctl -e -u jezziki-proxy.service &> /home/admin/logs/jezziki-proxy.service.$(timestamp).log

for i in {1..5}
do
  journalctl -e -u jezziki.node@$i.service &> /home/admin/logs/jezziki.node$i.$(timestamp).service.log
done

minimumKB=50000 # 50MB

fileArr=()

filePath=/home/admin/prod/jezziki-app/logs/

for logfile in $filePath*
do
    fileArr+=(${logfile##*/}) # get file name only, param expansion https://stackoverflow.com/questions/5920333/how-to-check-size-of-a-file-using-bash
done

for value in "${fileArr[@]}"
do
  path="$filePath$value" # param expansion possible https://www.cyberciti.biz/faq/bash-remove-last-character-from-string-line-word/
  # actualSize=$(wc -c <"$path") # file length in bytes
  kbSize=$(du -k "$path" | cut -f 1)
  if [ $kbSize -ge $minimumKB ]; then
      fileName=${value%.log} # without ext https://stackoverflow.com/questions/27658675/how-to-remove-last-n-characters-from-a-string-in-bash
      echo "$value file saved in /home/admin/logs" 
      cp $path /home/admin/logs/$fileName.$(timestamp).log
      echo "" > $path # do not remove the file or stdout is lost, just empty
  fi

done

filePath=/home/admin/prod/jezziki-proxy/logs/

for file in $filePath*
do
  fileName=${file#$filePath}
  kbSize=$(du -k "$path" | cut -f 1)
  if [ $kbSize -ge $minimumKB ]; then
    echo $fileName saved in /home/admin/logs
    cp $file /home/admin/logs/$fileName.$(timestamp).log
    echo "" > $file # do not remove the file or stdout is lost, just empty
  fi
done

rm -r /home/admin/prod/jezziki-proxy/whois/* >/dev/null 2>&1