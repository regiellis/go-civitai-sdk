#!/bin/bash

# Civitai API Tester - Cross-platform build script
set -e

APP_NAME="civitai-tester"
VERSION="1.0.0"
BUILD_DIR="builds"

echo "======================================================================"
echo "⚠️  SECURITY WARNING: DEVELOPMENT/TESTING TOOL ONLY"
echo "======================================================================"
echo "Building Civitai API Tester v${VERSION}..."
echo "⚠️  REMINDER: Do not use these binaries in production environments!"
echo "======================================================================"

# Clean previous builds
rm -rf ${BUILD_DIR}
mkdir -p ${BUILD_DIR}

# Build for different platforms
platforms=(
    "windows/amd64"
    "windows/arm64"
    "linux/amd64"
    "linux/arm64"
    "darwin/amd64"
    "darwin/arm64"
)

for platform in "${platforms[@]}"; do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    
    output_name=${APP_NAME}
    if [ $GOOS = "windows" ]; then
        output_name+='.exe'
    fi
    
    output_path=${BUILD_DIR}/${APP_NAME}-${GOOS}-${GOARCH}
    if [ $GOOS = "windows" ]; then
        output_path+='.exe'
    fi
    
    echo "Building for ${GOOS}/${GOARCH}..."
    env GOOS=$GOOS GOARCH=$GOARCH go build -ldflags="-s -w" -o $output_path .
    
    if [ $? -ne 0 ]; then
        echo "Error building for ${GOOS}/${GOARCH}"
        exit 1
    fi
done

echo ""
echo "Build completed successfully!"
echo "Binaries available in ./${BUILD_DIR}/"
ls -la ${BUILD_DIR}/