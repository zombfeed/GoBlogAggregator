package main

import (
	"context"
	"fmt"
)

func handlerReset(s *state, cmd command) error {
	if len(cmd.Args) > 0 {
		return fmt.Errorf("usage: %s", cmd.Name)
	}
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		return fmt.Errorf("unable to reset user database: %w", err)
	}
	fmt.Println("User database succcesfully reset!")
	return nil
}
