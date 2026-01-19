package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/zombfeed/GoBlogAggregator/internal/database"
)

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.Name)
	}
	feed, err := s.db.GetFeedByURL(context.Background(), cmd.Args[0])
	if err != nil {
		return fmt.Errorf("coud not get feed: %w", err)
	}
	feedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("could not create feed follow: %w", err)
	}
	printFeedFollow(feedFollow.UserName, feedFollow.FeedName)
	return nil
}

func handlerListFeedFollowers(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %s", cmd.Name)
	}
	follows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}
	if len(follows) == 0 {
		fmt.Println("No feed follows found for this user.")
		return nil
	}

	fmt.Printf("Feed follows for user %s:\n", user.Name)
	for _, feeds := range follows {
		fmt.Printf("* %s\n", feeds.FeedName)
	}
	fmt.Println("********************************************")
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.Name)
	}
	feed, err := s.db.GetFeedByURL(context.Background(), cmd.Args[0])
	if err != nil {
		return fmt.Errorf("failed to find feed: %w", err)
	}
	err = s.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		FeedID: feed.ID,
		UserID: user.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to delete feed: %w", err)
	}
	fmt.Println("Feed deleted successfully")
	return nil
}

func printFeedFollow(username, feedname string) {
	fmt.Printf("* User:			%s\n", username)
	fmt.Printf("* Feed:			%s\n", feedname)
}
