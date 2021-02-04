CREATE DATABASE IF NOT EXISTS `app` CHARACTER SET utf8 COLLATE utf8_general_ci;
GRANT ALL ON `app`.* TO 'app'@'%' identified by 'app';
GRANT REPLICATION SLAVE ON *.* TO 'app'@'%' identified by 'app';