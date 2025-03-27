package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ayulemd/gator/internal/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("username required")
	}

	user, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			fmt.Println(pqErr.Code, pqErr.Message)
			os.Exit(1)
		}
	}

	if user.Name == "" {
		log.Printf("user does not exist")
		os.Exit(1)
	}

	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return errors.New("unable to set user")
	}

	fmt.Printf("User set to %s\n", user.Name)

	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("username required")
	}

	userParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
		Name:      cmd.args[0],
	}

	user, err := s.db.CreateUser(context.Background(), userParams)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				fmt.Println(pqErr.Message)
			}
		} else {
			fmt.Println("unknown error when creating user")
		}

		os.Exit(1)
	}

	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return err
	}

	log.Printf("Created user: %+v", user)

	return nil
}

func handlerResetUsers(s *state, cmd command) error {
	err := s.db.ResetUsers(context.Background())
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			fmt.Println(pqErr.Message)
		} else {
			fmt.Println("unknown error when reseting users")
		}

		os.Exit(1)
	}

	fmt.Println("Users reset")

	return nil
}

func handlerGetUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			fmt.Println(pqErr.Message)
		} else {
			fmt.Println("unknown error when listing users")
		}

		os.Exit(1)
	}

	for _, user := range users {
		if user.Name == s.cfg.CurrentUserName {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Println("*", user.Name)
		}
	}

	return nil
}

func handlerAgg(s *state, cmd command) error {
	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("error fetching feed: %w", err)
	}

	fmt.Printf("%+v\n", feed)

	return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 2 {
		return errors.New("CLI usage: gator addfeed [name] [url]")
	}

	feedParams := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
		Name:      sql.NullString{String: cmd.args[0], Valid: true},
		Url:       sql.NullString{String: cmd.args[1], Valid: true},
		UserID:    uuid.NullUUID{UUID: user.ID, Valid: true},
	}

	feed, err := s.db.CreateFeed(context.Background(), feedParams)
	if err != nil {
		return fmt.Errorf("error creating feed: %w", err)
	}

	newCmd := command{
		name: "addfeed",
		args: []string{cmd.args[1]},
	}

	handlerFollow(s, newCmd, user)

	fmt.Printf("%+v\n", feed)

	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.ListFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("error listing feeds: %w", err)
	}

	for _, feed := range feeds {
		fmt.Printf("Feed: %s\nURL: %s\nAdded by: %s\n\n", feed.FeedName.String, feed.Url.String, feed.UserName)
	}

	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return errors.New("CLI usage: gator follow [url]")
	}

	urlString := sql.NullString{
		String: cmd.args[0],
		Valid:  true,
	}

	feedId, err := s.db.FeedIdFromUrl(context.Background(), urlString)
	if err != nil {
		return fmt.Errorf("error fetching feed ID from URL: %w", err)
	}

	feedFollowParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
		UserID:    uuid.NullUUID{UUID: user.ID, Valid: true},
		FeedID:    uuid.NullUUID{UUID: feedId, Valid: true},
	}

	followedFeed, err := s.db.CreateFeedFollow(context.Background(), feedFollowParams)
	if err != nil {
		return fmt.Errorf("error following feed: %w", err)
	}

	fmt.Printf("%s followed \"%s\"\n", followedFeed.UserName, followedFeed.FeedName.String)

	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	followedFeeds, err := s.db.GetFeedFollowsForUser(context.Background(), uuid.NullUUID{UUID: user.ID, Valid: true})
	if err != nil {
		return fmt.Errorf("error listing followed feeds: %w", err)
	}

	fmt.Println("User:", user.Name)
	fmt.Println("Feeds followed:")

	for _, feed := range followedFeeds {
		fmt.Println("*", feed.FeedName.String)
	}

	return nil
}
