package main

import (
	"context"
	"fmt"

	"github.com/QueenQ14/gator/internal/database"
	"github.com/QueenQ14/gator/internal/rss"
)

func printFeed(feed database.Feed, user database.User) {
	fmt.Println("Feed:")
	fmt.Println(" * ID: ", feed.ID)
	fmt.Println(" * Name: ", feed.Name)
	fmt.Println(" * URL: ", feed.Url)
	fmt.Println(" * User: ", user.Name)
	fmt.Println(" * Created: ", feed.CreatedAt)
	fmt.Println(" * Updated: ", feed.UpdatedAt)
	fmt.Println("===============================")
}

func scrapeFeeds(s *state) error {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("error fetching next feed: %v", err)
	}
	err = s.db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		return fmt.Errorf("error marking the feed as fetched: %v", err)
	}
	rssFeed, err := rss.FetchFeed(context.Background(), feed.Url)
	if err != nil {
		return fmt.Errorf("error fetching the feed: %v", err)
	}
	for _, item := range rssFeed.Channel.Item {
		fmt.Println(item.Title)
	}
	return nil
}
