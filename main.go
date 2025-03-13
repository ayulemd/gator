package main

import (
	"log"
	"os"

	"github.com/ayulemd/gator/internal/config"
)

type state struct {
	cfg *config.Config
}

func main() {
	gatorConfig, err := config.Read()
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	appState := state{&gatorConfig}
	commandsMap := make(map[string]func(*state, command) error)

	appCommands := commands{commandsMap}
	appCommands.register("login", handlerLogin)

	cliArgs := os.Args

	if len(cliArgs) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
	}

	cmd := command{cliArgs[1], cliArgs[2:]}

	err = appCommands.run(&appState, cmd)
	if err != nil {
		log.Fatal(err)
	}
}
