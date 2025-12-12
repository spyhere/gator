package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/spyhere/gator/internal/database"
	"github.com/spyhere/gator/internal/rss"
)

const TIME_FORMAT = "15:04:05 02-01-2006"

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
	fmt.Println(user.Name)
	fmt.Println(user.UpdatedAt.Format(TIME_FORMAT))
	return nil
}

func handleReset(state *state, _ command) error {
	err := state.db.ClearUsers(context.Background())
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

func handleAddFeed(state *state, cmd command) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("Expecting 2 arguments: [name] and [url], instead got %v\n", cmd.args)
	}
	username := state.cfg.CurrentUserName
	user, err := state.db.GetUser(context.Background(), username)
	if err != nil {
		fmt.Printf("No users found for '%s'! Register this user first!\n", username)
		return err
	}
	feedName, feedUrl := cmd.args[0], cmd.args[1]
	now := time.Now()
	newFeed := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Name:      feedName,
		Url:       feedUrl,
		UserID:    user.ID,
	}
	feed, err := state.db.CreateFeed(context.Background(), newFeed)
	if err != nil {
		return err
	}
	fmt.Printf(
		"New feed with name '%s' successfuly created at %s\n",
		feed.Name,
		feed.CreatedAt.Format(TIME_FORMAT),
	)
	return nil
}

func handleFeeds(state *state, _ command) error {
	feeds, err := state.db.GetFeeds(context.Background())
	if err != nil {
		return nil
	}
	content := [][]string{}
	for _, it := range feeds {
		content = append(content, []string{it.Name, it.Url, it.Creator})
	}
	res, err := formatContentWithTitle([]string{"Name", "URL", "Creator"}, content)
	if err != nil {
		return err
	}
	fmt.Println(res)
	return nil
}

func handleFollow(state *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("Expecting 1 argument: [url]")
	}
	ctx := context.Background()
	username := state.cfg.CurrentUserName
	user, err := state.db.GetUser(ctx, username)
	if err != nil {
		return err
	}
	url := cmd.args[0]
	feed, err := state.db.GetFeed(ctx, url)
	if err != nil {
		return err
	}
	now := time.Now()
	feedFollowParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    user.ID,
		FeedID:    feed.ID,
	}
	feedFollow, err := state.db.CreateFeedFollow(ctx, feedFollowParams)
	if err != nil {
		return err
	}
	body := [][]string{{feedFollow.UserName, feedFollow.FeedName, feedFollow.UpdatedAt.Format(TIME_FORMAT)}}
	res, err := formatContentWithTitle([]string{"User", "Feed", "Updated At"}, body)
	if err != nil {
		return err
	}
	fmt.Println(res)
	return nil
}
