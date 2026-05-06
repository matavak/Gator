package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"
	"github.com/matavak/gator/internal/database"
)

func scrapeFeeds(s *state) error {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("Could not get next feed to fetch:%v", err)
	}
	err = s.db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		return fmt.Errorf("could not mark feed as fetched %v", err)
	}
	rssFeed, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		return err
	}
	TimeLayouts := []string{
		time.RFC1123Z,
		time.RFC1123,
		time.RFC822Z,
		time.RFC822,
		time.RFC3339,
	}
	for _, item := range rssFeed.Channel.Item {
		pubDate := sql.NullTime{
			Valid: false,
		}
		for _, layout := range TimeLayouts {
			timeVal, err := time.Parse(layout, item.PubDate)
			if err == nil {
				pubDate = sql.NullTime{
					Time:  timeVal,
					Valid: true,
				}
			}
		}
		desc := sql.NullString{
			Valid: false,
		}
		if item.Description != "" {
			desc = sql.NullString{
				Valid:  true,
				String: item.Description,
			}
		}
		_, err = s.db.CreatePost(context.Background(), database.CreatePostParams{
			FeedID:      feed.ID,
			Title:       item.Title,
			Url:         item.Link,
			Description: desc,
			PublishedAt: pubDate,
		})
		if err != nil {
			var pqErr *pq.Error
			if errors.As(err, &pqErr) {
				if pqErr.Code == "23505" {
					continue
				}
			}
			fmt.Printf("error creating post,%v\n", err)
		}
	}

	return nil
}

func handleAgg(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("wrong arguments.usage: %s <time between reqs> ", cmd.Name)
	}
	timeBetweenRequests, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("Could not parse time between reqs %v", err)
	}
	fmt.Printf("collecting feeds every %v\n", timeBetweenRequests)
	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		err = scrapeFeeds(s)
		if err != nil {
			return err
		}
	}
}
