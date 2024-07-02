#!/bin/bash

cd mariadb-container
docker build -t mariadb-custom -f Dockerfile.mariadb .

cd ../mongodb-container
docker build -t mongodb-custom -f Dockerfile.mongodb .

cd ../mysql-container
docker build -t mysql-custom -f Dockerfile.mysql .

