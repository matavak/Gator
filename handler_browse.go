package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/matavak/gator/internal/database"
)

func handleBrowse(s *state, cmd command, dbUser database.User) error {
	browseLimit := 2
	if len(cmd.Args) == 1 {
		specifiedLimit, err := strconv.Atoi(cmd.Args[0])
		if err != nil {
			return fmt.Errorf("could not parse browse limit %v", err)
		}
		browseLimit = specifiedLimit
	}
	if len(cmd.Args) > 1 {
		return fmt.Errorf("incorrect arguemtns, usage:%s [limit]", cmd.Name)
	}
	userPosts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: dbUser.ID,
		Limit:  int32(browseLimit),
	})
	if err != nil {
		return fmt.Errorf("could not get posts for user %v", err)
	}
	PrintModelPretty(userPosts)
	return nil
}
