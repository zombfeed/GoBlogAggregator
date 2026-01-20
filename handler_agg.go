package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
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
		log.Printf("Could not collect feed %s: %v", feed.Name, err)
		return
	}
	for _, item := range feedData.Channel.Item {
		savePost(db, item, feed)
		log.Printf("Feed %s collected, %v posts found", feed.Name, len(feedData.Channel.Item))
	}
}

//	func parseTimeLayouts(time_to_parse string) time.Time {
//		time_formats := []string{
//			time.Layout,
//			time.ANSIC,
//			time.UnixDate,
//			time.RubyDate,
//			time.RFC822,
//			time.RFC822Z,
//			time.RFC850,
//			time.RFC1123,
//			time.RFC1123Z,
//			time.RFC3339,
//			time.RFC3339Nano,
//			time.Kitchen,
//		}
//
//		for _, layout := range time_formats {
//			t, err := time.Parse(layout, time_to_parse)
//			if err != nil {
//				continue
//			} else {
//				return t
//			}
//		}
//		return time.
//	}
func savePost(db *database.Queries, item RSSItem, feed database.Feed) {
	description := sql.NullString{String: "", Valid: false}
	if len(item.Description) > 0 {
		description = sql.NullString{String: item.Description, Valid: true}
	}
	pubdate, err := time.Parse(time.RFC1123, item.PubDate)
	if err != nil {
		log.Printf("Could not parse published date: %v", err)
		return
	}
	_, err = db.CreatePost(context.Background(), database.CreatePostParams{
		ID:          uuid.New(),
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
		Title:       item.Title,
		Url:         item.Link,
		Description: description,
		PublishedAt: pubdate,
		FeedID:      feed.ID,
	})
	if err != nil {
		log.Printf("Could not save post: %v", err)
		return
	}
}
