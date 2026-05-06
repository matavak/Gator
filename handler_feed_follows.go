package main

import (
	"context"
	"fmt"

	"github.com/matavak/gator/internal/database"
)

func handleUnfollow(s *state, cmd command, dbUser database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("wrong arguments provided.usage:%s  <feed url>", cmd.Name)
	}
	feedURL := cmd.Args[0]

	feed, err := s.db.GetFeedByURL(context.Background(), feedURL)
	if err != nil {
		return fmt.Errorf("Error getting feed:%v", err)
	}

	err = s.db.UnfollowFeed(context.Background(), database.UnfollowFeedParams{
		UserID: dbUser.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("could not unfollow feed %v", err)
	}

	return nil
}

func handleFollow(s *state, cmd command, dbUser database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("wrong arguments provided.usage:%s  <feed url>", cmd.Name)
	}
	feedURL := cmd.Args[0]

	feed, err := s.db.GetFeedByURL(context.Background(), feedURL)
	if err != nil {
		return fmt.Errorf("Error getting feed:%v", err)
	}

	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		UserID: dbUser.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("could not create feedFollow %v", err)
	}
	fmt.Printf("Created feed follow. Feed Name:%s \n User name:%s \n", feed.Name, dbUser.Name)
	return nil
}

func handleFollowing(s *state, cmd command, dbUser database.User) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("wrong arguments provided.usage:%s", cmd.Name)
	}

	userFeeds, err := s.db.GetFeedFollowsForUser(context.Background(), dbUser.ID)
	if err != nil {
		return fmt.Errorf("could not retreive user feed follows:%v", err)
	}
	for _, feed := range userFeeds {
		fmt.Println(feed.FeedName)
	}
	return nil
}
