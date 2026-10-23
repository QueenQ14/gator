package main

import (
	"fmt"
)

type command struct {
	name string
	args []string
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("please enter a valid username")
	}
	err := s.Config.SetUser(cmd.args[0])
	if err != nil {
		return err
	}
	fmt.Println("The user has been set to: ", cmd.args[0])
	return nil
}
