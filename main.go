package main

import (
	"fmt"
	"log"

	"github.com/ayulemd/gator/internal/config"
)

func main() {
	gatorConfig, err := config.Read()
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	fmt.Printf("%v\n", gatorConfig)

	gatorConfig.SetUser("ayulemd")

	gatorConfig, err = config.Read()
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	fmt.Printf("%v\n", gatorConfig)
}
