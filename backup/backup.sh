#!/bin/bash

#Backup
BACKUP_DIR="/backup"
DATE=`date +%d-%m-%Y-%H-%M`
OUTPUT_FILE="$BACKUP_DIR/$DATE.sql"
OUTPUT_FILE_TGZ="$OUTPUT_FILE.tar.gz"
USER=root
PASSWORD=root

if [ ! -z $PASSWORD ]; then
    mysqldump -u$USER -p$PASSWORD --all-databases > "$OUTPUT_FILE"
else
    mysqldump -u$USER --all-databases > "$OUTPUT_FILE"
fi

tar -czf $OUTPUT_FILE.tar.gz $OUTPUT_FILE
rm $OUTPUT_FILE
echo "Backup : $OUTPUT_FILE.tar.gz"

#Send to servers

SCP_USER="hub-backup"
SERVERS=( "127.0.0.1" )
REMOTE_BACKUP_FOLDER="/hub-backup"

for server in "${SERVERS[@]}"
do
    :
    scp $OUTPUT_FILE_TGZ $SCP_USER@$server:$REMOTE_BACKUP_FOLDER
    echo $server
done
