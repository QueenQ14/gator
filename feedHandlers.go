package main

import (
	"context"
	"fmt"
	"time"

	"github.com/QueenQ14/gator/internal/database"
	"github.com/google/uuid"
)

func handlerAddfeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("Usage: %s \"<name>\" \"<url>\"", cmd.name)
	}
	title := cmd.args[0]
	feedURL := cmd.args[1]
	params := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      title,
		Url:       feedURL,
		UserID:    user.ID,
	}
	feed, err := s.db.CreateFeed(context.Background(), params)
	if err != nil {
		return fmt.Errorf("error creating feed: %v", err)
	}
	row, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("error creating feed-follow: %v", err)
	}
	fmt.Println(row)
	return nil
}

func handlerDelFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("Usage: %s <feed_url>", cmd.name)
	}
	url := cmd.args[0]
	rowsAffected, err := s.db.DeleteFeedByUrl(context.Background(), url)
	if err != nil {
		return fmt.Errorf("error deleting feed: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("feed '%s' does not exist", url)
	}
	fmt.Printf("feed '%s' has been deleted successfully\n", url)
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
		printFeed(feed, user)
	}
	return nil
}
