package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
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

func handlerBrowse(s *state, cmd command, user database.User) error {
	limit := 2
	if len(cmd.args) > 0 {
		parsedLimit, err := strconv.Atoi(cmd.args[0])
		if err != nil {
			return fmt.Errorf("Usage: %s [limit] <- default=2", cmd.name)
		}
		limit = parsedLimit
	}
	page := 1
	scanner := bufio.NewScanner(os.Stdin)
	for {
		offset := (page - 1) * limit
		posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
			UserID: user.ID,
			Limit:  int32(limit),
			Offset: int32(offset),
		})
		if err != nil {
			return fmt.Errorf("error retrieving posts: %v", err)
		}
		if len(posts) == 0 {
			if page == 1 {
				fmt.Println("no posts found")
				return nil
			}
			fmt.Println("\nno more posts available")
			page--
			continue
		}
		fmt.Printf("\n================= Page %d ================\n\n", page)
		for _, post := range posts {
			fmt.Printf("%s from %s\n", post.PublishedAt.Time.Format("Mon Jan 2"), post.FeedName)
			fmt.Printf("--- %s ---\n", post.Title)
			fmt.Printf("    %v\n", post.Description.String)
			fmt.Printf("Link: %s\n", post.Url)
			fmt.Println("========================================")
		}
		fmt.Print("\n[n]ext | [p]revious | [q]quit: ")
		if !scanner.Scan() {
			break
		}
		choice := strings.TrimSpace(strings.ToLower(scanner.Text()))
		switch choice {
		case "n", "next":
			if len(posts) < limit {
				fmt.Println("you are on the last page")
			} else {
				page++
			}
		case "p", "prev":
			if page > 1 {
				page--
			} else {
				fmt.Println("you are on the first page")
			}
		case "q", "quit", "exit":
			fmt.Println("exiting browse...")
			return nil
		default:
			fmt.Println("invalid input.Use 'n','p' or 'q'")
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading input: %v", err)
	}
	return nil
}
