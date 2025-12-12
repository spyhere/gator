package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/spyhere/gator/internal/database"
	"github.com/spyhere/gator/internal/rss"
)

func handlerLogin(state *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("error: expecting [username] argument for login command!")
	}
	username := cmd.args[0]
	_, err := state.db.GetUser(context.Background(), username)
	if err != nil {
		fmt.Printf("No users found for '%s'! Register this user first!\n", username)
		return err
	}
	err = state.cfg.SetUser(username)
	if err != nil {
		return err
	}
	fmt.Printf("New user '%s' has been set!\n", username)
	return nil
}

func handlerRegister(state *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("error: expecting [username] argument for register command!")
	}
	username := cmd.args[0]
	now := time.Now()
	newUser := database.CreateUserParams{
		ID:        uuid.New(),
		Name:      username,
		CreatedAt: now,
		UpdatedAt: now,
	}
	user, err := state.db.CreateUser(context.Background(), newUser)
	if err != nil {
		fmt.Printf("User '%s' already exists!\n", username)
		return err
	}
	state.cfg.SetUser(user.Name)
	fmt.Println("New user has been created.")
	fmt.Println(user.ID)
	fmt.Println(user.Name)
	fmt.Println(user.UpdatedAt)
	fmt.Println(user.CreatedAt)
	return nil
}

func handleReset(state *state, _ command) error {
	err := state.db.TruncateUsers(context.Background())
	if err != nil {
		return err
	}
	fmt.Println("Successfuly reseted the 'users' table.")
	return nil
}

func handleUsers(state *state, _ command) error {
	users, err := state.db.GetUsers(context.Background())
	if err != nil {
		return err
	}
	res := ""
	for _, it := range users {
		isCurrent := ""
		if it.Name == state.cfg.CurrentUserName {
			isCurrent = "(current)"
		}
		res += fmt.Sprintf("* %s %s\n", it.Name, isCurrent)
	}
	fmt.Print(res)
	return nil
}

func handleAgg(state *state, _ command) error {
	rssFeed, err := rss.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}
	fmt.Printf("%v\n", *rssFeed)
	return nil
}
