package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/QueenQ14/gator/internal/database"
	"github.com/QueenQ14/gator/internal/rss"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("please enter a valid username")
	}
	val, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("user does not exist")
	}
	err = s.Config.SetUser(val.Name)
	if err != nil {
		return err
	}
	fmt.Println("The user has been set to: ", cmd.args[0])
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("please provide a name")
	}
	name := cmd.args[0]
	params := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
	}
	val, err := s.db.CreateUser(context.Background(), params)
	if err != nil {
		return fmt.Errorf("user already exists")
	}
	fmt.Println(val.ID, " | ", val.CreatedAt, " | ", val.UpdatedAt, " | ", val.Name)
	return s.Config.SetUser(val.Name)
}

func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		return err
	}
	fmt.Println("the users have been reset")
	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}
	for _, user := range users {
		user_str := "* " + user.Name
		if user.Name == s.Config.UserName {
			user_str += " (current)"
		}
		fmt.Println(user_str)
	}
	return nil
}

func handlerAggregator(s *state, cmd command) error {
	feedURL := "https://www.wagslane.dev/index.xml"
	feed, err := rss.FetchFeed(context.Background(), feedURL)
	if err != nil {
		return err
	}
	fmt.Printf("Feed: %v\n", feed)
	return nil
}

func handlerAddfeed(s *state, cmd command) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("please enter valid arguments")
	}
	title := cmd.args[0]
	feedURL := cmd.args[1]
	current_user, err := s.db.GetUser(context.Background(), s.Config.UserName)
	if err != nil {
		return fmt.Errorf("error fetching user: %v", err)
	}
	params := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      title,
		Url:       feedURL,
		UserID:    current_user.ID,
	}
	feed, err := s.db.CreateFeed(context.Background(), params)
	if err != nil {
		return fmt.Errorf("error creating feed: %v", err)
	}
	fmt.Printf("Feed: %v\n", feed)
	return nil
}

func handlerListFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("error fetching feeds: %v", err)
	}
	for _, feed := range feeds {
		user, err := s.db.GetUserByID(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("error retrieving user for feed: %v", feed.Name)
		}
		printfeed(feed, user)
	}
	return nil
}
