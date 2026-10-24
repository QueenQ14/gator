package main

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/QueenQ14/gator/internal/database"
	"github.com/QueenQ14/gator/internal/rss"
	"github.com/google/uuid"
)

func handlerAggregator(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("Usage: %s <time>", cmd.name)
	}
	time_between_reqs := cmd.args[0]
	timeBetweenRequests, err := time.ParseDuration(time_between_reqs)
	if err != nil {
		return fmt.Errorf("please enter valid time value: %v", err)
	}
	ticker := time.NewTicker(timeBetweenRequests)
	fmt.Printf("Collecting feeds every %v...\n", timeBetweenRequests)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
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
		publish_date, err := time.Parse(time.RFC1123Z, item.PubDate)
		nt := sql.NullTime{}
		if err != nil {
			nt = sql.NullTime{Valid: false}
		}
		nt = sql.NullTime{
			Time:  publish_date,
			Valid: true,
		}
		desc := sql.NullString{Valid: false}
		if item.Description != "" {
			desc = sql.NullString{
				String: item.Description,
				Valid:  true,
			}
		}
		params := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Url:         item.Link,
			Description: desc,
			PublishedAt: nt,
			FeedID:      feed.ID,
		}
		// Creating the Post
		_, err = s.db.CreatePost(context.Background(), params)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique constraint") {
				continue
			}
			fmt.Printf("couldn't create post: %v", err)
			continue
		}
	}
	fmt.Printf("Feed %s collected, %d posts found\n", feed.Name, len(rssFeed.Channel.Item))
	return nil
}
