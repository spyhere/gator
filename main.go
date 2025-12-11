package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/spyhere/gator/internal/config"
	"github.com/spyhere/gator/internal/database"
)

func main() {
	c, err := config.ReadConfigJSON()
	if err != nil {
		log.Fatal(err)
	}
	db, err := sql.Open("postgres", c.DbURL)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)
	state := state{
		cfg: &c,
		db:  dbQueries,
	}
	commands := commands{
		all: map[string]commandHandler{},
	}
	registerCommands(&commands)

	newCommand, err := extractCommand(os.Args)
	if err != nil {
		log.Fatal(err)
	}
	if err := commands.run(&state, newCommand); err != nil {
		log.Fatal(err)
	}
}
