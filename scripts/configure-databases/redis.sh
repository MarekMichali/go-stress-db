#!/bin/bash

docker start redis

docker exec -i redis redis-cli config set slowlog-log-slower-than 0

docker exec -i redis redis-cli config set slowlog-max-len 200000

docker exec -i redis redis-cli slowlog reset

docker exec -i redis redis-cli slowlog get

docker exec -i redis redis-cli flushall
