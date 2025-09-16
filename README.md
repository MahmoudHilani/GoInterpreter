# GoInterpreter

A minimal Go interpreter / learning project, built to reinforce language fundamentals, interpreter architecture, and integrating with a live web-front end. Deployed for interactive use on mahmoudhilani.com.

## Features

Implements a working interpreter in Go

Supports lexical analysis, parsing, AST evaluation, and basic error handling

Includes a test suite (unit tests) to validate core components

Integrated with a web server so users can input Go-like code interactively via a browser

Deployed under tmux for persistent terminal sessions, monitored via nginx for HTTP access and proxying to the GoInterpreter web service

## Tech & Skills Demonstrated

Go programming language (tooling, modules, package structure)

Language runtimes: lexing, parsing, AST walking/evaluation

Web backend integration (serving HTTP, receiving and executing code snippets safely)

DevOps / deployment skills: use of tmux, nginx reverse proxy, server setup and configuration

Testing and code correctness

## Deployment & Integration

Runs as a long-lived process inside tmux on the server, so the session remains active even when disconnected

Exposed via nginx which handles routing / TLS / HTTP requests and proxies them to the GoInterpreter backend

Interactive front end / webserver connected (websockets) to portfolio so visitors can try out code live on mahmoudhilani.com
