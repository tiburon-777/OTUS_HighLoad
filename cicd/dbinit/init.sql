CREATE DATABASE IF NOT EXISTS `app` CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
GRANT ALL ON `app`.* TO 'app'@'%' identified by 'app';