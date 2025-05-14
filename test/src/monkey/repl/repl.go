package repl

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/MahmoudHilani/GoInterpreter/test/src/monkey/evaluator"
	"github.com/MahmoudHilani/GoInterpreter/test/src/monkey/lexer"
	"github.com/MahmoudHilani/GoInterpreter/test/src/monkey/object"
	"github.com/MahmoudHilani/GoInterpreter/test/src/monkey/parser"
	"github.com/gorilla/websocket"
)

const PROMPT = ">> "
type InterpretRequest struct {
	Code string `json:"code"`
}

type InterpretResponse struct {
	Result string `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}


func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	for {
		fmt.Print(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func StartAPI(c *websocket.Conn)  {
	env := object.NewEnvironment()
	var res string 
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, 
				websocket.CloseGoingAway, 
				websocket.CloseNormalClosure) {
				log.Printf("WebSocket read error: %v", err)
			}
			break
		}
		log.Printf("WebSocket message received: %d, %s", mt, message)

		var req InterpretRequest
		
		err = json.Unmarshal(message, &req)
		if err != nil {
		errorResponse := createErrorResponse(err)
		log.Printf("JSON parse error: %v", err)
		
		// Send the error response back to the client
		err = c.WriteJSON(errorResponse)
		if err != nil {
			log.Printf("Error sending error response: %v", err)
			break
		}
		continue
	}
		l := lexer.New(req.Code)
		p := parser.New(l)
		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			 printParserErrorsAPI(p.Errors(), c)
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			res = evaluated.Inspect()
		}
		if evaluated == nil {
			res = "continue"
		}
		resp := InterpretResponse{
			Result: res,
		}
		err = c.WriteJSON(resp)
		if err != nil {
			log.Printf("Error writing response: %v", err)
			break
		}
	}
	
}

const MONKEY_FACE = `
              __,__
      .---. .-"   "-. .---.
     / .. \/ .-. .-. \/ .. \
    |  |  ' | / Y \ | '  |  |
    |  \  \ \ 0 | 0 / /  /  |
     \'- ,\.-"""""""-./, -'/
      ' '-' /_ ^ ^ _\ '-' '
           | \._ _./ |
           \ \ '~' / /
           '._'-=-'_.'
             '-----'
`
func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, MONKEY_FACE)
	io.WriteString(out, "Woops! We ran into some monkey business here!\n")
	io.WriteString(out, " parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}

func printParserErrorsAPI(errors []string, c *websocket.Conn) {
	e := strings.Join(errors, "\r\n")
	monkeyFace := strings.ReplaceAll(MONKEY_FACE, "\n", "\r\n")
	out := monkeyFace + "Woops! We ran into some monkey business here!\r\n" + e
	resp := InterpretResponse{
		Result: out,
	}
	err := c.WriteJSON(resp)
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

func createErrorResponse(err error) map[string]interface{} {
	var errorResponse map[string]interface{}
	var syntaxError *json.SyntaxError
	var unmarshalError *json.UnmarshalTypeError
	switch {
	case errors.As(err, &syntaxError):
		errorResponse = map[string]interface{}{
			"error": fmt.Sprintf("Request contains badly-formed JSON (at position %d)", syntaxError.Offset),
			"code":  "bad_request",
		}
	
	case errors.Is(err, io.ErrUnexpectedEOF):
		errorResponse = map[string]interface{}{
			"error": "Request contains badly-formed JSON",
			"code":  "bad_request",
		}
	
	case errors.As(err, &unmarshalError):
		errorResponse = map[string]interface{}{
			"error": fmt.Sprintf("Request contains an invalid value for the %q field (at position %d)", 
				unmarshalError.Field, unmarshalError.Offset),
			"code":  "bad_request",
		}
	
	case strings.HasPrefix(err.Error(), "json: unknown field "):
		fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
		errorResponse = map[string]interface{}{
			"error": fmt.Sprintf("Request contains unknown field: %s", fieldName),
			"code":  "bad_request",
		}
	
	case errors.Is(err, io.EOF):
		errorResponse = map[string]interface{}{
			"error": "Request must not be empty",
			"code":  "bad_request",
		}
	
	default:
		errorResponse = map[string]interface{}{
			"error": "Internal server error",
			"code":  "server_error",
		}
		log.Printf("Unhandled error: %v", err)
	}
	
	return errorResponse
	}

