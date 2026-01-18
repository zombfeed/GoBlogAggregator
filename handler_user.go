package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/zombfeed/GoBlogAggregator/internal/database"
)

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}
	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		uuid.New(),
		time.Now(),
		time.Now(),
		cmd.Args[0],
	})
	if err != nil {
		return fmt.Errorf("could not create user: %w", err)
	}
	err = s.config.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("could not set current user: %w", err)
	}
	fmt.Printf("User successfully created!")
	printUser(user)
	return nil
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}
	name := cmd.Args[0]
	if _, err := s.db.GetUser(context.Background(), name); err != nil {
		return fmt.Errorf("unable to find user: %w", err)
	}
	err := s.config.SetUser(name)
	if err != nil {
		return fmt.Errorf("unable to set current user: %w", err)
	}
	fmt.Printf("User switched successfully!")
	return nil
}

func printUser(user database.User) {
	fmt.Printf("* ID:				%v\n", user.ID)
	fmt.Printf("* NAME:			%v\n", user.Name)
}
