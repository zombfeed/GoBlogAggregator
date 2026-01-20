package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/zombfeed/GoBlogAggregator/internal/database"
)

func handlerBrowse(s *state, cmd command, user database.User) error {
	if len(cmd.Args) > 1 {
		return fmt.Errorf("usage: %s *<optional_limit>", cmd.Name)
	}
	var limit int32
	if len(cmd.Args) == 0 {
		limit = 2
	} else {
		i64, err := strconv.ParseInt(cmd.Args[0], 10, 32)
		if err != nil {
			return fmt.Errorf("invalid limit: %w", err)
		}
		limit = int32(i64)
	}
	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  limit,
	})
	if err != nil {
		return fmt.Errorf("could not get posts")
	}
	fmt.Printf("Found %d posts:\n", len(posts))
	for _, post := range posts {
		printPosts(post)
	}
	return nil
}

func printPosts(post database.Post) {
	fmt.Printf("* Title: 	%s\n", post.Title)
	fmt.Printf("* Description:		%s\n", post.Description.String)
	fmt.Println("**************************************************")
}
