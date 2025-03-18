#!/bin/bash

# Exit on error
set -e

echo "Building Monkey Interpreter WASM module..."

# Set environment variables for WebAssembly
export GOOS=js
export GOARCH=wasm

# Build the WASM module with build tags
go build -tags "js,wasm" -o interpreter.wasm ./wasm_main.go

# Copy the wasm_exec.js file from the Go installation
cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" .

echo "Build complete!"
echo "Files created:"
echo "- interpreter.wasm"
echo "- wasm_exec.js"
echo ""
echo "To test the WASM module, open example.html in a web browser."
echo "Note: You may need to serve the files using a local web server due to browser security restrictions."
echo "For example: python -m http.server 8080" 