package main

import (
	"github.com/spyhere/gator/internal/config"
	"github.com/spyhere/gator/internal/database"
)

type state struct {
	cfg *config.Config
	db  *database.Queries
}

type command struct {
	name string
	args []string
}

type commandHandler func(*state, command) error

type commands struct {
	all map[string]commandHandler
}
