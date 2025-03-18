# Monkey Language Interpreter as WebAssembly

This is a WebAssembly build of the Monkey programming language interpreter. It allows you to run Monkey code directly in a web browser.

## Building the WASM Module

### Local Build

To build the WebAssembly module locally, you need to have Go installed with WebAssembly support.

```bash
# Set the GOOS and GOARCH environment variables for WebAssembly
GOOS=js GOARCH=wasm go build -tags "js,wasm" -o interpreter.wasm ./wasm_main.go
```

### Automated Build with GitHub Actions

This project includes a GitHub Actions workflow that automatically builds the WebAssembly module and creates a release with the necessary files whenever changes are pushed to the main branch.

The workflow:

1. Sets up a Go environment
2. Builds the WebAssembly module
3. Copies the required wasm_exec.js file
4. Creates a GitHub release with both files

You can find the latest release in the GitHub repository's Releases section.

## Using the WASM Module in a Web Project

1. Download the latest `interpreter.wasm` and `wasm_exec.js` files from the GitHub Releases.

2. Include the following HTML and JavaScript in your web page:

```html
<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8" />
    <title>Monkey Interpreter</title>
    <script src="wasm_exec.js"></script>
    <script>
      const go = new Go();
      WebAssembly.instantiateStreaming(
        fetch("interpreter.wasm"),
        go.importObject
      ).then((result) => {
        go.run(result.instance);
        console.log("WASM loaded successfully");
      });

      function runMonkeyCode() {
        const code = document.getElementById("code").value;
        const result = goInterpret(code);
        document.getElementById("result").textContent = result;
      }
    </script>
  </head>
  <body>
    <h1>Monkey Language Interpreter</h1>
    <textarea id="code" rows="10" cols="50">
let x = 5;
x + 10;</textarea
    >
    <br />
    <button onclick="runMonkeyCode()">Run</button>
    <pre id="result"></pre>
  </body>
</html>
```

## Monkey Language Examples

Here are some examples of Monkey code you can run:

```
// Variables and arithmetic
let x = 5;
let y = 10;
x + y;

// Functions
let add = fn(a, b) { a + b; };
add(5, 5);

// Arrays
let arr = [1, 2, 3, 4, 5];
arr[2];

// Hashes (dictionaries)
let hash = {"name": "Monkey", "age": 1};
hash["name"];

// Conditionals
if (5 > 3) {
    return "greater";
} else {
    return "less";
}
```

## Troubleshooting

If you encounter the error "could not import syscall/js", make sure you're building with the correct environment variables and build tags:

```bash
GOOS=js GOARCH=wasm go build -tags "js,wasm" -o interpreter.wasm ./wasm_main.go
```

The `syscall/js` package is only available when building for WebAssembly. The build tags (`js,wasm`) ensure that the code is only compiled when targeting WebAssembly.

For more information on Go's WebAssembly support, see the [official documentation](https://github.com/golang/go/wiki/WebAssembly).
