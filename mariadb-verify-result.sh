#!/bin/bash

docker update --cpus 2 mariadb
op_type=$1

if [[ "$op_type" == "insert" ]]; then
    result=$(docker exec -i mariadb bash <<'EOF'
    mariadb -p"123456" << 'EOSQL'
    SELECT SUM(count_star) as total_inserts, FORMAT_PICO_TIME(SUM(sum_timer_wait)) as total_time FROM performance_schema.events_statements_summary_by_digest WHERE digest_text LIKE 'insert%';
EOF
    )
elif [[ "$op_type" == "select" ]]; then
    result=$(docker exec -i mariadb bash <<'EOF'
    mariadb -p"123456" << 'EOSQL'
    SELECT SUM(count_star) as total_executions, FORMAT_PICO_TIME(SUM(sum_timer_wait)) as total_time FROM performance_schema.events_statements_summary_by_digest WHERE digest_text LIKE 'SELECT *%';
EOSQL
EOF
    )
elif [[ "$op_type" == "update" ]]; then
    result=$(docker exec -i mariadb bash <<'EOF'
    mariadb -p"123456" << 'EOSQL'
    SELECT SUM(count_star) as total_executions, FORMAT_PICO_TIME(SUM(sum_timer_wait)) as total_time FROM performance_schema.events_statements_summary_by_digest WHERE digest_text LIKE 'update%';
EOSQL
EOF
    )
elif [[ "$op_type" == "delete" ]]; then
    result=$(docker exec -i mariadb bash <<'EOF'
    mariadb -p"123456" << 'EOSQL'
    SELECT SUM(count_star) as total_executions, FORMAT_PICO_TIME(SUM(sum_timer_wait)) as total_time FROM performance_schema.events_statements_summary_by_digest WHERE digest_text LIKE 'delete%';EOSQL
EOF
    )
else
    echo "Invalid operation type"
fi

echo $result 
