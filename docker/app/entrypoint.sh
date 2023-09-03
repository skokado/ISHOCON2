#!/bin/bash
sleep 3
mysql -h db -u root -pishocon -e 'CREATE DATABASE IF NOT EXISTS ishocon2;' && \
mysql -h db -u root -pishocon -e "CREATE USER IF NOT EXISTS ishocon IDENTIFIED BY 'ishocon';" && \
mysql -h db -u root -pishocon -e 'GRANT ALL ON *.* TO ishocon;' && \
cd /data && tar -jxvf ishocon2.dump.tar.bz2 && mysql -h db -u root -pishocon ishocon2 < /data/ishocon2.dump

echo 'setup completed.'

exec "$@"
