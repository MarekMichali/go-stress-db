#!/bin/bash

# Start the MariaDB container
docker start mariadbtest

# Execute commands inside the MariaDB container without interactive mode
# by using a heredoc to pass commands to docker exec
docker exec -i mariadbtest bash << 'EOF'
# Run MariaDB with the password provided via an environment variable
mariadb -p"123456" << 'EOSQL'
# Run SQL commands
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