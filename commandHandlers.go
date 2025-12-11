package main

import "fmt"

func handlerLogin(state *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("error: expecting [username] argument for login command!")
	}
	username := cmd.args[0]
	err := state.c.SetUser(username)
	if err != nil {
		return err
	}
	fmt.Printf("New user '%s' has been set!\n", username)
	return nil
}
