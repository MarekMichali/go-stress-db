#!/bin/bash

docker start mariadbtest

docker exec -i mariadbtest bash << 'EOF'

mariadb -p"123456" << 'EOSQL'

USE videos;
DROP TABLE IF EXISTS videos;
SET GLOBAL tmp_table_size = 1024 * 1024 * 1024 * 4;
SET GLOBAL max_heap_table_size = 1024 * 1024 * 1024 * 4;
SET tmp_table_size = 1024 * 1024 * 1024 * 4;
SET max_heap_table_size = 1024 * 1024 * 1024 * 4;

CREATE TABLE videos (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(10) NOT NULL,
    data VARCHAR(16000)
) ENGINE=MEMORY;
EOSQL
EOF