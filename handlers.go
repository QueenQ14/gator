package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/QueenQ14/gator/internal/database"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("please enter a valid username")
	}
	val, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("user does not exist")
	}
	err = s.Config.SetUser(val.Name)
	if err != nil {
		return err
	}
	fmt.Println("The user has been set to: ", cmd.args[0])
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("please provide a name")
	}
	name := cmd.args[0]
	params := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
	}
	val, err := s.db.CreateUser(context.Background(), params)
	if err != nil {
		return fmt.Errorf("user already exists")
	}
	fmt.Println(val.ID, " | ", val.CreatedAt, " | ", val.UpdatedAt, " | ", val.Name)
	return s.Config.SetUser(val.Name)
}

func handlerReset(s *state, cmd command) error {
	err := s.db.ResetUsers(context.Background())
	if err != nil {
		return err
	}
	fmt.Println("the users have been reset")
	return nil
}
