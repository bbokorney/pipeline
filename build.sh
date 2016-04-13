#!/bin/bash

set -x
set -e
rm -f pipeline
godep get -v
GOOS=linux godep go build -o pipeline
docker build -t pipeline .
