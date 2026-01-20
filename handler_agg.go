package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/zombfeed/GoBlogAggregator/internal/database"
)

func handlerAggregator(s *state, cmd command) error {
	if len(cmd.Args) < 1 || len(cmd.Args) > 2 {
		return fmt.Errorf("usage: %s <time_between_reqs> (1s, 1m, 1h, etc...)", cmd.Name)
	}
	timeBetweenRequests, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("invalid duration: %w", err)
	}
	log.Printf("Collecting feeds every %s...", timeBetweenRequests)
	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func scrapeFeeds(s *state) {
	feedFetched, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		log.Printf("Could not get feeds to fetch: %v", err)
		return
	}
	log.Println("Found a feed to fetch!")
	scrapeFeed(s.db, feedFetched)
}

func scrapeFeed(db *database.Queries, feed database.Feed) {
	_, err := db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("Could not mark feed %s as fetched: %v", feed.Name, err)
		return
	}

	feedData, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		log.Printf("Couled not collect feed %s: %v", feed.Name, err)
		return
	}
	for _, item := range feedData.Channel.Item {
		fmt.Printf("Found post: %s\n", item.Title)
	}
	log.Printf("Feed %s collected, %v posts found", feed.Name, len(feedData.Channel.Item))
}
