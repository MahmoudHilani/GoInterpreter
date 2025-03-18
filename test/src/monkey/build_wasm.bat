@echo off
echo Building Monkey Interpreter WASM module...

rem Set environment variables for WebAssembly
set GOOS=js
set GOARCH=wasm

rem Build the WASM module with build tags
go build -tags "js,wasm" -o interpreter.wasm ./wasm_main.go

rem Copy the wasm_exec.js file from the Go installation
for /f "tokens=*" %%g in ('go env GOROOT') do (set GOROOT=%%g)
copy "%GOROOT%\misc\wasm\wasm_exec.js" .

echo Build complete!
echo Files created:
echo - interpreter.wasm
echo - wasm_exec.js
echo.
echo To test the WASM module, open example.html in a web browser.
echo Note: You may need to serve the files using a local web server due to browser security restrictions.
echo For example: python -m http.server 8080 