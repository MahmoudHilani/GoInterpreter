// code for local machine
// package main

// import (
// 	"fmt"
// 	"os"
// 	"os/user"

// 	"github.com/MahmoudHilani/GoInterpreter/test/src/monkey/repl"
// )

// func main() {
// 	user, err := user.Current()
// 	if err != nil {
// 		panic(err)

// 	}
// 	fmt.Printf("Hello %s! Start yapping\n", user.Username)
// 	repl.Start(os.Stdin, os.Stdout)
// }

package main

import (
	"log"
	"net/http"
	"os"

	"github.com/MahmoudHilani/GoInterpreter/test/src/monkey/repl"
	"github.com/gorilla/websocket"
)

type InterpretRequest struct {
	Code string `json:"code"`
}

type InterpretResponse struct {
	Result string `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}

var upgrader = websocket.Upgrader{}

func handleSocket(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	upgrader.ReadBufferSize = 1024
	upgrader.WriteBufferSize = 1024
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Websocket upgrade error: %v", err)
		return
	}
	defer c.Close()
	c.SetReadLimit(1048576)
	repl.StartAPI(c)
}



func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/socket", handleSocket)

	log.Printf("Starting server on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}