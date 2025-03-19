//go:build js && wasm
// +build js,wasm

// wasm_main.go
// package main

// import (
// 	"strings"
// 	"syscall/js"

// 	"github.com/MahmoudHilani/GoInterpreter/test/src/monkey/evaluator"
// 	"github.com/MahmoudHilani/GoInterpreter/test/src/monkey/lexer"
// 	"github.com/MahmoudHilani/GoInterpreter/test/src/monkey/object"
// 	"github.com/MahmoudHilani/GoInterpreter/test/src/monkey/parser"
// )

// func main() {
// 	c := make(chan struct{})
	
// 	// Register JavaScript functions
// 	js.Global().Set("goInterpret", js.FuncOf(interpret))
	
// 	<-c
// }

// func interpret(this js.Value, args []js.Value) interface{} {
// 	if len(args) < 1 {
// 		return "Error: No code provided"
// 	}
	
// 	code := args[0].String()
	
// 	// Set up the environment
// 	env := object.NewEnvironment()
	
// 	// Create lexer and parser
// 	l := lexer.New(code)
// 	p := parser.New(l)
	
// 	// Parse the program
// 	program := p.ParseProgram()
	
// 	// Check for parser errors
// 	if len(p.Errors()) != 0 {
// 		errorMessages := strings.Join(p.Errors(), "\n")
// 		return "Parser Errors:\n" + errorMessages
// 	}
	
// 	// Evaluate the program
// 	evaluated := evaluator.Eval(program, env)
	
// 	// Return the result
// 	if evaluated != nil {
// 		return evaluated.Inspect()
// 	}
	
// 	return "null"
// }