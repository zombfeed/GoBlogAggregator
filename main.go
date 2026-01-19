package main

import (
	"context"
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/zombfeed/GoBlogAggregator/internal/config"
	"github.com/zombfeed/GoBlogAggregator/internal/database"
)

const dbURL = "postgres://postgres:postgres@localhost:5432/gator"

type state struct {
	db     *database.Queries
	config *config.Config
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.config.CurrentUserName)
		if err != nil {
			return err
		}
		return handler(s, cmd, user)
	}
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("error connecing to db %v", err)
	}
	defer db.Close()

	dbQueries := database.New(db)
	s := state{dbQueries, &cfg}
	c := commands{registeredCommands: make(map[string]func(*state, command) error)}
	c.register("login", handlerLogin)
	c.register("register", handlerRegister)
	c.register("reset", handlerReset)
	c.register("users", handlerUsers)
	c.register("agg", handlerAggregator)
	c.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	c.register("feeds", handlerListFeeds)
	c.register("follow", middlewareLoggedIn(handlerFollow))
	c.register("following", middlewareLoggedIn(handlerListFeedFollowers))
	if len(os.Args) < 2 {
		log.Fatalf("usage: cli <command> [args...]")
	}
	cmd := command{
		Name: os.Args[1],
		Args: os.Args[2:],
	}
	if err := c.run(&s, cmd); err != nil {
		log.Fatal(err)
	}
}
