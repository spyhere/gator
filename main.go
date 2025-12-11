package main

import (
	"log"
	"os"

	"github.com/spyhere/gator/internal/config"
)

func main() {
	c, err := config.ReadConfigJSON()
	if err != nil {
		log.Fatal(err)
	}
	state := state{
		c: &c,
	}
	commands := commands{
		all: map[string]commandHandler{},
	}
	commands.register("login", handlerLogin)

	newCommand, err := extractCommand(os.Args)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	if err := commands.run(&state, newCommand); err != nil {
		log.Fatal(err)
	}
}
