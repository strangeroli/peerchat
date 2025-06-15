#!/bin/bash

# Script to run specific linters that GitHub CI uses
# This avoids typecheck issues while still checking the important linters

echo "Running golangci-lint with specific linters..."

# Run errcheck
echo "Running errcheck..."
if ! /home/tux/go/bin/golangci-lint run --disable-all --enable=errcheck --no-config; then
    echo "❌ errcheck failed"
    exit 1
fi

# Run ineffassign  
echo "Running ineffassign..."
if ! /home/tux/go/bin/golangci-lint run --disable-all --enable=ineffassign --no-config; then
    echo "❌ ineffassign failed"
    exit 1
fi

# Run staticcheck
echo "Running staticcheck..."
if ! /home/tux/go/bin/golangci-lint run --disable-all --enable=staticcheck --no-config; then
    echo "❌ staticcheck failed"
    exit 1
fi

# Run unused
echo "Running unused..."
if ! /home/tux/go/bin/golangci-lint run --disable-all --enable=unused --no-config; then
    echo "❌ unused failed"
    exit 1
fi

echo "✅ All linters passed!"
