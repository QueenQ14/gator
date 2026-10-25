package main

import (
	"fmt"

	"github.com/QueenQ14/gator/internal/database"
)

func printFeed(feed database.Feed, user database.User) {
	fmt.Println("Feed:")
	fmt.Println(" * ID: ", feed.ID)
	fmt.Println(" * Name: ", feed.Name)
	fmt.Println(" * URL: ", feed.Url)
	fmt.Println(" * User: ", user.Name)
	fmt.Println(" * Created: ", feed.CreatedAt)
	fmt.Println(" * Updated: ", feed.UpdatedAt)
	fmt.Println(" * Last Fetched: ", feed.LastFetchedAt)
	fmt.Println("===============================")
}

func registerAllCommands(cmds *commands) {
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAggregator)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddfeed))
	cmds.register("delfeed", middlewareLoggedIn(handlerDelFeed))
	cmds.register("feeds", handlerListFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	cmds.register("browse", middlewareLoggedIn(handlerBrowse))
}
