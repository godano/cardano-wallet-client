#!/usr/bin/env bash
# Requirements:
# Install oapi-codegen: https://github.com/deepmap/oapi-codegen
# Install goimports: go get golang.org/x/tools/cmd/goimports
# Install gofumpt: go get mvdan.cc/gofumpt

# Download latest Swagger definition of cardano-wallet
wget -O swagger.yaml "https://input-output-hk.github.io/cardano-wallet/api/edge/swagger.yaml"

function generate() {
    part="$1"
    oapi-codegen -generate "$part" \
        -package "wallet" \
        swagger.yaml > "wallet/generated-$part.go"
}

# Generated the different code parts
mkdir -p wallet
generate types
generate client
generate spec
generate server # Server not strictly necessary, but included for completeness

# Fix errors in the generated code - remove invalid type name prefixes
sed -i -e 's/200_//g' -e 's/202_//g' wallet/*.go

# Format code and fix imports
goimports -w wallet/*.go
gofumpt -w wallet/*.go
