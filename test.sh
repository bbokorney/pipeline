#/bin/bash

set -x
set -e
echo "Creating dependencies"
docker-compose -f test.yml -p test_run up -d
echo "Running tests"
host=$(echo $DOCKER_HOST | awk -F/ '{print $3}' | awk -F: '{print $1}')
port=$(docker-compose port pipeline 4322 | awk -F: {'print $2'})
export PIPELINE_URL="http://$host:$port"
go test
docker-compose -f test.yml -p test_run kill
docker-compose -f test.yml -p test_run rm  -f
