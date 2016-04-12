#/bin/bash

echo "Creating dependencies"
docker-compose up -d
echo "Running tests"
export PIPELINE_URL="http://dockerhost:4322"
go test
