package main

import (
	"fmt"
	"os"

	"github.com/zombfeed/GoBlogAggregator/internal/config"
)

type handlerFunc func(*state, command) error

type state struct {
	config *config.Config
}

type command struct {
	Name string
	Args []string
}

type commands struct {
	commands map[string]handlerFunc
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("username is required")
	}
	err := s.config.SetUser(cmd.Args[0])
	if err != nil {
		return err
	}
	fmt.Printf("User has been set to %s\n", cmd.Args[0])
	return nil
}

func (c *commands) run(s *state, cmd command) error {
	toRun, ok := c.commands[cmd.Name]
	if !ok {
		return fmt.Errorf("the given command does not exist: %s", cmd.Name)
	}
	err := toRun(s, cmd)
	if err != nil {
		return err
	}
	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.commands[name] = f
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		return
	}

	s := state{&cfg}
	handlers := make(map[string]handlerFunc)
	c := commands{commands: handlers}
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
