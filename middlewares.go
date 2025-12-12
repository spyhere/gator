package main

import (
	"context"
	"fmt"
)

func middlewareLoggedIn(handler commandHandlerLoggedIn) commandHandler {
	return func(s *state, c command) error {
		username := s.cfg.CurrentUserName
		user, err := s.db.GetUser(context.Background(), username)
		if err != nil {
			fmt.Printf("No users found for '%s'! Register this user first!\n", username)
			return err
		}
		return handler(s, c, user)
	}
}
