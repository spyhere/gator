package main

import "github.com/spyhere/gator/internal/config"

type state struct {
	c *config.Config
}

type command struct {
	name string
	args []string
}

type commandHandler func(*state, command) error

type commands struct {
	all map[string]commandHandler
}
