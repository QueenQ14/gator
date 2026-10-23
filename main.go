package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/QueenQ14/gator/internal/config"
	"github.com/QueenQ14/gator/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	// Reading config
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	// SQL DB Connection Setup and initialising SQLC DB Queries
	dbUrl := cfg.DBUrl
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatalf("error connecting with db: %v", err)
	}
	dbQueries := database.New(db)

	// Setting current state
	current_state := &state{
		db:     dbQueries,
		Config: &cfg,
	}

	// Setting up command registrations
	cmds := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	main_args := os.Args
	if len(main_args) < 2 {
		fmt.Println("not enough arguments were provided")
		os.Exit(1)
	}
	cmd := command{name: main_args[1], args: main_args[2:]}
	// fmt.Println("cmdname: ", cmd.name, "|| args: ", cmd.args)
	err = cmds.run(current_state, cmd)
	if err != nil {
		log.Fatal(err)
	}
}
