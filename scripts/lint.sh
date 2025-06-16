#!/bin/bash

echo "🔍 Running Go linters..."

# Check formatting first
echo "1. Checking code formatting..."
if ! ./scripts/check-format.sh; then
    echo "❌ Fix formatting first"
    exit 1
fi

# Run go vet
echo "2. Running go vet..."
if timeout 30 go vet ./...; then
    echo "✅ go vet passed"
else
    echo "❌ go vet failed or timed out"
    exit 1
fi

# Run go build to check compilation
echo "3. Testing compilation..."
if timeout 30 go build -o /tmp/test-build cmd/peerchat-cli/main.go; then
    echo "✅ Compilation successful"
    rm -f /tmp/test-build
else
    echo "❌ Compilation failed"
    exit 1
fi

echo "✅ All basic linting completed successfully!"
