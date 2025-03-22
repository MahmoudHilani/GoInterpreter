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

func handleInterpret(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers for all responses
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	
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
	
	var req InterpretRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Printf("JSON decode error: %v", err)
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("Received code: '%s'", req.Code)
	
	env := object.NewEnvironment()
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