package main

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"

	"github.com/matavak/gator/internal/database"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func PrintModelPretty[T any](model T) {
	data, err := json.MarshalIndent(model, "", " ")
	if err != nil {
		fmt.Printf("Could not print model %v\n", err)
		return
	}
	fmt.Println("============================================")
	fmt.Println(string(data))
	fmt.Println("============================================")
}

func handleAddFeed(s *state, cmd command, dbUser database.User) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("wrong arguments provided.usage:%s <feed name> <feed url>", cmd.Name)
	}
	feedName := cmd.Args[0]
	feedURL := cmd.Args[1]

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		Name:   feedName,
		Url:    feedURL,
		UserID: dbUser.ID,
	})
	if err != nil {
		return fmt.Errorf("Error creating feed:%v", err)
	}
	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		UserID: dbUser.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("could not create feedFollow %v", err)
	}
	fmt.Printf("Created feed:%v \n", feed)
	return nil
}

func handleRetrieveFeeds(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("wrong arguments provided.usage:%s ", cmd.Name)
	}
	feeds, err := s.db.GetAllFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("error retrieving all feeds: %v", err)
	}
	if len(feeds) == 0 {
		fmt.Printf("0 feeds found in DB")
		return nil
	}
	fmt.Printf("found %d feeds\n", len(feeds))
	for _, feed := range feeds {
		userID := feed.UserID
		user, err := s.db.GetUserByID(context.Background(), userID)
		if err != nil {
			return fmt.Errorf("error retrieving user: %v", err)
		}
		fmt.Printf("%s\n%s\n%s\n", feed.Name, feed.Url, user.Name)
	}
	return nil
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "gator")
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Could not read resp body %v", err)
	}
	rssFeed := &RSSFeed{}
	err = xml.Unmarshal([]byte(body), &rssFeed)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling xml:%v", err)
	}
	rssFeed.UnescapeStringRSSFeed()
	return rssFeed, nil
}

func (r *RSSFeed) UnescapeStringRSSFeed() {
	r.Channel.Title = html.UnescapeString(r.Channel.Title)
	r.Channel.Description = html.UnescapeString(r.Channel.Description)
	for i, v := range r.Channel.Item {
		r.Channel.Item[i].Title = html.UnescapeString(v.Title)
		r.Channel.Item[i].Description = html.UnescapeString(v.Description)
	}
}
