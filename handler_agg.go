package main

import (
	"context"
	"fmt"
)

func handlerAggregator(s *state, cmd command) error {
	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("could not fetch feed: %w", err)
	}
	fmt.Printf("feed: %v\n", feed)
	return nil
}
