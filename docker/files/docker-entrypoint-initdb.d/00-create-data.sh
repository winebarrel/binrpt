#!/bin/bash
set -e
mysql -u root -e "CREATE DATABASE test DEFAULT CHARACTER SET utf8mb4"
mysql -u root test -e "CREATE TABLE test (id bigint primary key not null auto_increment, num int not null, str varchar(64) not null)"
mysql -u root test -e "INSERT INTO test (num, str) VALUES (100, 'foo'), (200, 'bar'), (300, 'zoo')"
mysql -u root test -e "CREATE TABLE secure_test (id bigint primary key not null auto_increment, num int not null, str varchar(64) not null)"
mysql -u root test -e "INSERT INTO secure_test (num, str) VALUES (100, 'foo'), (200, 'bar'), (300, 'zoo')"
