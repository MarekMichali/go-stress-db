#!/bin/bash

CONN_IT="--conn=25 --it=20"

if [ $# -eq 0 ]; then
    echo "Please provide the output file name as an argument."
    exit 1
fi

output_file=$1

echo "------------------------------------------------------------"
echo ""

for ((i=1; i<=10; i++))
do
    source ./scripts/configure-databases/mysql.sh
    sleep 20
    source ./scripts/configure-databases/mysql.sh
    cpus=$(bc <<< "scale=2; $i * 0.01")
    echo "The value of cpus is: $cpus"

    docker update --cpus $cpus mysql
    sleep 2
    cd app
    go run . --db=mysql $CONN_IT --op=insert
    echo ""
    cd ..
    sleep 1
    echo -e "\n insert $cpus" >> $output_file
    source ./mysql-verify-result.sh insert >> $output_file
    echo "------------------------------------------------------------"

    docker update --cpus $cpus mysql
    sleep 2
    cd app
    go run . --db=mysql $CONN_IT --op=select
    echo ""
    cd ..
    sleep 1
    echo -e "\n select $cpus" >> $output_file
    source ./mysql-verify-result.sh select >> $output_file
    echo "------------------------------------------------------------"

    docker update --cpus $cpus mysql
    sleep 2
    cd app
    go run . --db=mysql $CONN_IT --op=update
    echo ""
    cd ..
    sleep 1
    echo -e "\n update $cpus" >> $output_file
    source ./mysql-verify-result.sh update >> $output_file
    echo "------------------------------------------------------------"

    docker update --cpus $cpus mysql
    sleep 2
    cd app
    go run . --db=mysql $CONN_IT --op=delete
    echo ""
    cd ..
    sleep 1
    echo -e "\n delete $cpus" >> $output_file
    source ./mysql-verify-result.sh delete >> $output_file
    echo "------------------------------------------------------------"
    docker stop mysql
done
