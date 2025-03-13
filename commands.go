package main

import "errors"

type command struct {
	name string
	args []string
}

type commands struct {
	commandsMap map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.commandsMap[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	handler, exists := c.commandsMap[cmd.name]
	if !exists {
		return errors.New("command does not exist")
	}

	err := handler(s, cmd)
	if err != nil {
		return err
	}

	return nil
}
