package main

import (
	"fmt"
	"os"

	"github.com/zombfeed/GoBlogAggregator/internal/config"
)

type state struct {
	config *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		return
	}

	s := state{&cfg}
	c := commands{registeredCommands: make(map[string]func(*state, command) error)}
	c.register("login", handlerLogin)

	if len(os.Args) < 2 {
		fmt.Println("error: no arguments given ")
		os.Exit(1)
	}
	cmd := command{Name: os.Args[1], Args: os.Args[2:]}
	if err := c.run(&s, cmd); err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}
