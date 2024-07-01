#!/bin/bash

container_name="$1"
cpus="$2"

if [[ -z "$container_name" || -z "$cpus" ]]; then
    echo "Usage: $0 <container_name> <cpus> <flags to pass to the app>"
    echo "Example: $0 mysql 1.0 --db=mysql --op=update"
    exit 1
fi

docker update --cpus "$cpus" "$container_name"
echo "The value of cpus is: $cpus"
#docker inspect --format='{{.HostConfig.NanoCpus}}' "$container_name"
sleep 1
cd app

go run . "${@:3}"