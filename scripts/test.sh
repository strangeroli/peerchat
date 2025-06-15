#!/bin/bash

# Xelvra P2P Messenger Test Script
# Runs comprehensive tests for all components

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🧪 Running Xelvra P2P Messenger Tests${NC}"
echo ""

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}✗ Go is not installed${NC}"
    exit 1
fi

# Run go mod tidy
echo -e "${YELLOW}📦 Tidying Go modules...${NC}"
go mod tidy
echo -e "${GREEN}✓ Go modules tidied${NC}"

# Run go vet
echo -e "${YELLOW}🔍 Running go vet...${NC}"
go vet ./...
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ go vet passed${NC}"
else
    echo -e "${RED}✗ go vet failed${NC}"
    exit 1
fi

# Run unit tests
echo -e "${YELLOW}🧪 Running unit tests...${NC}"
go test -v -race -coverprofile=coverage.out ./...
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Unit tests passed${NC}"
else
    echo -e "${RED}✗ Unit tests failed${NC}"
    exit 1
fi

# Generate coverage report
echo -e "${YELLOW}📊 Generating coverage report...${NC}"
go tool cover -html=coverage.out -o coverage.html
echo -e "${GREEN}✓ Coverage report generated: coverage.html${NC}"

# Build test
echo -e "${YELLOW}🔨 Testing build...${NC}"
go build -o /tmp/peerchat-cli-test ./cmd/peerchat-cli
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Build test passed${NC}"
    rm -f /tmp/peerchat-cli-test
else
    echo -e "${RED}✗ Build test failed${NC}"
    exit 1
fi

# Integration tests
echo -e "${YELLOW}🔗 Running integration tests...${NC}"
if [ -f "tests/integration_test.go" ]; then
    go test -v ./tests/...
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ Integration tests passed${NC}"
    else
        echo -e "${RED}✗ Integration tests failed${NC}"
        exit 1
    fi
else
    echo -e "${YELLOW}⚠ Integration tests not yet implemented${NC}"
fi

# CLI functional tests
echo -e "${YELLOW}⚙️ Running CLI functional tests...${NC}"

# Build CLI for testing
go build -o bin/peerchat-cli-test ./cmd/peerchat-cli

# Test help command
./bin/peerchat-cli-test --help > /dev/null
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ CLI help test passed${NC}"
else
    echo -e "${RED}✗ CLI help test failed${NC}"
    exit 1
fi

# Test version command
./bin/peerchat-cli-test version > /dev/null
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ CLI version test passed${NC}"
else
    echo -e "${RED}✗ CLI version test failed${NC}"
    exit 1
fi

# Test init command (in temporary directory)
TEMP_DIR=$(mktemp -d)
export HOME="$TEMP_DIR"
./bin/peerchat-cli-test init > /dev/null
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ CLI init test passed${NC}"
else
    echo -e "${RED}✗ CLI init test failed${NC}"
    exit 1
fi

# Cleanup
rm -rf "$TEMP_DIR"
rm -f bin/peerchat-cli-test

# Performance tests
echo -e "${YELLOW}⚡ Running performance tests...${NC}"
if [ -f "tests/performance_test.go" ]; then
    go test -bench=. -benchmem ./tests/...
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ Performance tests passed${NC}"
    else
        echo -e "${RED}✗ Performance tests failed${NC}"
        exit 1
    fi
else
    echo -e "${YELLOW}⚠ Performance tests not yet implemented${NC}"
fi

echo ""
echo -e "${GREEN}🎉 All tests passed successfully!${NC}"
echo ""
echo -e "${BLUE}📊 Test Summary:${NC}"
echo -e "  ✓ Go vet"
echo -e "  ✓ Unit tests"
echo -e "  ✓ Build test"
echo -e "  ✓ CLI functional tests"
if [ -f "coverage.out" ]; then
    COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}')
    echo -e "  📊 Test coverage: ${COVERAGE}"
fi
echo ""
echo -e "${BLUE}Reports generated:${NC}"
echo -e "  📄 coverage.html - Test coverage report"
