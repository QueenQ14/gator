package main

import (
	"context"
	"fmt"

	"github.com/QueenQ14/gator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(),s.Config.UserName)
		if err != nil {
			return fmt.Errorf("error fetching user: %v", err)
		}
		return handler(s, cmd, user)
	}
	
}
