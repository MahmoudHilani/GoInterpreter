// package main

// import (
// 	"fmt"
// 	"os"
// 	"os/user"

// 	"github.com/MahmoudHilani/GoInterpreter/test/src/monkey/repl"
// )

// func mmain() {
// 	user, err := user.Current()
// 	if err != nil {
// 		panic(err)

// 	}
// 	fmt.Printf("Hello %s! Start yapping\n", user.Username)
// 	repl.Start(os.Stdin, os.Stdout)
// }

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"log"
	"net/http"
	"os"

	"github.com/MahmoudHilani/GoInterpreter/test/src/monkey/object"
	"github.com/MahmoudHilani/GoInterpreter/test/src/monkey/repl"
)

type InterpretRequest struct {
	Code string `json:"code"`
}

type InterpretResponse struct {
	Result string `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}

var env = object.NewEnvironment()

func handleInterpret(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers for all responses
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	ct := r.Header.Get("Content-Type")
	if ct != "" {
		mediaType := strings.ToLower(strings.TrimSpace(strings.Split(ct, ";")[0]))
		if mediaType != "application/json" {
			msg := "Content-Type header is not application/json"
			http.Error(w, msg, http.StatusUnsupportedMediaType)
			return 
		}
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)
	
	// Handle preflight OPTIONS request
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	
	if r.Method != http.MethodPost {
		log.Printf("Invalid method: %s", r.Method)
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return	
	}

	// Debug request
	log.Printf("Content-Type: %s", r.Header.Get("Content-Type"))
	log.Printf("Request length: %d", r.ContentLength)
	
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	var req InterpretRequest
	err := dec.Decode(&req)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalError *json.UnmarshalTypeError
		var maxBytesError *http.MaxBytesError

		switch {
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			http.Error(w, msg, http.StatusBadRequest)

		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := fmt.Sprint("Request body contains badly-formed JSON")
			http.Error(w, msg, http.StatusBadRequest)

		case errors.As(err, &unmarshalError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalError.Field, unmarshalError.Offset)
			http.Error(w, msg, http.StatusBadRequest)

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains uknown field: %s", fieldName)
			http.Error(w, msg, http.StatusBadRequest)
		
		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			http.Error(w, msg, http.StatusBadRequest)

		case errors.As(err, &maxBytesError):
			msg := fmt.Sprintf("Request body must not be larger than %d bytes", maxBytesError.Limit)
			http.Error(w, msg, http.StatusRequestEntityTooLarge)

		default:
			log.Print(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		}
		return
	}
	log.Printf("Request: %s", req)
	log.Printf("Received code: '%s'", req.Code)
	
	
	result := repl.StartAPI(req.Code, env)
	log.Printf("Interpretation result: %s", result)

	resp := InterpretResponse{
		Result: result,
	}
	
	json.NewEncoder(w).Encode(resp)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/interpret", handleInterpret)

	log.Printf("Starting server on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}