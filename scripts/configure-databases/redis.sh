#!/bin/bash

docker start redis

docker exec -i redis redis-cli flushall
