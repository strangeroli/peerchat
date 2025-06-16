#!/bin/bash

# Check Go code formatting
echo "🔍 Checking Go code formatting..."

# Check if any files need formatting
UNFORMATTED=$(gofmt -l .)

if [ -n "$UNFORMATTED" ]; then
    echo "❌ Code is not formatted properly:"
    echo "$UNFORMATTED"
    echo ""
    echo "🔧 To fix formatting, run:"
    echo "  gofmt -w ."
    exit 1
else
    echo "✅ All Go code is properly formatted"
    exit 0
fi
