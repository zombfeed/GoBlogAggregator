package main

import "fmt"

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}
	err := s.config.SetUser(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("unable to set current user: %w", err)
	}
	fmt.Printf("User switched successfully!")
	return nil
}
