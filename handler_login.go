package main

import (
	"errors"
	"fmt"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("username required")
	}

	err := s.cfg.SetUser(cmd.args[0])
	if err != nil {
		return errors.New("unable to set user")
	}

	fmt.Printf("User set to %s\n", cmd.args[0])

	return nil
}
