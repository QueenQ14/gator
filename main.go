package main

import (
	"fmt"
	"log"
	"os"

	"github.com/QueenQ14/gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	current_state := &state{
		Config: &cfg,
	}
	cmds := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handlerLogin)
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
