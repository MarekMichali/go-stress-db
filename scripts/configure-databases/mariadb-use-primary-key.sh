#!/bin/bash

docker start mariadbtest

docker exec -i mariadbtest bash << 'EOF'

mariadb -p"123456" << 'EOSQL'

CREATE DATABASE IF NOT EXISTS videos;
USE videos;
DROP TABLE IF EXISTS videos;
SET max_heap_table_size = 1024 * 1024 * 1024 * 1;

CREATE TABLE videos (
    name VARCHAR(10) PRIMARY KEY NOT NULL,
    data VARCHAR(16000)
) ENGINE=MEMORY;
EOSQL
EOF