#!/bin/bash

CONN_IT="--conn=25 --it=20"

if [ $# -eq 0 ]; then
    echo "Please provide the output file name as an argument."
    exit 1
fi

output_file=$1

echo "------------------------------------------------------------"
echo ""

echo -e "" >> $output_file
for ((i=1; i<=10; i++))
do
    source ./scripts/configure-databases/redis.sh
    sleep 20
    source ./scripts/configure-databases/redis.sh
    cpus=$(bc <<< "scale=2; $i * 0.01")
    echo "The value of cpus is: $cpus"

    docker update --cpus $cpus redis
    sleep 2
    cd app
    go run . --db=redis $CONN_IT --op=insert
    echo ""
    cd ..
    sleep 1
    echo -e "\n insert $cpus" >> $output_file
    source ./redis-verify-result.sh set >> $output_file
    echo "------------------------------------------------------------"

    docker update --cpus $cpus redis
    sleep 2
    cd app
    go run . --db=redis $CONN_IT --op=select
    echo ""
    cd ..
    sleep 1
    echo -e "\n select $cpus" >> $output_file
    source ./redis-verify-result.sh get >> $output_file
    echo "------------------------------------------------------------"

    docker update --cpus $cpus redis
    sleep 2
    cd app
    go run . --db=redis $CONN_IT --op=update
    echo ""
    cd ..
    sleep 1
    echo -e "\n update $cpus" >> $output_file
    source ./redis-verify-result.sh set >> $output_file
    echo "------------------------------------------------------------"

    docker update --cpus $cpus redis
    sleep 2
    cd app
    go run . --db=redis $CONN_IT --op=delete
    echo ""
    cd ..
    sleep 1
    echo -e "\n delete $cpus" >> $output_file
    source ./redis-verify-result.sh del >> $output_file
    echo "------------------------------------------------------------"
    docker stop redis
done
