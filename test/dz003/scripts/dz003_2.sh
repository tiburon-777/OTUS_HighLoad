#!/bin/bash

mysql -uroot -proot -h mysql_master -t app -e 'CREATE TABLE IF NOT EXISTS test (`id` INT(11), `name` VARCHAR(255));'


let i=1
let noerr=0

while [ $noerr = 0 ]
do
  mysql -uroot -proot -h mysql_master -t app -e 'INSERT INTO test (`id`,`name`) VALUES ('$i',"string_value_'$1'")' ||  ((noerr=1 ))
  echo $i
  ((i++))
done
echo "Err"