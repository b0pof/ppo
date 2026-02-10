#!/bin/bash

echo "Fixing code style..."

goimports -w .
gofumpt -w .

golangci-lint run --fix ./...

echo "Code fixed"
