package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/ayulemd/gator/internal/config"
	"github.com/ayulemd/gator/internal/database"

	_ "github.com/lib/pq"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

func main() {
	gatorConfig, err := config.Read()
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}

	db, err := sql.Open("postgres", gatorConfig.DbUrl)
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}

	dbQueries := database.New(db)

	appState := &state{dbQueries, &gatorConfig}
	commandsMap := make(map[string]func(*state, command) error)

	appCommands := commands{commandsMap}
	appCommands.register("login", handlerLogin)
	appCommands.register("register", handlerRegister)
	appCommands.register("reset", handlerResetUsers)
	appCommands.register("users", handlerGetUsers)
	appCommands.register("agg", handlerAgg)
	appCommands.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	appCommands.register("feeds", handlerFeeds)
	appCommands.register("follow", middlewareLoggedIn(handlerFollow))
	appCommands.register("following", middlewareLoggedIn(handlerFollowing))
	appCommands.register("unfollow", middlewareLoggedIn(handlerUnfollow))

	cliArgs := os.Args

	if len(cliArgs) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
	}

	cmd := command{cliArgs[1], cliArgs[2:]}

	err = appCommands.run(appState, cmd)
	if err != nil {
		log.Fatal(err)
	}
}
