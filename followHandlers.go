package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/QueenQ14/gator/internal/database"
)

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("Usage: %s <url>", cmd.name)
	}
	url := cmd.args[0]
	feed, err := s.db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return fmt.Errorf("error fetching feed: %v", err)
	}

	params := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}
	row, err := s.db.CreateFeedFollow(context.Background(), params)
	if err != nil {
		return fmt.Errorf("%s encountered error following feed:  %v", user.Name, feed.Name)
	}
	fmt.Println(row.UserName, "is now following:", row.FeedName)
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	rows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("error fetching from feed-follow: %v", err)
	}
	fmt.Println(user.Name, "is following", len(rows), "feeds")
	for _, row := range rows {
		fmt.Println(" * ", row.FeedName)
	}
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("Usage: %s <url>", cmd.name)
	}
	feedURL := cmd.args[0]
	feed, err := s.db.GetFeedByUrl(context.Background(), feedURL)
	if err != nil {
		return fmt.Errorf("error fetching feed: %v", err)
	}

	params := database.DeleteFeedFollowForUserFeedParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}
	err = s.db.DeleteFeedFollowForUserFeed(context.Background(), params)
	if err != nil {
		return fmt.Errorf("error unfollowing feed: %v", err)
	}
	fmt.Println(user.Name, "has successfully unfollowed feed:", feed.Name)
	return nil
}
