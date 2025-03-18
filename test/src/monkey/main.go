package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/MahmoudHilani/GoInterpreter/test/src/monkey/repl"
)

func mmain() {
	user, err := user.Current()
	if err != nil {
		panic(err)

	}
	fmt.Printf("Hello %s! Start yapping\n", user.Username)
	repl.Start(os.Stdin, os.Stdout)
}