package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/zombfeed/GoBlogAggregator/internal/database"
)

func handlerAddFeed(s *state, cmd command) error {
	user, err := s.db.GetUser(context.Background(), s.config.CurrentUserName)
	if err != nil {
		return err
	}

	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: %s <feed name> <url>", cmd.Name)
	}

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      cmd.Args[0],
		Url:       cmd.Args[1],
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("could not get feed: %w", err)
	}
	printFeed(feed)
	return nil
}

func printFeed(feed database.Feed) {
	fmt.Printf("*NAME:			%s\n", feed.Name)
	fmt.Printf("*URL: 			%s\n", feed.Url)
}
