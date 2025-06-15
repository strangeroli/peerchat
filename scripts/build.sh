#!/bin/bash

# Xelvra P2P Messenger Build Script
# This script builds all components of the Xelvra messenger

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Project info
PROJECT_NAME="peerchat"
VERSION="0.1.0-alpha"
BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Build directories
BIN_DIR="bin"
DIST_DIR="dist"

echo -e "${BLUE}🔨 Building Xelvra P2P Messenger v${VERSION}${NC}"
echo -e "${BLUE}Build time: ${BUILD_TIME}${NC}"
echo -e "${BLUE}Git commit: ${GIT_COMMIT}${NC}"
echo ""

# Create directories
mkdir -p ${BIN_DIR} ${DIST_DIR}

# Build flags
LDFLAGS="-X main.version=${VERSION} -X main.buildTime=${BUILD_TIME} -X main.gitCommit=${GIT_COMMIT}"

echo -e "${YELLOW}📦 Building CLI application...${NC}"

# Build CLI for current platform
go build -ldflags "${LDFLAGS}" -o ${BIN_DIR}/peerchat-cli ./cmd/peerchat-cli
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ CLI built successfully: ${BIN_DIR}/peerchat-cli${NC}"
else
    echo -e "${RED}✗ CLI build failed${NC}"
    exit 1
fi

# Build API server
echo -e "${YELLOW}📦 Building API server...${NC}"
if [ -d "cmd/peerchat-api" ]; then
    go build -ldflags "${LDFLAGS}" -o ${BIN_DIR}/peerchat-api ./cmd/peerchat-api
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ API server built successfully: ${BIN_DIR}/peerchat-api${NC}"
    else
        echo -e "${RED}✗ API server build failed${NC}"
        exit 1
    fi
else
    echo -e "${YELLOW}⚠ API server not yet implemented${NC}"
fi

# Make binaries executable
chmod +x ${BIN_DIR}/*

# Show binary info
echo ""
echo -e "${BLUE}📊 Build Summary:${NC}"
ls -lh ${BIN_DIR}/
echo ""

# Test the CLI
echo -e "${YELLOW}🧪 Testing CLI...${NC}"
./${BIN_DIR}/peerchat-cli version
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ CLI test passed${NC}"
else
    echo -e "${RED}✗ CLI test failed${NC}"
    exit 1
fi

echo ""
echo -e "${GREEN}🎉 Build completed successfully!${NC}"
echo -e "${BLUE}Binaries available in: ${BIN_DIR}/${NC}"
echo ""
echo -e "${BLUE}Next steps:${NC}"
echo -e "  ./${BIN_DIR}/peerchat-cli init     # Initialize identity"
echo -e "  ./${BIN_DIR}/peerchat-cli start    # Start P2P node"
echo -e "  ./${BIN_DIR}/peerchat-cli --help   # Show help"
