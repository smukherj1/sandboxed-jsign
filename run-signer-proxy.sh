#!/usr/bin/bash

set -eu

SIGN_KEY="data/secrets/sign.key"
KEY="data/secrets/server.key"
TS_CERT="data/secrets/ts.crt"
KMS_CERT="data/secrets/googlekms.crt"
TS_ADDR=":9001"
KMS_ADDR=":9000"

go run bin/signer-proxy/main.go \
  --key=${KEY} \
  --sign-key=${SIGN_KEY} \
  --kms-cert=${KMS_CERT} \
  --ts-cert=${TS_CERT} \
  --ts-addr=${TS_ADDR} \
  --kms-addr=${KMS_ADDR}
