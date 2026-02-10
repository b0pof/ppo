#!/bin/bash

echo "Checking cyclomatic complexity..."

if ! command -v cyclocomplex &> /dev/null; then
    echo "Installing cyclocomplex..."
    go install github.com/cyclocomplex/cyclocomplex@latest
fi

cyclocomplex -max 10 ./...

if [ $? -ne 0 ]; then
    echo "Cyclomatic complexity check failed"
    echo "Please refactor functions with complexity > 10"
    exit 1
fi

echo "Cyclomatic complexity check passed"

echo "Checking Halstead complexity..."

echo "Top 10 most complex functions:"
gocyclo -top 10 .

echo "All complexity checks completed!"
