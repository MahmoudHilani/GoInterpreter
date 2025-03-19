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
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)	
	}



	var req InterpretRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	env := object.NewEnvironment()
	result := repl.StartAPI(req.Code, env)

	resp := InterpretResponse{
		Result: result,
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

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




