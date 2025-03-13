#!/usr/bin/bash

go build -o out/signer-proxy bin/signer-proxy/main.go

docker build -f Dockerfile -t local-java:latest .