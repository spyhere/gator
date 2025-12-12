package main

import (
	"fmt"

	"github.com/spyhere/gator/internal/database"
)

type command struct {
	name string
	args []string
}

type commandHandler func(*state, command) error
type commandHandlerLoggedIn func(*state, command, database.User) error

type commands struct {
	all map[string]commandHandler
}

func (c *commands) run(s *state, cmd command) error {
	if _, ok := c.all[cmd.name]; !ok {
		return fmt.Errorf("Command '%s' is not found!", cmd.name)
	}
	return c.all[cmd.name](s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) error {
	c.all[name] = f
	return nil
}
