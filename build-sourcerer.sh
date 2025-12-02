#!/bin/bash
# build-sourcerer.sh - Build sourcerer from scratch

set -e  # Exit on any error

echo "Building Sourcerer from scratch..."

# Navigate to sourcerer directory
cd ~/github/n0shoes/sourcerer-mcp

echo "Cleaning previous builds..."
go clean -cache
#rm sourcerer

echo "Downloading dependencies..."
go mod download

echo "Tidying go.mod..."
go mod tidy

echo "Building sourcerer binary..."
go build -o sourcerer cmd/sourcerer/main.go

echo "Verifying binary..."
if [ -f sourcerer ]; then
    ls -lh sourcerer
    echo "✅ Build successful!"
    echo "Binary location: $(pwd)/sourcerer"
else
    echo "❌ Build failed - binary not found"
    exit 1
fi
