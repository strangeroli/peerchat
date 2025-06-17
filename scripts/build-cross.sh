#!/bin/bash

# Xelvra P2P Messenger Cross-Platform Build Script
# Builds binaries for multiple platforms and creates distribution packages

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Project info
PROJECT_NAME="peerchat"
VERSION="0.4.0-alpha"
BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Build directories
DIST_DIR="dist"
BIN_DIR="bin"

# Target platforms
PLATFORMS=(
    "linux/amd64"
    "linux/arm64"
    "windows/amd64"
    "darwin/amd64"
    "darwin/arm64"
)

echo -e "${BLUE}🌍 Cross-platform build for Xelvra P2P Messenger v${VERSION}${NC}"
echo -e "${BLUE}Build time: ${BUILD_TIME}${NC}"
echo -e "${BLUE}Git commit: ${GIT_COMMIT}${NC}"
echo ""

# Create directories
mkdir -p ${DIST_DIR} ${BIN_DIR}

# Build flags
LDFLAGS="-X main.version=${VERSION} -X main.buildTime=${BUILD_TIME} -X main.gitCommit=${GIT_COMMIT} -s -w"

# Build for each platform
for platform in "${PLATFORMS[@]}"; do
    IFS='/' read -r -a platform_split <<< "$platform"
    GOOS="${platform_split[0]}"
    GOARCH="${platform_split[1]}"
    
    echo -e "${YELLOW}📦 Building for ${GOOS}/${GOARCH}...${NC}"
    
    # Set binary extension for Windows
    binary_ext=""
    if [ "$GOOS" = "windows" ]; then
        binary_ext=".exe"
    fi
    
    # Build CLI
    output_name="${DIST_DIR}/peerchat-cli-${VERSION}-${GOOS}-${GOARCH}${binary_ext}"
    env GOOS=$GOOS GOARCH=$GOARCH go build -ldflags "${LDFLAGS}" -o "$output_name" ./cmd/peerchat-cli
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ Built: $output_name${NC}"
        
        # Create archive
        archive_name="${DIST_DIR}/peerchat-${VERSION}-${GOOS}-${GOARCH}"
        if [ "$GOOS" = "windows" ]; then
            # Create ZIP for Windows
            zip -j "${archive_name}.zip" "$output_name" README.md LICENSE 2>/dev/null || true
            echo -e "${GREEN}✓ Archive: ${archive_name}.zip${NC}"
        else
            # Create tar.gz for Unix-like systems
            tar -czf "${archive_name}.tar.gz" -C "${DIST_DIR}" "$(basename "$output_name")" -C .. README.md LICENSE 2>/dev/null || true
            echo -e "${GREEN}✓ Archive: ${archive_name}.tar.gz${NC}"
        fi
    else
        echo -e "${RED}✗ Build failed for ${GOOS}/${GOARCH}${NC}"
    fi
done

# Create checksums
echo ""
echo -e "${YELLOW}🔐 Generating checksums...${NC}"
cd ${DIST_DIR}
sha256sum * > checksums.sha256 2>/dev/null || shasum -a 256 * > checksums.sha256
cd ..
echo -e "${GREEN}✓ Checksums generated: ${DIST_DIR}/checksums.sha256${NC}"

# Show build summary
echo ""
echo -e "${BLUE}📊 Build Summary:${NC}"
ls -lh ${DIST_DIR}/
echo ""

echo -e "${GREEN}🎉 Cross-platform build completed successfully!${NC}"
echo -e "${BLUE}Distribution packages available in: ${DIST_DIR}/${NC}"
