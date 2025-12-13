package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
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
	fmt.Println("Successfully reset the 'users' table.")
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

func handleAgg(state *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("Expecting [time between requests] argument")
	}
	timeBetweenReq, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return fmt.Errorf("Incorrect time interval!\nExample: 300ms, 5m, 1h")
	}
	fmt.Println("Collecting feeds every", timeBetweenReq)
	ticker := time.NewTicker(timeBetweenReq)
	for ; ; <-ticker.C {
		if err = scrapeFeeds(state); err != nil {
			return err
		}
	}
}

func handleAddFeed(state *state, cmd command, user database.User) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("Expecting 2 arguments: [name] and [url], instead got %v\n", cmd.args)
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
	feedFollowParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    user.ID,
		FeedID:    feed.ID,
	}
	_, err = state.db.CreateFeedFollow(context.Background(), feedFollowParams)
	if err != nil {
		return err
	}
	body := [][]string{{feed.Name, feed.Url, feed.UpdatedAt.Format(TIME_FORMAT)}}
	res, err := createStringifiedTable([]string{"Name", "Url", "Updated At"}, body)
	if err != nil {
		return err
	}
	fmt.Println(res)
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
	res, err := createStringifiedTable([]string{"Name", "URL", "Creator"}, content)
	if err != nil {
		return err
	}
	fmt.Println(res)
	return nil
}

func handleFollow(state *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("Expecting 1 argument: [url]")
	}
	ctx := context.Background()
	url := cmd.args[0]
	feed, err := state.db.GetFeedByUrl(ctx, url)
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
	res, err := createStringifiedTable([]string{"User", "Feed", "Updated At"}, body)
	if err != nil {
		return err
	}
	fmt.Println(res)
	return nil
}

func handleFollowing(state *state, _ command, user database.User) error {
	feeds, err := state.db.GetUserFeeds(context.Background(), user.ID)
	if err != nil {
		return err
	}
	body := [][]string{}
	for _, it := range feeds {
		body = append(body, []string{it.FeedName, it.Url, it.UpdatedAt.Format(TIME_FORMAT)})
	}
	res, err := createStringifiedTable([]string{"Feed", "Url", "Updated At"}, body)
	if err != nil {
		return err
	}
	fmt.Println(res)
	return nil
}

func handleUnfollow(state *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("Expecting [url] argument.")
	}
	ctx := context.Background()
	url := cmd.args[0]
	feed, err := state.db.GetFeedByUrl(ctx, url)
	if err != nil {
		return err
	}
	deleteUserFeedParams := database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}
	if err = state.db.DeleteFeedFollow(ctx, deleteUserFeedParams); err != nil {
		return fmt.Errorf("You are not following current feed")
	}
	fmt.Println("You have successfully unfollowed", url)
	return nil
}

func scrapeFeeds(state *state) error {
	ctx := context.Background()
	feed, err := state.db.GetNextFeedToFetch(ctx)
	if err != nil {
		return err
	}
	rssFeed, err := rss.FetchFeed(ctx, feed.Url)
	if err != nil {
		return err
	}
	err = state.db.MarkFeedFetched(ctx, feed.ID)
	fmt.Println(rssFeed.Channel.Title)
	for _, it := range rssFeed.Channel.Item {
		pubDate, err := time.Parse("Mon, 02 Jan 2006 15:04:05 -0700", it.PubDate)
		if err != nil {
			return fmt.Errorf("Error parsing:\n%s", err.Error())
		}
		_, err = state.db.CreatePost(ctx, database.CreatePostParams{
			ID:          uuid.New(),
			Title:       it.Title,
			Url:         it.Link,
			Description: sql.NullString{String: it.Description, Valid: true},
			FeedID:      feed.ID,
			PublishedAt: pubDate,
		})
		var pqError *pq.Error
		if errors.As(err, &pqError) {
			if pqError.Code.Name() != "unique_violation" {
				return pqError
			}
		}
	}
	log.Printf("Grabbed feed for %s\nTotal titles: %d\n", rssFeed.Channel.Title, len(rssFeed.Channel.Item))
	return nil
}
