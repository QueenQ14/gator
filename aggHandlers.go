package main

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/QueenQ14/gator/internal/database"
	"github.com/QueenQ14/gator/internal/rss"
	"github.com/google/uuid"
)

func handlerAggregator(s *state, cmd command) error {
	if len(cmd.args) < 1 || len(cmd.args) > 2 {
		return fmt.Errorf("Usage: %s <time> [concurrency] <- default=3", cmd.name)
	}
	time_between_reqs := cmd.args[0]
	timeBetweenRequests, err := time.ParseDuration(time_between_reqs)
	if err != nil {
		return fmt.Errorf("please enter valid time value: %v", err)
	}
	concurrency := 3
	if len(cmd.args) == 2 {
		parsedConcurrency, err := strconv.Atoi(cmd.args[1])
		if err != nil {
			return fmt.Errorf("invalid concurrency value: %v", err)
		}
		if parsedConcurrency <= 0 {
			return fmt.Errorf("concurrency value must be greater than 0")
		}
		concurrency = parsedConcurrency
	}
	ticker := time.NewTicker(timeBetweenRequests)
	defer ticker.Stop()
	fmt.Printf("Collecting feeds every %v with %d workers...\n", timeBetweenRequests, concurrency)
	for ; ; <-ticker.C {
		fmt.Println("refreshing...")
		feeds, err := s.db.GetNextFeedsToFetch(context.Background(), int32(concurrency))
		if err != nil {
			fmt.Println("Error fetching feeds:", err)
			continue
		}

		var wg sync.WaitGroup
		for _, feed := range feeds {
			wg.Add(1)
			go func(f database.Feed) {
				defer wg.Done()
				if err := scrapeFeed(s, f); err != nil {
					fmt.Printf("Error scraping %s: %v\n", f.Name, err)
				}
			}(feed)
		}
		wg.Wait()
	}
}

func scrapeFeed(s *state, feed database.Feed) error {
	err := s.db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		return fmt.Errorf("error marking the feed as fetched: %v", err)
	}
	rssFeed, err := rss.FetchFeed(context.Background(), feed.Url)
	if err != nil {
		return fmt.Errorf("error fetching the feed: %v", err)
	}
	for _, item := range rssFeed.Channel.Item {
		publish_date, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			publish_date, err = time.Parse(time.RFC1123, item.PubDate)
		}
		nt := sql.NullTime{}
		if err == nil {
			nt = sql.NullTime{
			Time:  publish_date,
			Valid: true,
			}
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
			fmt.Printf("couldn't create post: %v\n", err)
			continue
		}
	}
	fmt.Printf("Feed %s collected, %d posts found\n", feed.Name, len(rssFeed.Channel.Item))
	return nil
}
