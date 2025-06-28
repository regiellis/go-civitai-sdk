#!/bin/bash

# Go CivitAI SDK - Library Structure Verification
#
# Copyright (c) 2025 Regi Ellis
# Licensed under Restricted Use License - Non-Commercial Only

echo "🔍 Verifying Go CivitAI SDK Library Structure"
echo "============================================="

# Check if we're in the right directory
if [ ! -f "go.mod" ]; then
    echo "❌ Error: go.mod not found. Please run this script from the SDK root directory."
    exit 1
fi

echo "✅ Directory structure check:"
echo "   📂 Root directory: $(pwd)"

# Check required files
required_files=(
    "go.mod"
    "README.md"
    ".gitignore"
    "client.go"
    "types.go"
    "exceptions.go"
    "client_test.go"
    "types_test.go"
    "integration_test.go"
)

for file in "${required_files[@]}"; do
    if [ -f "$file" ]; then
        echo "   ✅ $file"
    else
        echo "   ❌ Missing: $file"
    fi
done

# Check required directories
required_dirs=(
    "examples"
    "cmd"
)

for dir in "${required_dirs[@]}"; do
    if [ -d "$dir" ]; then
        echo "   ✅ $dir/"
    else
        echo "   ❌ Missing directory: $dir/"
    fi
done

echo ""
echo "📦 Go module verification:"
go mod verify
if [ $? -eq 0 ]; then
    echo "✅ go.mod is valid"
else
    echo "❌ go.mod verification failed"
fi

echo ""
echo "🧪 Running tests:"
go test -v -short
test_result=$?

echo ""
echo "🔨 Building test program:"
cd cmd/test
go build -o test-sdk
build_result=$?

if [ $build_result -eq 0 ]; then
    echo "✅ Test program built successfully"
    echo "📋 Test program info:"
    ls -la test-sdk
    echo ""
    echo "🚀 Running quick validation (with timeout):"
    timeout 10s ./test-sdk || echo "⚠️  Test timed out or failed (this is expected with API issues)"
else
    echo "❌ Test program build failed"
fi

cd ../..

echo ""
echo "📚 Examples verification:"
cd examples
for example in *.go; do
    if [ -f "$example" ]; then
        echo "   📝 Checking syntax: $example"
        go build -o /dev/null "$example"
        if [ $? -eq 0 ]; then
            echo "   ✅ $example builds successfully"
        else
            echo "   ❌ $example has build errors"
        fi
    fi
done
cd ..

echo ""
echo "📊 Summary:"
if [ $test_result -eq 0 ] && [ $build_result -eq 0 ]; then
    echo "✅ Go CivitAI SDK structure is valid and ready for use!"
    echo "📖 Module: $(grep '^module' go.mod)"
    echo "🏷️  Version: $(grep '^go' go.mod)"
else
    echo "⚠️  Some issues found. Check the output above."
fi

echo ""
echo "📋 Quick usage:"
echo "   go get github.com/regiellis/go-civitai-sdk"
echo "   import \"github.com/regiellis/go-civitai-sdk\""
