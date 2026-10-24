package main

import (
	"fmt"

	"github.com/QueenQ14/gator/internal/database"
)

func printfeed(feed database.Feed, user database.User) {
	fmt.Println("Feed:")
	fmt.Println(" * ID: ", feed.ID)
	fmt.Println(" * Name: ", feed.Name)
	fmt.Println(" * URL: ", feed.Url)
	fmt.Println(" * User: ", user.Name)
	fmt.Println(" * Created: ", feed.CreatedAt)
	fmt.Println(" * Updated: ", feed.UpdatedAt)
	fmt.Println("===============================")
}
