CREATE DATABASE IF NOT EXISTS `binrpt`;
USE `binrpt`;
CREATE TABLE IF NOT EXISTS `replica_status` (
  `id` int(11) NOT NULL,
  `file` varchar(255) NOT NULL,
  `position` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB;
