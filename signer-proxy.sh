#!/usr/bin/sh

set -eu

go run bin/signer-proxy/main.go --key=data/secrets/server.key --cert=data/secrets/server.crt
