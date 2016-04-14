#!/bin/bash

set -x
set -e
rm -f pipeline
chmod +x ./test.sh
godep get -v -t
GOOS=linux godep go build -o pipeline
docker build -t pipeline-test .
