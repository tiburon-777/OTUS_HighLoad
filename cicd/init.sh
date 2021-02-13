#!/bin/bash

until docker exec mysql_master sh -c 'export MYSQL_PWD=root; mysql -u root -e ";"'
do
    echo "Waiting for mysql_master database connection..."
    sleep 4
done

priv_stmt='CREATE DATABASE IF NOT EXISTS app CHARACTER SET utf8 COLLATE utf8_general_ci; GRANT ALL ON app.* TO "app"@"%" IDENTIFIED BY "app"; GRANT REPLICATION SLAVE ON *.* TO "mydb_slave_user"@"%" IDENTIFIED BY "mydb_slave_pwd"; FLUSH PRIVILEGES;'
docker exec mysql_master sh -c "export MYSQL_PWD=root; mysql -u root -e '$priv_stmt'"

until docker exec mysql_slave1 sh -c 'export MYSQL_PWD=root; mysql -u root -e ";"'
do
    echo "Waiting for mysql_slave1 database connection..."
    sleep 4
done

until docker exec mysql_slave2 sh -c 'export MYSQL_PWD=root; mysql -u root -e ";"'
do
    echo "Waiting for mysql_slave2 database connection..."
    sleep 4
done

priv_stmt='CREATE DATABASE IF NOT EXISTS app CHARACTER SET utf8 COLLATE utf8_general_ci; GRANT ALL ON app.* TO "app"@"%" IDENTIFIED BY "app"; FLUSH PRIVILEGES;'
docker exec mysql_slave1 sh -c "export MYSQL_PWD=root; mysql -u root -e '$priv_stmt'"
docker exec mysql_slave2 sh -c "export MYSQL_PWD=root; mysql -u root -e '$priv_stmt'"

docker-ip() {
    docker inspect --format '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' "$@"
}

MS_STATUS=`docker exec mysql_master sh -c 'export MYSQL_PWD=root; mysql -u root -e "SHOW MASTER STATUS"' | grep mysq`
CURRENT_LOG=`echo $MS_STATUS | awk '{print $1}'`
CURRENT_POS=`echo $MS_STATUS | awk '{print $2}'`

start_slave_stmt="CHANGE MASTER TO MASTER_HOST='$(docker-ip mysql_master)',MASTER_USER='mydb_slave_user',MASTER_PASSWORD='mydb_slave_pwd',MASTER_LOG_FILE='$CURRENT_LOG',MASTER_LOG_POS=$CURRENT_POS; START SLAVE;"
start_slave_cmd='export MYSQL_PWD=root; mysql -u root -e "'
start_slave_cmd+="$start_slave_stmt"
start_slave_cmd+='"'

docker exec mysql_slave1 sh -c "$start_slave_cmd"
docker exec mysql_slave1 sh -c "export MYSQL_PWD=root; mysql -u root -e 'SHOW SLAVE STATUS \G' | grep Slave_"

docker exec mysql_slave2 sh -c "$start_slave_cmd"
docker exec mysql_slave2 sh -c "export MYSQL_PWD=root; mysql -u root -e 'SHOW SLAVE STATUS \G' | grep Slave_"