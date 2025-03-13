#!/usr/bin/bash

ARGS="-jar jsign.jar --storetype GOOGLECLOUD --keystore projects/bundle/locations/default/keyRings/key-2025 --storepass=default -a gulugulu -t localhost:8080 --proxyUrl localhost:8081 hello.ps1"

docker run \
  --rm \
  --runtime=runsc \
  --network=none \
  -v ./deps/jsign.jar:/wd/jsign.jar \
  -v ./data/hello.ps1:/wd/hello.ps1 \
  -w /wd \
  local-java:latest \
  ${ARGS}