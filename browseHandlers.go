package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/QueenQ14/gator/internal/database"
)

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
		fmt.Print("\n[n]next | [p]previous | [q]quit: ")
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
