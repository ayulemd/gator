package main

import (
	"context"
	"fmt"

	"github.com/ayulemd/gator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return fmt.Errorf("error fetching current user ID: %w", err)
		}

		err = handler(s, cmd, user)
		if err != nil {
			return fmt.Errorf("error executing handler: %w", err)
		}

		return nil
	}
}
