#!/usr/bin/bash

args=(
  "-jar" "jsign.jar"
  "--storetype" "GOOGLECLOUD"
  "--keystore" "projects/default/locations/default/keyRings/default"
  "--storepass=default"
  "--alias" "default/cryptoKeyVersions/0"
  "--certfile=sign.crt"
  "-t" "10.0.0.3"
  "hello.ps1"
)

ARGS="${args[@]}"

docker run \
  --rm \
  --runtime=runsc \
  --network=docker-br0 \
  --add-host=cloudkms.googleapis.com:10.0.0.2 \
  -v ./deps/jsign.jar:/wd/jsign.jar \
  -v ./data/hello.ps1:/wd/hello.ps1 \
  -v ./data/secrets/sign.crt:/wd/sign.crt \
  -w /wd \
  local-java:latest \
  ${ARGS}