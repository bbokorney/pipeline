#!/bin/bash

rm -f pipeline
godep get
GOOS=linux godep go build -o pipeline
docker build -t pipeline .
