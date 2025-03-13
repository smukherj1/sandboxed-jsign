#!/usr/bin/sh

set -eu

KEY="data/secrets/server.key"
CERT="data/secrets/server.crt"
TS_ADDR=":8080"
KMS_ADDR=":8081"

go run bin/signer-proxy/main.go \
  --key=${KEY} \
  --cert=${CERT} \
  --ts-addr=${TS_ADDR} \
  --kms-addr=${KMS_ADDR}
