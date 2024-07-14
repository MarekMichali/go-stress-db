#!/bin/bash

echo "------------------------------------------------------------"
echo ""



for ((i=1; i<=10; i++))
do
    source ./scripts/configure-databases/mysql.sh
    sleep 5
    source ./scripts/configure-databases/mysql.sh
    cpus=$(bc <<< "scale=2; $i * 0.01")
    echo "The value of cpus is: $cpus"

    docker update --cpus $cpus mysql
    sleep 2
    cd app
    go run . --db=mysql --conn=25 --it=20 --op=insert
    echo ""
    cd ..
    echo -e "\n insert $cpus" >> output.txt
    source ./mysql-verify-result.sh insert >> output.txt
    echo "------------------------------------------------------------"

    docker update --cpus $cpus mysql
    sleep 2
    cd app
    go run . --db=mysql --conn=25 --it=20 --op=select
    echo ""
    cd ..
    echo -e "\n select $cpus" >> output.txt
    source ./mysql-verify-result.sh select >> output.txt
    echo "------------------------------------------------------------"

    docker update --cpus $cpus mysql
    sleep 2
    cd app
    go run . --db=mysql --conn=25 --it=20 --op=update
    echo ""
    cd ..
    echo -e "\n update $cpus" >> output.txt
    source ./mysql-verify-result.sh update >> output.txt
    echo "------------------------------------------------------------"

    docker update --cpus $cpus mysql
    sleep 2
    cd app
    go run . --db=mysql --conn=25 --it=20 --op=delete
    echo ""
    cd ..
    echo -e "\n delete $cpus" >> output.txt
    source ./mysql-verify-result.sh delete >> output.txt
    echo "------------------------------------------------------------"
    docker stop mysql
done


