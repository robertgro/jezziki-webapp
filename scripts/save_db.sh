#!/bin/bash
# relies on .pgpass in ~

export PGPASSFILE=/home/admin/.pgpass

timestamp() {
  date +"%Y-%m-%d_%H-%M-%S" # current time
  # format & credits https://stackoverflow.com/questions/17066250/create-timestamp-variable-in-bash-script
}

filePath=/home/admin/dumps/jezziki.$(timestamp).sql

touch $filePath

pg_dump -U postgres jezziki > $filePath