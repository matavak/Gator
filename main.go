package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/matavak/gator/internal/config"
	"github.com/matavak/gator/internal/database"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatalf("could not open postgres db: %v", err)
	}
	dbQueries := database.New(db)
	programState := &state{
		Conf: &cfg,
		db:   dbQueries,
	}
	programCommands := commands{
		Commands: map[string]func(*state, command) error{},
	}
	programCommands.register("login", handlerLogin)
	programCommands.register("register", handlerRegister)
	programCommands.register("reset", handlerReset)
	programCommands.register("users", handlerGetUsers)
	programCommands.register("agg", handleAgg)
	programCommands.register("addfeed", middlewareLoggedIn(handleAddFeed))
	programCommands.register("feeds", handleRetrieveFeeds)
	programCommands.register("follow", middlewareLoggedIn(handleFollow))
	programCommands.register("following", middlewareLoggedIn(handleFollowing))
	programCommands.register("unfollow", middlewareLoggedIn(handleUnfollow))
	programCommands.register("browse", middlewareLoggedIn(handleBrowse))
	args := os.Args
	if len(args) < 2 {
		log.Fatal("less than 2 arguments found in os.Args")
	}
	args = args[1:]
	cmd := command{
		Name: args[0],
		Args: args[1:],
	}
	err = programCommands.run(programState, cmd)
	if err != nil {
		log.Fatalf("error running command: %v", err)
	}
}

func middlewareLoggedIn(handler func(s *state, cmd command, dbUser database.User) error) func(*state, command) error {
	return func(s *state, c command) error {
		currentUser := s.Conf.CurrentUserName
		dbUser, err := s.db.GetUser(context.Background(), currentUser)
		if err != nil {
			return fmt.Errorf("Error getting username:%s from DB. ERR:%v", currentUser, err)
		}
		return handler(s, c, dbUser)
	}
}

type state struct {
	Conf *config.Config
	db   *database.Queries
}
