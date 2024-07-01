#!/bin/bash

container_name="$1"
cpus="$2"

if [[ -z "$container_name" || -z "$cpus" ]]; then
    echo "------------------------------------------------------------"
    echo "Usage: ./run-test.sh <container_name> <cpus> <flags to pass to the app>"
    echo "Example: ./run-test.sh mysql 1.0 --db=mysql --op=update"
    echo ""
    echo "For more details about the flags, use the --help flag"
    echo "Example: ./run-test.sh mysql 1.0 --help"
    echo "------------------------------------------------------------"
    exit 1
fi

echo "------------------------------------------------------------"
echo ""
docker update --cpus "$cpus" "$container_name"
echo "The value of cpus is: $cpus"
#docker inspect --format='{{.HostConfig.NanoCpus}}' "$container_name"
sleep 1
cd app

go run . "${@:3}"
echo ""
echo "------------------------------------------------------------"