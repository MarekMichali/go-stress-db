#!/bin/bash

container_name="$1"
cpus="$2"

if [[ -z "$container_name" || -z "$cpus" ]]; then
    echo "Usage: $0 <container_name> <cpus> <flags to pass to the app>"
    exit 1
fi

docker update --cpus "$cpus" "$container_name"
docker inspect --format='{{.HostConfig.NanoCpus}}' "$container_name"
sleep 1
cd app

go run . "${@:3}"